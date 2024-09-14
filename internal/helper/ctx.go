package helper

import (
	"context"
	"errors"
	"gophKeeper/internal/server/constant"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func GetCtxUserID(ctx context.Context) (uuid.UUID, error) {
	u, ok := ctx.Value(constant.CtxUserID).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user not found at ctx")
	}
	return u, nil
}

func GetGRPCCtxToken(ctx context.Context) (token []byte) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(constant.TokenKey); len(values) > 0 {
			token = []byte(values[0])
		}
	}
	return
}
