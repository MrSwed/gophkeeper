package grpc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	pb "gophKeeper/internal/proto"
	"gophKeeper/internal/server/config"
	errs "gophKeeper/internal/server/errors"
	"gophKeeper/internal/server/model"
	"gophKeeper/internal/server/service"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type data struct {
	pb.UnimplementedDataServer
	s   service.Data
	log *zap.Logger
	c   *config.Config
}

var _ pb.DataServer = (*data)(nil)

func NewDataServer(s service.Data, c *config.Config, log *zap.Logger) *data {
	return &data{
		s:   s,
		log: log,
		c:   c,
	}
}

func (g *data) List(ctx context.Context, in *pb.ListRequest) (out *pb.ListResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, g.c.GRPCOperationTimeout)
	defer cancel()
	var (
		q    = new(model.ListQuery)
		list model.List
	)
	q.Offset = in.Offset
	q.Limit = in.Limit
	list, err = g.s.ListSelf(ctx, q)
	if err != nil {
		return
	}
	out = &pb.ListResponse{
		Total: list.Total,
		Items: make([]*pb.ItemShort, len(list.Items)),
	}
	for i, item := range list.Items {
		out.Items[i] = &pb.ItemShort{
			Key:       item.Key,
			CreatedAt: timestamppb.New(item.CreatedAt),
		}
		if item.UpdatedAt != nil {
			out.Items[i].UpdatedAt = timestamppb.New(*item.UpdatedAt)
		}
		if item.Description != nil {
			out.Items[i].Description = *item.Description
		}
	}
	return
}

func (g *data) SyncItem(ctx context.Context, in *pb.ItemSync) (out *pb.ItemSync, err error) {
	out = in
	ctx, cancel := context.WithTimeout(ctx, g.c.GRPCOperationTimeout)
	defer cancel()
	syncKey := in.GetKey()
	if syncKey == "" {
		err = errs.ErrorSyncNoKey
		return
	}

	var item *model.Item
	item, err = g.s.GetSelfItem(ctx, syncKey)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}

	if !item.IsNew() && in.CreatedAt != nil && !in.CreatedAt.AsTime().Equal(item.CreatedAt) {
		err = fmt.Errorf("%w key: %s", errs.ErrorSyncSameKey, syncKey)
		return
	}

	// it is same data
	if in.CreatedAt.AsTime().Equal(item.CreatedAt) && (item.UpdatedAt == nil || in.UpdatedAt.AsTime().Equal(*item.UpdatedAt)) {
		return
	}

	// incoming data is newest - update server store
	if item.CreatedAt.IsZero() || (in.GetUpdatedAt().IsValid() && ((item.UpdatedAt != nil &&
		in.UpdatedAt.AsTime().After(*item.UpdatedAt)) ||
		item.UpdatedAt == nil)) {
		if in.GetDescription() != "" {
			item.Description = new(string)
			*item.Description = in.GetDescription()
		}
		item.CreatedAt = in.GetCreatedAt().AsTime()
		item.Blob = in.GetBlob()
		if in.GetUpdatedAt().IsValid() {
			if item.UpdatedAt == nil {
				item.UpdatedAt = new(time.Time)
			}
			*item.UpdatedAt = in.GetUpdatedAt().AsTime()
		}
		err = g.s.SaveSelfItem(ctx, item)
		if err != nil {
			g.log.Error("save item failed", zap.Error(err))
		}
		return
	}

	// incoming is oldest or empty, return from server store
	out.Blob = item.Blob
	out.Description = ""
	if item.Description != nil {
		out.Description = *item.Description
	}
	out.CreatedAt = timestamppb.New(item.CreatedAt)
	if item.UpdatedAt != nil {
		out.UpdatedAt = timestamppb.New(*item.UpdatedAt)
	}
	return
}
