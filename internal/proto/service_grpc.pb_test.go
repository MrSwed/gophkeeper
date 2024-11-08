package proto

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnimplemented(t *testing.T) {

	t.Run("generated grpc proto not implemented", func(t *testing.T) {
		ctx := context.Background()
		data := UnimplementedDataServer{}
		var err error
		_, err = UnimplementedAuthServer{}.RegisterClient(ctx, nil)
		assert.Error(t, err, status.Errorf(codes.Unimplemented, "method RegisterClient not implemented"))

		_, err = data.List(ctx, nil)
		assert.Error(t, err, status.Errorf(codes.Unimplemented, "method List not implemented"))
		_, err = data.SyncItem(ctx, nil)
		assert.Error(t, err, status.Errorf(codes.Unimplemented, "method SyncItem not implemented"))

		user := UnimplementedUserServer{}
		_, err = user.SyncUser(ctx, nil)
		assert.Error(t, err, status.Errorf(codes.Unimplemented, "method SyncUser not implemented"))
		_, err = user.DeleteUser(ctx, nil)
		assert.Error(t, err, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented"))
	})
}
