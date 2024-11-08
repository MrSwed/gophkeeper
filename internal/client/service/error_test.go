package service

import (
	"errors"
	"fmt"
	"gophKeeper/internal/client/model"
	"gophKeeper/internal/client/model/type/card"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServiceError(t *testing.T) {
	type args struct {
		e error
	}
	tests := []struct {
		args args
		name string
	}{
		{
			name: "test error stub",
			args: args{fmt.Errorf("error load current user profile: %v", "some error")},
		},
		{
			name: "test error stub",
			args: args{errors.New("error db_file - is not set")},
		},
		{
			name: "test error stub",
			args: args{fmt.Errorf("open sqlite error %s dbFile %s\n", "some error", "some dbFile")},
		},
		{
			name: "test error stub",
			args: args{fmt.Errorf("db update error: %s dbFile %s\n", "some error", "some dbFile")},
		},
		{
			name: "test error stub",
			args: args{fmt.Errorf("usrCfgDir error: %s \n", "some error")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewServiceError(tt.args.e)

			err := srv.Save(&card.Model{})
			assert.Equal(t, err, tt.args.e, "Save()")

			_, err = srv.Get("")
			assert.Equal(t, err, tt.args.e, "Get()")

			_, err = srv.List(model.ListQuery{})
			assert.Equal(t, err, tt.args.e, "List()")

			err = srv.Delete("")
			assert.Equal(t, err, tt.args.e, "Delete()")

			err = srv.ChangePasswd()
			assert.Equal(t, err, tt.args.e, "ChangePasswd()")

			_, err = srv.GetToken()
			assert.Equal(t, err, tt.args.e, "GetToken()")
		})
	}
}
