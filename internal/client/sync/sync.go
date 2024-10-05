package sync

import (
	"bytes"
	"context"
	cfg "gophKeeper/internal/client/config"
	"gophKeeper/internal/client/model"
	"gophKeeper/internal/client/service"
	pb "gophKeeper/internal/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SyncService interface {
	List(context.Context, model.ListRequest) (model.ListResponse, error)
	SyncItem(context.Context, *model.ItemSync) error

	SyncUser(context.Context, string) error
	DeleteUser(context.Context) error
	Close() error
}

var _ SyncService = (*syncService)(nil)

type syncService struct {
	s       service.Service
	conn    *grpc.ClientConn
	callOpt []grpc.CallOption
}

func NewSyncService(ctx context.Context, addr string, token []byte, s service.Service) (context.Context, SyncService, error) {
	sync := &syncService{
		s: s,
	}
	var err error
	ctx, sync.conn, sync.callOpt, err = dial(ctx, addr, map[string]string{
		pb.TokenKey: string(token),
	})
	return ctx, sync, err
}

func (sync syncService) Close() error {
	return sync.conn.Close()
}

func (sync syncService) List(ctx context.Context, request model.ListRequest) (list model.ListResponse, err error) {
	req := &pb.ListRequest{
		Limit:  request.Limit,
		Offset: request.Offset,
	}
	var res *pb.ListResponse
	client := pb.NewDataClient(sync.conn)
	res, err = client.List(ctx, req, sync.callOpt...)
	if err != nil {
		return
	}
	list = model.ListResponse{
		Total: res.Total,
	}
	list.Items = make([]model.DBItem, len(res.Items))
	for i, item := range res.Items {
		list.Items[i] = model.DBItem{
			Key:         item.Key,
			Description: item.Description,
		}
		if item.CreatedAt != nil {
			list.Items[i].CreatedAt = item.CreatedAt.AsTime()
		}
		if item.UpdatedAt != nil {
			list.Items[i].UpdatedAt = new(time.Time)
			*list.Items[i].UpdatedAt = item.UpdatedAt.AsTime()
		}
	}
	// todo run goroutines
	return
}

func (sync syncService) SyncItem(ctx context.Context, item *model.ItemSync) (err error) {
	pbItem := &pb.ItemSync{
		Key:         item.Key,
		Description: item.Description,
	}
	if !item.CreatedAt.IsZero() {
		pbItem.CreatedAt = timestamppb.New(item.CreatedAt)
	}
	if item.UpdatedAt != nil {
		pbItem.UpdatedAt = timestamppb.New(*item.UpdatedAt)
	}
	if item.Blob != nil {
		pbItem.Blob = item.Blob
	}

	client := pb.NewDataClient(sync.conn)
	pbItem, err = client.SyncItem(ctx, pbItem, sync.callOpt...)
	if err != nil {
		return
	}
	item.Description = pbItem.Description
	item.CreatedAt = pbItem.CreatedAt.AsTime()
	if pbItem.UpdatedAt != nil {
		if item.UpdatedAt == nil {
			item.UpdatedAt = new(time.Time)
		}
		*item.UpdatedAt = pbItem.UpdatedAt.AsTime()
	}

	// todo save it

	return
}

func (sync syncService) SyncUser(ctx context.Context, newPass string) (err error) {
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

	client := pb.NewUserClient(sync.conn)
	getUser, err = client.SyncUser(ctx, user, sync.callOpt...)
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

func (sync syncService) DeleteUser(ctx context.Context) (err error) {
	client := pb.NewUserClient(sync.conn)
	_, err = client.DeleteUser(ctx, &pb.NoMessage{}, sync.callOpt...)
	return

}
