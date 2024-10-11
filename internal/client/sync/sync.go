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
		startTime  = time.Now()
		serverList *pb.ListResponse
		clientList out.List

		g          *errgroup.Group
		numWorkers = runtime.NumCPU()
		syncList   syncList
	)
	g, ctx = errgroup.WithContext(ctx)

	// Stage 1. collect list of
	// get server list
	g.Go(func() (err error) {
		var request = &pb.ListRequest{
			Limit:  cfg.PageSize,
			Offset: 0,
			// Orderby: "key",
		}
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			serverList, err = sc.dataClient.List(ctx, request, sc.callOpt...)
			if err != nil {
				return
			}
			for _, item := range serverList.GetItems() {
				if item.UpdatedAt == nil {
					if item.CreatedAt != nil {
						syncList.ToSync(item.Key, item.CreatedAt)
					}
				} else {
					syncList.ToSync(item.Key, item.UpdatedAt)
				}
			}
			if serverList.Total <= request.Offset+request.Limit {
				break
			}
			request.Offset += request.Limit
		}
		return
	})

	// get client list
	g.Go(func() (err error) {
		var request = model.ListQuery{
			Limit:  cfg.PageSize,
			Offset: 0,
			// Orderby: "key",
			SyncAt: startTime.String(),
		}

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
	})
	err = g.Wait()

	keysQueue := make(chan string)
	syncCount := int64(0)
	syncList.Range(func(key, _ interface{}) bool {
		keysQueue <- key.(string)
		syncCount++
		if syncCount >= syncList.Len() {
			close(keysQueue)
		}
		return true
	})
	itemSendQueue := make(chan model.DBRecord)
	itemSaveQueue := make(chan model.DBRecord)
	resultQueue := make(chan struct{})
	// run workers for get locals data
	for i := 0; i < numWorkers; i++ {
		g.Go(func() (err error) {
			defer close(itemSendQueue)
			for key := range keysQueue {
				item, er := sc.s.GetRaw(key)
				if er != nil && !errors.Is(er, sql.ErrNoRows) {
					err = errors.Join(err, er)
					continue
				}
				if item.Key == "" {
					item.Key = key
				}
				itemSendQueue <- item
			}
			return
		})
	}

	// run workers for send to sync
	for i := 0; i < numWorkers; i++ {
		g.Go(func() (err error) {
			defer close(itemSaveQueue)
			for itemSend := range itemSendQueue {
				itemGetPb, er := sc.dataClient.SyncItem(ctx, itemSend.ToItemSync())
				if er != nil {
					err = errors.Join(err, er)
					continue
				}
				var itemGet model.DBRecord
				itemGet.FromItemSync(itemGetPb)
				if itemSend.UpdatedAt != itemGet.UpdatedAt || itemSend.CreatedAt.IsZero() {
					itemSaveQueue <- itemGet
				}
			}
			return
		})
	}

	// run workers for save local
	for i := 0; i < numWorkers; i++ {
		g.Go(func() (err error) {
			defer close(resultQueue)
			for itemGet := range itemSaveQueue {
				er := sc.s.SaveRaw(itemGet)
				if er != nil {
					err = errors.Join(err, er)
					continue
				}
				resultQueue <- struct{}{}
			}
			return
		})
	}
	if er := g.Wait(); err != nil {
		err = errors.Join(err, er)
	}
	var resultCount int
	for _ = range resultQueue {
		resultCount++
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
