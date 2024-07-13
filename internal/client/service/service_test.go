package service

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"gophKeeper/internal/client/config"
	clMigrate "gophKeeper/internal/client/migrate"
	"gophKeeper/internal/client/model"
	"gophKeeper/internal/client/model/type/auth"
	"gophKeeper/internal/client/model/type/bin"
	"gophKeeper/internal/client/model/type/card"
	"gophKeeper/internal/client/model/type/text"
	"gophKeeper/internal/client/storage"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type serviceStoreTestSuite struct {
	suite.Suite
	db  *sqlx.DB
	srv Service
}

func (suite *serviceStoreTestSuite) SetupSuite() {
	storePath := filepath.Join(suite.T().TempDir(), config.AppName)
	err := os.MkdirAll(storePath, os.ModePerm)
	require.NoError(suite.T(), err)
	dbFile := filepath.Join(storePath, "store.db")

	config.User.Set("encryption_key", "SomePhraseEncryptionKey")

	suite.db, err = sqlx.Open("sqlite3", dbFile)
	require.NoError(suite.T(), err)

	_, err = clMigrate.Migrate(suite.db.DB)
	switch {
	case errors.Is(err, migrate.ErrNoChange):
	default:
		require.NoError(suite.T(), err)
	}

	r := storage.NewStorage(suite.db, storePath)
	suite.srv = NewService(r)

	require.NoError(suite.T(), err)

}

func (suite *serviceStoreTestSuite) TearDownSuite() {
	err := suite.db.Close()
	require.NoError(suite.T(), err)
	require.NoError(suite.T(), os.RemoveAll(suite.T().TempDir()))
}

func TestHandlersFileStoreTest(t *testing.T) {
	suite.Run(t, new(serviceStoreTestSuite))
}

func (suite *serviceStoreTestSuite) Test_service() {
	t := suite.T()

	type wantErr struct {
		save, get, list, del bool
	}
	type args struct {
		list *model.ListQuery
		del  bool
	}
	tests := []struct {
		name     string
		dataItem model.Model
		args     args
		wantErr  wantErr
	}{
		{
			name: "test auth",
			dataItem: &auth.Model{
				Common: model.Common{Key: "test-set-auth-1"},
				Data: &auth.Data{
					Login:    "login1",
					Password: "password1",
				},
			},
			args: args{
				list: &model.ListQuery{},
				del:  true,
			},
			wantErr: wantErr{},
		},
		{
			name: "test card",
			dataItem: &card.Model{
				Common: model.Common{Key: "test-set-card-1"},
				Data: &card.Data{
					Exp:    "11/05",
					Number: "0000000000000000",
					Name:   "CardHolder Name",
					CVV:    "000",
				},
			},
			args: args{
				list: &model.ListQuery{},
				del:  true,
			},
			wantErr: wantErr{},
		},
		{
			name: "test Bin",
			dataItem: &bin.Model{
				Common: model.Common{Key: "test-set-Bin-1"},
				Data: &bin.Data{
					Bin: []byte("SOME BYTE SLICE"),
				},
			},
			args: args{
				list: &model.ListQuery{},
				del:  true,
			},
			wantErr: wantErr{},
		},
		{
			name: "test bin",
			dataItem: &bin.Model{
				Common: model.Common{Key: "test-set-Bin-1"},
				Data: &bin.Data{
					Bin: []byte("SOME BYTE SLICE"),
				},
			},
			args: args{
				list: &model.ListQuery{},
				del:  true,
			},
			wantErr: wantErr{},
		},
		{
			name: "test text",
			dataItem: &text.Model{
				Common: model.Common{Key: "test-set-Bin-1"},
				Data: &text.Data{
					Text: "some text\ntext some\nmultiline",
				},
			},
			args: args{
				list: &model.ListQuery{},
				del:  true,
			},
			wantErr: wantErr{},
		},
	}
	for _, tt := range tests {
		if tt.dataItem != nil {
			t.Run(tt.name+" save", func(t *testing.T) {
				err := suite.srv.Save(tt.dataItem)
				assert.Equal(t, tt.wantErr.save, (err != nil),
					fmt.Sprintf("SaveStore() error = %v, wantErr %v", err, tt.wantErr.save))
			})
		}

		if tt.args.list != nil {
			t.Run(tt.name+" list", func(t *testing.T) {
				gotData, err := suite.srv.List(*tt.args.list)
				if (err != nil) != tt.wantErr.list {
					t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr.list)
					return
				}
				if tt.dataItem != nil {
					assert.Greater(t, gotData.Total, 0)
					exist := false
					for _, i := range gotData.Items {
						if i.Key == tt.dataItem.GetKey() {
							exist = true
							break
						}
					}
					assert.True(t, exist, tt.dataItem.GetKey()+" not exist in list ", gotData)
				}
			})
		}

		if tt.dataItem != nil {
			t.Run(tt.name+" get saved", func(t *testing.T) {
				gotItemData, err := suite.srv.Get(tt.dataItem.GetKey())

				assert.Equal(t, tt.wantErr.get, (err != nil),
					fmt.Sprintf("GetStored() error = %v, wantErr %v", err, tt.wantErr.get))

				assert.NotNil(t, gotItemData)

				assert.Equal(t, gotItemData.Key, tt.dataItem.GetKey())
				assert.Equal(t, gotItemData.Description, tt.dataItem.GetDescription())

				if !reflect.DeepEqual(gotItemData.Data.GetData(), tt.dataItem.GetData()) {
					t.Errorf("GetStored() gotData = %v, want %v", gotItemData.Data, tt.dataItem.GetData())
				}
			})
		}

		if tt.args.del {
			t.Run(tt.name+" delete", func(t *testing.T) {
				err := suite.srv.Delete(tt.dataItem.GetKey())

				assert.Equal(t, tt.wantErr.del, (err != nil),
					fmt.Sprintf("Delete() error = %v, wantErr %v", err, tt.wantErr.del))

				_, err = suite.srv.Get(tt.dataItem.GetKey())
				if err == nil || !errors.Is(err, sql.ErrNoRows) {
					t.Errorf("Delete failed.  %s steel alife, error: %v ", tt.dataItem.GetKey(), err)
				}
			})
		}
	}
}
