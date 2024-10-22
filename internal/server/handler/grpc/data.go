/*
This package provides the implementation of the gRPC data server for the GophKeeper application.
It defines methods for handling data operations such as listing and synchronizing items.

Main functionalities include:

- Listing stored items with pagination and ordering options.
- Synchronizing individual items between the client and server.
*/
package grpc

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "gophKeeper/internal/proto"
	"gophKeeper/internal/server/config"
	errs "gophKeeper/internal/server/errors"
	"gophKeeper/internal/server/model"
	"gophKeeper/internal/server/service"
)

// data implements the DataServer interface defined in the protobuf file.
// It provides methods for data management, including listing and synchronizing items.
type data struct {
	pb.UnimplementedDataServer
	s   service.Data
	log *zap.Logger
	c   *config.Config
}

// Ensure that data implements the DataServer interface.
var _ pb.DataServer = (*data)(nil)

// NewDataServer creates a new instance of the data server.
func NewDataServer(s service.Data, c *config.Config, log *zap.Logger) *data {
	return &data{
		s:   s,
		log: log,
		c:   c,
	}
}

// List handles the listing of stored items.
// It takes a context and a ListRequest as input and returns a ListResponse
// containing the total number of items and a list of item summaries.
func (g *data) List(ctx context.Context, in *pb.ListRequest) (out *pb.ListResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, g.c.GRPCOperationTimeout)
	defer cancel()
	var (
		q    = new(model.ListQuery)
		list model.List
	)
	q.Offset = in.GetOffset()
	q.Limit = in.GetLimit()
	if in.GetOrderby() != "" {
		q.OrderBy = in.GetOrderby()
	}
	list, err = g.s.ListSelf(ctx, q)
	if err != nil {
		err = status.Error(codes.Internal, err.Error())
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

// SyncItem handles the synchronization of an individual item.
// It takes a context and an ItemSync request as input, and returns an ItemSync response
// with the synchronized item data or an error if synchronization fails.
func (g *data) SyncItem(ctx context.Context, in *pb.ItemSync) (out *pb.ItemSync, err error) {
	out = in
	ctx, cancel := context.WithTimeout(ctx, g.c.GRPCOperationTimeout)
	defer cancel()
	syncKey := in.GetKey()
	if syncKey == "" {
		err = status.Error(codes.InvalidArgument, errs.ErrorSyncNoKey.Error())
		return
	}
	var item *model.Item
	item, err = g.s.GetSelfItem(ctx, syncKey)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		err = status.Error(codes.Internal, err.Error())
		return
	}
	if !item.IsNew() && in.GetCreatedAt() != nil && !in.GetCreatedAt().AsTime().Equal(item.CreatedAt) {
		err = status.Errorf(codes.Canceled, "%s key: %s", errs.ErrorSyncCreatedDate, syncKey)
		return
	}
	// If the incoming data is the same, do nothing
	if in.GetCreatedAt().AsTime().Equal(item.CreatedAt) &&
		(item.UpdatedAt == nil || in.GetUpdatedAt().AsTime().Equal(*item.UpdatedAt)) {
		return
	}
	// If incoming data is newer, update the server store
	if item.CreatedAt.IsZero() || (in.GetUpdatedAt().IsValid() && ((item.UpdatedAt != nil &&
		in.GetUpdatedAt().AsTime().After(*item.UpdatedAt)) ||
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
			err = status.Error(codes.Internal, err.Error())
		}
		return
	}
	// If incoming data is older or empty, return from server store
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
