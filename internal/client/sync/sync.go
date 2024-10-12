package sync

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	cfg "gophKeeper/internal/client/config"
	"gophKeeper/internal/client/model"
	"gophKeeper/internal/client/model/out"
	"gophKeeper/internal/client/service"
	pb "gophKeeper/internal/proto"
	"runtime"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SyncService interface {
	SyncData(ctx context.Context) (err error)

	SyncUser(context.Context, string) error
	DeleteUser(context.Context) error
	Close() error
}

var _ SyncService = (*syncService)(nil)

type syncService struct {
	s          service.Service
	conn       *grpc.ClientConn
	callOpt    []grpc.CallOption
	dataClient pb.DataClient
}

func NewSyncService(ctx context.Context, addr string, token []byte, s service.Service) (context.Context, SyncService, error) {
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
		g          *errgroup.Group
		numWorkers = runtime.NumCPU()
		syncList   = &syncList{
			startTime: time.Now(),
		}
	)
	g, ctx = errgroup.WithContext(ctx)

	// Stage 1. collect list with sync needed items
	// get server list
	g.Go(sc.getRemoteCollect(ctx, syncList))

	// get client list
	g.Go(sc.getLocalCollect(ctx, syncList))

	// wait for full complete sync list
	err = g.Wait()
	// send keys for sync

	errCh := make(chan error)
	defer close(errCh)

	keysQueue := syncList.KeyQueue()

	localItemQueues := make([]chan model.DBRecord, numWorkers)
	// run workers for get locals data
	for i := 0; i < numWorkers; i++ {
		localItemQueues[i] = sc.getItemChan(ctx, keysQueue, errCh) // return localItemsQueue
	}

	// send to sync
	// syncItemsQueue := make(chan model.DBRecord)
	remoteItemsQueue := sc.syncItemsHandler(ctx, errCh, localItemQueues...)

	// run workers for save local
	resultQueue := sc.saveUpdatedItems(ctx, remoteItemsQueue, errCh)

	var resultCount int
	for {
		select {
		case _, ok := <-resultQueue:
			if ok {
				resultCount++
			} else {
				resultQueue = nil
			}
		case er, ok := <-errCh:
			if ok {
				err = errors.Join(err, er)
			} else {
				errCh = nil
			}
		}
		if resultQueue == nil && errCh == nil {
			break
		}
	}
	cfg.User.Set("sync.status.data.last_sync_at", time.Now())
	cfg.User.Set("sync.status.data.updated", resultCount)
	return
}

func (sc syncService) SyncUser(ctx context.Context, newPass string) (err error) {
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
	if updatedAt := cfg.User.GetTime("sync.user.updated_at"); !updatedAt.IsZero() {
		user.UpdatedAt = timestamppb.New(updatedAt)
	}

	client := pb.NewUserClient(sc.conn)
	getUser, err = client.SyncUser(ctx, user, sc.callOpt...)
	if err != nil {
		return
	}
	var updated bool
	if !bytes.Equal(user.PackedKey, getUser.PackedKey) {
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
	if getUser.CreatedAt != nil && (user.CreatedAt == nil || user.CreatedAt != getUser.CreatedAt) {
		cfg.User.Set("sync.user.created_at", getUser.CreatedAt.AsTime())
		updated = true
	}
	if getUser.UpdatedAt != nil {
		cfg.User.Set("sync.user.updated_at", getUser.UpdatedAt.AsTime())
		updated = true
	}
	if updated {
		cfg.User.Set("sync.status.user.last_sync_at", time.Now())
	}
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
				// Orderby: "key",
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
				// Orderby: "key",
				SyncAt: syncList.startTime.String(),
			}
		)
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
					syncList.ToSync(item.Key, timestamppb.New(item.CreatedAt))
				} else {
					syncList.ToSync(item.Key, timestamppb.New(*item.UpdatedAt))
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

func (sc syncService) getItemChan(ctx context.Context, keysQueue chan string,
	errCh chan error) (lcItemQueue chan model.DBRecord) {
	lcItemQueue = make(chan model.DBRecord)
	go func() {
		defer close(lcItemQueue)
		for key := range keysQueue {
			select {
			case <-ctx.Done():
			default:
				item, er := sc.s.GetRaw( /*ctx,*/ key)
				if er != nil && !errors.Is(er, sql.ErrNoRows) {
					// go func() {
					errCh <- er
					// }()
					continue
				}
				if item.Key == "" {
					item.Key = key
				}
				lcItemQueue <- item
			}
		}
	}()

	return lcItemQueue
}

func (sc syncService) syncItemsHandler(ctx context.Context, errCh chan error,
	localItemQueues ...chan model.DBRecord) (remoteItemsQueue chan model.DBRecord) {
	var (
		g *errgroup.Group
	)
	g, ctx = errgroup.WithContext(ctx)
	remoteItemsQueue = make(chan model.DBRecord)
	for _, lcItemQueue := range localItemQueues {
		// go less 1.22
		lcItemQueue := lcItemQueue
		g.Go(func() (err error) {
			for itemSend := range lcItemQueue {
				itemGetPb, er := sc.dataClient.SyncItem(ctx, itemSend.ToItemSync())
				if er != nil {
					err = errors.Join(err, er)
					continue
				}
				var itemGet model.DBRecord
				itemGet.FromItemSync(itemGetPb)
				remoteItemsQueue <- itemGet
				// we do not need check it again, because we filtered the list at prev stage
				// if itemSend.UpdatedAt != itemGet.UpdatedAt || itemSend.CreatedAt.IsZero() {
				// 	remoteItemsQueue <- itemGet
				// }
			}
			return
		})
	}
	go func() {
		err := g.Wait()
		close(remoteItemsQueue)
		if err != nil {
			errCh <- err
		}
	}()

	return remoteItemsQueue
}

func (sc syncService) saveUpdatedItems(ctx context.Context, itemSaveQueue chan model.DBRecord, errCh chan error) (resultQueue chan struct{}) {
	resultQueue = make(chan struct{})
	go func() {
		defer close(resultQueue)
		for itemGet := range itemSaveQueue {
			select {
			case <-ctx.Done():
			default:
				er := sc.s.SaveRaw( /*ctx,*/ itemGet)
				if er != nil {
					go func() {
						errCh <- er
					}()
					continue
				}
				resultQueue <- struct{}{}
			}
		}
	}()

	return resultQueue
}
