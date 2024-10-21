package grpc

import (
	"context"
	pb "gophKeeper/internal/proto"
	"gophKeeper/internal/server/config"
	errs "gophKeeper/internal/server/errors"
	"gophKeeper/internal/server/model"
	"gophKeeper/internal/server/service"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type user struct {
	pb.UnimplementedUserServer
	s   service.User
	log *zap.Logger
	c   *config.Config
}

var _ pb.UserServer = (*user)(nil)

func NewUserServer(s service.User, c *config.Config, log *zap.Logger) *user {
	return &user{
		s:   s,
		log: log,
		c:   c,
	}
}

func (g *user) SyncUser(ctx context.Context, in *pb.UserSync) (out *pb.UserSync, err error) {
	out = in
	defer func() { out.Password = "" }()
	pass := model.PassChangeRequest{Password: in.GetPassword()}
	if err = pass.Validate(); err != nil {
		err = status.Error(codes.InvalidArgument, err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(ctx, g.c.GRPCOperationTimeout)
	defer cancel()
	syncKey := in.GetEmail()
	if syncKey == "" {
		err = status.Error(codes.Internal, errs.ErrorSyncNoKey.Error())
		return
	}

	var storedUser model.User
	storedUser, err = g.s.GetSelf(ctx)
	if err != nil {
		err = status.Error(codes.Internal, err.Error())
		return
	}

	// it is same data
	if storedUser.PackedKey != nil &&
		in.GetCreatedAt().AsTime().Equal(storedUser.CreatedAt) &&
		(storedUser.UpdatedAt == nil || in.GetUpdatedAt().AsTime().Equal(*storedUser.UpdatedAt)) {
		return
	}

	// incoming data is newest - update server store
	if (in.GetUpdatedAt().IsValid() && ((storedUser.UpdatedAt != nil &&
		in.UpdatedAt.AsTime().After(*storedUser.UpdatedAt)) ||
		storedUser.UpdatedAt == nil)) || storedUser.CreatedAt.IsZero() {
		storedUser.Description = nil
		if in.GetDescription() != "" {
			storedUser.Description = new(string)
			*storedUser.Description = in.GetDescription()
		}
		storedUser.CreatedAt = in.GetCreatedAt().AsTime()
		storedUser.PackedKey = in.GetPackedKey()
		storedUser.Password = in.GetPassword()
		storedUser.UpdatedAt = nil
		if in.GetUpdatedAt().IsValid() {
			storedUser.UpdatedAt = new(time.Time)
			*storedUser.UpdatedAt = in.GetUpdatedAt().AsTime()
		}
		err = g.s.SaveSelf(ctx, &storedUser)
		if err != nil {
			err = status.Error(codes.Internal, err.Error())
		}

		return
	}

	// incoming is oldest, return from server store
	out.PackedKey = storedUser.PackedKey
	out.Description = ""
	if storedUser.Description != nil {
		out.Description = *storedUser.Description
	}
	out.CreatedAt = timestamppb.New(storedUser.CreatedAt)
	out.UpdatedAt = nil
	if storedUser.UpdatedAt != nil {
		out.UpdatedAt = timestamppb.New(*storedUser.UpdatedAt)
	}

	return
}

func (g *user) DeleteUser(ctx context.Context, _ *pb.NoMessage) (out *pb.OkResponse, err error) {
	err = g.s.DeleteSelf(ctx)
	out = &pb.OkResponse{Ok: err == nil}
	if err != nil {
		err = status.Error(codes.Internal, err.Error())
	}

	return
}
