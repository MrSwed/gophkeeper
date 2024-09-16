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
	"google.golang.org/grpc/peer"
)

type auth struct {
	pb.UnimplementedAuthServer
	s   service.Service
	log *zap.Logger
	c   *config.Config
}

var _ pb.AuthServer = (*auth)(nil)

func NewAuthServer(s service.Service, c *config.Config, log *zap.Logger) *auth {
	return &auth{
		s:   s,
		log: log,
		c:   c,
	}
}

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
			// register new one with login data, rest at sync
			u := &model.User{
				Email:    in.GetEmail(),
				Password: in.GetPassword(),
			}
			err = g.s.SaveSelf(ctx, u)
			if err != nil {
				return
			}
			out.AppToken, err = g.s.GetClientToken(ctx, req)
			return
		}
		return
	}
	return
}
