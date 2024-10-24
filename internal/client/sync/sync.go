package sync

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"runtime"
	"sync"
	"time"

	cfg "gophKeeper/internal/client/config"
	"gophKeeper/internal/client/model"
	"gophKeeper/internal/client/model/out"
	"gophKeeper/internal/client/service"
	pb "gophKeeper/internal/proto"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service interface {
	SyncData(ctx context.Context) (err error)

	SyncUser(context.Context, string) (bool, error)
	DeleteUser(context.Context) error
	Close() error
}

var _ Service = (*syncService)(nil)

type syncService struct {
	conn       *grpc.ClientConn
	callOpt    []grpc.CallOption
	dataClient pb.DataClient
	s          service.Service
}

type dbRecQueue struct {
	item model.DBRecord
	err  error
}

func NewSyncService(ctx context.Context, addr string, token []byte, s service.Service) (context.Context, Service, error) {
	sc := &syncService{
		s: s,
	}
	var err error
	ctx, sc.conn, sc.callOpt, err = dial(ctx, addr, map[string]string{
		pb.TokenKey: string(token),
	})
	sc.dataClient = pb.NewDataClient(sc.conn)

	return ctx, sc, err
}

func (sc syncService) Close() error {
	return sc.conn.Close()
}

func (sc syncService) SyncData(ctx context.Context) (err error) {
	var (
		numWorkers = runtime.NumCPU()
		syncList   = &syncList{
			startTime: time.Now(),
		}
	)
	g, gctx := errgroup.WithContext(ctx)

	// Stage 1. collect list with sync needed items
	// get server list
	g.Go(sc.getRemoteCollect(gctx, syncList))

	// get client list
	g.Go(sc.getLocalCollect(gctx, syncList))

	// wait for full complete sync list
	err = g.Wait()
	// send keys for sync

	keysQueue := syncList.KeyQueue()

	numWorkers = min(numWorkers, int(syncList.Len()))
	localItemQueues := make([]chan dbRecQueue, numWorkers)
	// run workers for get locals data
	for i := 0; i < numWorkers; i++ {
		localItemQueues[i] = sc.getItemChan(ctx, keysQueue) // return localItemsQueue
	}

	// send to sync
	remoteItemsQueue := sc.syncItemsHandler(ctx, localItemQueues...)

	// run workers for save local
	resultQueue := sc.saveUpdatedItems(ctx, remoteItemsQueue)

	var resultCount int
	for er := range resultQueue {
		if er == nil {
			resultCount++
		} else {
			err = errors.Join(err, er)
		}
	}
	cfg.User.Set("sync.status.data.last_sync_at", time.Now())
	cfg.User.Set("sync.status.data.updated", resultCount)

	return
}

func (sc syncService) SyncUser(ctx context.Context, newPass string) (updated bool, err error) {
	var getUser *pb.UserSync
	user := &pb.UserSync{
		Email:       cfg.User.GetString("email"),
		PackedKey:   []byte(cfg.User.GetString("packed_key")),
		Description: cfg.User.GetString("sync.user.description"),
		Password:    newPass,
	}
	if createdAt := cfg.User.GetTime("sync.user.created_at"); !createdAt.IsZero() {
		user.CreatedAt = timestamppb.New(createdAt)
	}
	var updatedAt time.Time
	if newPass != "" {
		user.UpdatedAt = timestamppb.New(time.Now())
	} else if updatedAt = cfg.User.GetTime("sync.user.updated_at"); !updatedAt.IsZero() {
		user.UpdatedAt = timestamppb.New(updatedAt)
	}

	client := pb.NewUserClient(sc.conn)
	getUser, err = client.SyncUser(ctx, user, sc.callOpt...)
	if err != nil {
		return
	}
	if getUser.PackedKey != nil && !bytes.Equal(user.PackedKey, getUser.PackedKey) {
		user.PackedKey = getUser.PackedKey
		cfg.User.Set("packed_key", user.PackedKey)
		updated = true
	}
	if user.Description != getUser.Description {
		cfg.User.Set("sync.user.description", getUser.Description)
		updated = true
	}
	if getUser.Email != "" && user.Email != getUser.Email {
		cfg.User.Set("email", getUser.Email)
		updated = true
	}
	if getUser.CreatedAt != nil && (user.CreatedAt == nil || user.CreatedAt.AsTime() != getUser.CreatedAt.AsTime()) {
		cfg.User.Set("sync.user.created_at", getUser.CreatedAt.AsTime())
		updated = true
	}
	if updated || (getUser.UpdatedAt != nil &&
		(user.UpdatedAt == nil || user.UpdatedAt.AsTime() != getUser.UpdatedAt.AsTime())) {
		cfg.User.Set("sync.user.updated_at", getUser.UpdatedAt.AsTime())
		updated = true
	}
	cfg.User.Set("sync.status.user.last_sync_at", time.Now())

	return
}

func (sc syncService) DeleteUser(ctx context.Context) (err error) {
	client := pb.NewUserClient(sc.conn)
	_, err = client.DeleteUser(ctx, &pb.NoMessage{}, sc.callOpt...)
	return

}

func (sc syncService) getRemoteCollect(ctx context.Context, syncList *syncList) func() (err error) {
	return func() (err error) {
		var (
			remoteList *pb.ListResponse
			request    = &pb.ListRequest{
				Limit:  cfg.PageSize,
				Offset: 0,
			}
		)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			remoteList, err = sc.dataClient.List(ctx, request, sc.callOpt...)
			if err != nil {
				return
			}
			for _, item := range remoteList.GetItems() {
				if item.UpdatedAt == nil {
					if item.CreatedAt != nil {
						syncList.ToSync(item.Key, item.CreatedAt)
					}
				} else {
					syncList.ToSync(item.Key, item.UpdatedAt)
				}
			}
			if remoteList.Total <= request.Offset+request.Limit {
				break
			}
			request.Offset += request.Limit
		}
		return
	}
}

func (sc syncService) getLocalCollect(ctx context.Context, syncList *syncList) func() (err error) {
	return func() (err error) {
		var (
			clientList out.List
			request    = model.ListQuery{
				Limit:  cfg.PageSize,
				Offset: 0,
				SyncAt: syncList.startTime.Format(time.DateTime),
			}
		)
		_, z := time.Now().Zone()
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			clientList, err = sc.s.List(request)
			if err != nil {
				return
			}
			for _, item := range clientList.Items {
				if item.UpdatedAt == nil {
					syncList.ToSync(item.Key, timestamppb.New(item.CreatedAt.Add(-time.Duration(z)*time.Second)))
				} else {
					syncList.ToSync(item.Key, timestamppb.New((*item.UpdatedAt).Add(-time.Duration(z)*time.Second)))
				}
			}
			if clientList.Total <= request.Offset+request.Limit {
				break
			}
			request.Offset += request.Limit
		}
		return
	}
}

func (sc syncService) getItemChan(ctx context.Context, keysQueue chan string, /*,
errCh chan error*/) (lcItemQueue chan dbRecQueue) {
	lcItemQueue = make(chan dbRecQueue)
	go func() {
		defer close(lcItemQueue)
		for key := range keysQueue {
			select {
			case <-ctx.Done():
				return
			default:
				item, er := sc.s.GetRaw( /*ctx,*/ key)
				if er != nil {
					if errors.Is(er, sql.ErrNoRows) {
						er = nil
					}
				}
				if item.Key == "" {
					item.Key = key
				}
				lcItemQueue <- dbRecQueue{item, er}
			}
		}
	}()

	return lcItemQueue
}

func (sc syncService) syncItemsHandler(ctx context.Context,
	localItemQueues ...chan dbRecQueue) (remoteItemsQueue chan dbRecQueue) {
	var (
		g sync.WaitGroup
	)
	remoteItemsQueue = make(chan dbRecQueue)
	for _, lcItemQueue := range localItemQueues {
		// go less 1.22
		lcItemQueue := lcItemQueue
		g.Add(1)
		go func() {
			defer g.Done()
			for itemSend := range lcItemQueue {
				itemGetPb, er := sc.dataClient.SyncItem(ctx, itemSend.item.ToItemSync())
				var itemGet model.DBRecord
				if er == nil {
					itemGet.FromItemSync(itemGetPb)
				}
				remoteItemsQueue <- dbRecQueue{item: itemGet}
			}
		}()
	}
	go func() {
		g.Wait()
		close(remoteItemsQueue)
	}()

	return remoteItemsQueue
}

func (sc syncService) saveUpdatedItems(ctx context.Context, forSaveQueue chan dbRecQueue /*, errCh chan error*/) (resultQueue chan error) {
	resultQueue = make(chan error)
	go func() {
		defer close(resultQueue)
		for itemGet := range forSaveQueue {
			select {
			case <-ctx.Done():
			default:
				resultQueue <- sc.s.SaveRaw( /*ctx,*/ itemGet.item)
			}
		}
	}()

	return resultQueue
}
