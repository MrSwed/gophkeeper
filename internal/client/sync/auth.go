package sync

import (
	"context"
	"errors"

	"gophKeeper/internal/client/model"
	pb "gophKeeper/internal/proto"

	"google.golang.org/grpc"
)

func RegisterClient(ctx context.Context, addr string, request model.RegisterClientRequest) (token []byte, err error) {
	req := &pb.RegisterClientRequest{
		Email:    request.Email,
		Password: request.Password,
	}
	var (
		result *pb.ClientToken
		conn   *grpc.ClientConn
	)
	ctx, conn, _, err = dial(ctx, addr, nil)
	defer func() {
		err = errors.Join(err, conn.Close())
	}()

	client := pb.NewAuthClient(conn)
	result, err = client.RegisterClient(ctx, req)
	if err != nil {
		return
	}
	token = result.AppToken
	return
}
