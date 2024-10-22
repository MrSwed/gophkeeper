/*
This package provides the implementation of the gRPC user server for the GophKeeper application.
It defines methods for handling user operations such as synchronization and deletion.

Main functionalities include:

- Synchronizing user data between the client and server.
- Deleting user accounts from the server.
*/
package grpc

import (
	"context"
	"time"

	pb "gophKeeper/internal/proto"
	"gophKeeper/internal/server/config"
	errs "gophKeeper/internal/server/errors"
	"gophKeeper/internal/server/model"
	"gophKeeper/internal/server/service"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// user implements the UserServer interface defined in the protobuf file.
// It provides methods for user management, including synchronization and deletion.
type user struct {
	pb.UnimplementedUserServer
	s   service.User
	log *zap.Logger
	c   *config.Config
}

// Ensure that user implements the UserServer interface.
var _ pb.UserServer = (*user)(nil)

// NewUserServer creates a new instance of the user server.
func NewUserServer(s service.User, c *config.Config, log *zap.Logger) *user {
	return &user{
		s:   s,
		log: log,
		c:   c,
	}
}

// SyncUser handles the synchronization of user data.
// It takes a context and a UserSync request as input, and returns a UserSync response
// with the synchronized user data or an error if synchronization fails.
func (g *user) SyncUser(ctx context.Context, in *pb.UserSync) (out *pb.UserSync, err error) {
	out = in
	defer func() { out.Password = "" }() // Clear password from output
	pass := model.PassChangeRequest{Password: in.GetPassword()}
	if err = pass.Validate("Password"); err != nil {
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
	// If the incoming data is the same, do nothing
	if storedUser.PackedKey != nil &&
		in.GetCreatedAt().AsTime().Equal(storedUser.CreatedAt) &&
		(storedUser.UpdatedAt == nil || in.GetUpdatedAt().AsTime().Equal(*storedUser.UpdatedAt)) {
		return
	}
	// If incoming data is newer, update the server store
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
	// If incoming data is older, return from server store
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

// DeleteUser handles the deletion of the user account.
// It takes a context and a NoMessage request as input, and returns an OkResponse
// indicating whether the deletion was successful or not.
func (g *user) DeleteUser(ctx context.Context, _ *pb.NoMessage) (out *pb.OkResponse, err error) {
	err = g.s.DeleteSelf(ctx)
	out = &pb.OkResponse{Ok: err == nil}
	if err != nil {
		err = status.Error(codes.Internal, err.Error())
	}
	return
}
