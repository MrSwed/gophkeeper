package helper

import (
	"context"
	"errors"
	"gophKeeper/internal/server/constant"

	"github.com/google/uuid"
)

func GetCtxUserID(ctx context.Context) (uuid.UUID, error) {
	u, ok := ctx.Value(constant.CtxUserID).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user not found at ctx")
	}
	return u, nil
}
