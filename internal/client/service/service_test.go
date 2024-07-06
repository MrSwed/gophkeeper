package service

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"gophKeeper/internal/client/config"
	clMigrate "gophKeeper/internal/client/migrate"
	"gophKeeper/internal/client/model/input"
	"gophKeeper/internal/client/model/out"
	"gophKeeper/internal/client/storage"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/mattn/go-sqlite3"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const testSQLdata = `
insert into storage (key, description, created_at, updated_at, filename)
values  ('testKeyForGet', 'Some description', '2024-07-06 20:00:39', '2024-07-06 20:00:39', '20240706200039-testKeyForGet'),
        ('KeyForDelete', null, '2024-07-06 20:02:20', '2024-07-06 20:02:20', '20240706200220-KeyForDelete'),
        ('testkey2', null, '2024-07-06 20:02:20', '2024-07-06 20:02:20', '20240706200220-testKey2'),
        ('testkey3', null, '2024-07-06 20:02:21', '2024-07-06 20:02:21', '20240706200221-testKey3');
`

type serviceStoreTestSuite struct {
	suite.Suite
	db  *sqlx.DB
	srv Service
}

func parseTime(t *testing.T, s string) (tm time.Time) {
	var e error
	tm, e = time.Parse(time.DateTime, s)
	require.NoError(t, e)
	return

}

func (suite *serviceStoreTestSuite) SetupSuite() {
	storePath := filepath.Join(suite.T().TempDir(), config.AppName)
	err := os.MkdirAll(storePath, os.ModePerm)
	require.NoError(suite.T(), err)
	dbFile := filepath.Join(storePath, "store.db")

	suite.db, err = sqlx.Open("sqlite3", dbFile)
	require.NoError(suite.T(), err)

	_, err = clMigrate.Migrate(suite.db.DB)
	switch {
	case errors.Is(err, migrate.ErrNoChange):
	default:
		require.NoError(suite.T(), err)
	}

	_, err = suite.db.Exec(testSQLdata)

	require.NoError(suite.T(), err)

	r := storage.NewStorage(suite.db, storePath)
	suite.srv = NewService(r)
}

func (suite *serviceStoreTestSuite) TearDownSuite() {
	err := suite.db.Close()
	require.NoError(suite.T(), err)
	require.NoError(suite.T(), os.RemoveAll(suite.T().TempDir()))
}

func TestHandlersFileStoreTest(t *testing.T) {
	suite.Run(t, new(serviceStoreTestSuite))
}

func (suite *serviceStoreTestSuite) Test_service_Delete() {
	t := suite.T()
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test 1",
			args:    args{key: "KeyForDelete"},
			wantErr: false,
		},
		{
			name:    "test 2",
			args:    args{key: "someUnknownKey"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := suite.srv.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (suite *serviceStoreTestSuite) Test_service_Get() {
	t := suite.T()

	type args struct {
		key string
	}
	tests := []struct {
		name     string
		args     args
		wantData out.Item
		wantErr  bool
	}{
		{
			name: "test 1",
			args: args{key: "testKeyForGet"},
			wantData: out.Item{
				DBItem: storage.DBItem{
					Key:         "testKeyForGet",
					Description: &[]string{"Some description"}[0],
					CreatedAt:   parseTime(t, `2024-07-06 20:00:39`).Format(time.DateTime),
					UpdatedAt:   parseTime(t, `2024-07-06 20:00:39`).Format(time.DateTime),
				},
				Data: nil, // / todo
			},
			wantErr: false,
		},
		{
			name:    "test 2",
			args:    args{key: "someUnknown"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := suite.srv.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("Get() gotData = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func (suite *serviceStoreTestSuite) Test_service_List() {
	t := suite.T()

	tests := []struct {
		name     string
		query    input.ListQuery
		wantData out.List
		wantErr  bool
	}{
		{
			name: "query test",
			query: input.ListQuery{
				Key: "test",
			},
			wantData: out.List{
				Items: []storage.ListItem{
					{
						DBItem: storage.DBItem{
							Key:         "testKeyForGet",
							Description: &[]string{"Some description"}[0],
							CreatedAt:   parseTime(t, "2024-07-06 20:00:39").Format(time.DateTime),
							UpdatedAt:   parseTime(t, "2024-07-06 20:00:39").Format(time.DateTime),
						},
					},
					{
						DBItem: storage.DBItem{
							Key:         "testkey2",
							Description: nil,
							CreatedAt:   parseTime(t, "2024-07-06 20:02:20").Format(time.DateTime),
							UpdatedAt:   parseTime(t, "2024-07-06 20:02:20").Format(time.DateTime),
						},
					},
					{
						DBItem: storage.DBItem{
							Key:         "testkey3",
							Description: nil,
							CreatedAt:   parseTime(t, "2024-07-06 20:02:21").Format(time.DateTime),
							UpdatedAt:   parseTime(t, "2024-07-06 20:02:21").Format(time.DateTime),
						},
					},
				},
				Total: 3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotData, err := suite.srv.List(tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("List() gotData = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func (suite *serviceStoreTestSuite) Test_service_Set() {
	t := suite.T()

	tests := []struct {
		name    string
		args    input.Model
		wantErr bool
	}{
		{
			name: "test set auth 1",
			args: &input.Auth{
				Common:   input.Common{Key: "test-set-auth-1"},
				Login:    "login1",
				Password: "password1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := suite.srv.Set(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
