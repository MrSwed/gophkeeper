package sync

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func dial(ctx context.Context, addr string, meta map[string]string) (
	ctxOut context.Context,
	conn *grpc.ClientConn,
	callOpt []grpc.CallOption,
	err error) {
	ctxOut = ctx
	conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	if meta != nil {
		metaD := metadata.New(meta)
		ctxOut = metadata.NewOutgoingContext(ctx, metaD)
		callOpt = append(callOpt, grpc.Header(&metaD))
	}

	return
}
