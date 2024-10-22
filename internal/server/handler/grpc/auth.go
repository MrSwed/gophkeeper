/*
This package provides the implementation of the gRPC authentication server for the GophKeeper application.
It defines the methods for handling client registration and authentication.

Main functionalities include:

- Registering a new client.
- Retrieving client tokens for authenticated sessions.
*/
package grpc

import (
	"context"
	"database/sql"
	"errors"

	pb "gophKeeper/internal/proto"
	"gophKeeper/internal/server/config"
	"gophKeeper/internal/server/model"
	"gophKeeper/internal/server/service"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// auth implements the AuthServer interface defined in the protobuf file.
// It provides methods for client authentication and registration.
type auth struct {
	pb.UnimplementedAuthServer
	s   service.Service
	log *zap.Logger
	c   *config.Config
}

// Ensure that auth implements the AuthServer interface.
var _ pb.AuthServer = (*auth)(nil)

// NewAuthServer creates a new instance of the auth server.
func NewAuthServer(s service.Service, c *config.Config, log *zap.Logger) *auth {
	return &auth{
		s:   s,
		log: log,
		c:   c,
	}
}

// RegisterClient handles the registration of a new client.
// It takes a context and a RegisterClientRequest as input and returns a ClientToken
// or an error if the registration fails.
func (g *auth) RegisterClient(ctx context.Context, in *pb.RegisterClientRequest) (out *pb.ClientToken, err error) {
	ctx, cancel := context.WithTimeout(ctx, g.c.GRPCOperationTimeout)
	defer cancel()
	out = new(pb.ClientToken)
	var (
		meta, _ = peer.FromContext(ctx)
	)
	req := model.AuthRequest{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
		Meta:     meta,
	}
	out.AppToken, err = g.s.GetClientToken(ctx, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Register a new client with login data; the rest will be handled at sync.
			err = req.Validate()
			if err != nil {
				err = status.Error(codes.InvalidArgument, err.Error())
				return
			}
			u := &model.User{
				Email:    in.GetEmail(),
				Password: in.GetPassword(),
			}
			err = g.s.SaveSelf(ctx, u)
			if err != nil {
				err = status.Error(codes.Unauthenticated, err.Error())
				return
			}
			out.AppToken, err = g.s.GetClientToken(ctx, req)
		}
		if err != nil {
			err = status.Error(codes.Unauthenticated, err.Error())
		}
	}
	return
}
