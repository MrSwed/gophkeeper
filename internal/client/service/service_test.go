package service

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"gophKeeper/internal/client/config"
	cfg "gophKeeper/internal/client/config"
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
	db       *sqlx.DB
	srv      Service
	oldStdin *os.File
	user     string
	userBak  string
	pass     string
}

var testDataPath string = filepath.Join("..", "..", "..", "testdata")

func (suite *serviceStoreTestSuite) SetupSuite() {
	suite.user = "test-" + time.Now().Format("20060102150405")

	storePath := filepath.Join(suite.T().TempDir(), config.AppName, suite.user)
	// err := os.MkdirAll(storePath, os.ModePerm)
	// require.NoError(suite.T(), err)
	dbFile := filepath.Join(storePath, "store.db")
	profiles := cfg.Glob.GetStringMap("profiles")
	profiles[suite.user] = cfg.NewGlobProfileItem(storePath)
	cfg.Glob.Set("profiles", profiles)

	suite.userBak = cfg.Glob.GetString("profile")
	cfg.Glob.Set("profile", suite.user)
	err := cfg.UserLoad()
	require.NoError(suite.T(), err)

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

	suite.oldStdin = os.Stdin
	suite.pass = "SomeUserPassword"
	input := []byte(strings.Join([]string{suite.pass, suite.pass, ""}, "\n"))
	rp, wp, err := os.Pipe()
	require.NoError(suite.T(), err)
	_, err = wp.Write(input)
	require.NoError(suite.T(), err)
	err = wp.Close()
	require.NoError(suite.T(), err)
	os.Stdin = rp
	_, err = suite.srv.GetToken()
	require.NoError(suite.T(), err)
}

func (suite *serviceStoreTestSuite) TearDownSuite() {
	err := suite.db.Close()
	profiles := cfg.Glob.GetStringMap("profiles")
	delete(profiles, suite.user)
	cfg.Glob.Set("profiles", profiles)
	cfg.Glob.Set("profile", suite.userBak)

	require.NoError(suite.T(), err)
	require.NoError(suite.T(), os.RemoveAll(suite.T().TempDir()))
	os.Stdin = suite.oldStdin
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
		save model.Model
		del  bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr wantErr
	}{
		{
			name: "test auth",
			args: args{
				save: &auth.Model{
					Common: model.Common{Key: "test-set-auth-1"},
					Data: &auth.Data{
						Login:    "login1",
						Password: "password1",
					},
				},
				list: &model.ListQuery{},
				del:  true,
			},
			wantErr: wantErr{},
		},
		{
			name: "test auth 2",
			args: args{
				save: &auth.Model{
					Common: model.Common{Key: "test-set-auth-1", Description: "Some description"},
					Data: &auth.Data{
						Login:    "login1",
						Password: "password1",
					},
				},
				list: &model.ListQuery{},
				del:  true,
			},
			wantErr: wantErr{},
		},
		{
			name: "test card",
			args: args{
				save: &card.Model{
					Common: model.Common{Key: "test-set-card-1"},
					Data: &card.Data{
						Exp:    "11/05",
						Number: "0000000000000000",
						Name:   "CardHolder Name",
						CVV:    "000",
					},
				},
				list: &model.ListQuery{},
				del:  true,
			},
			wantErr: wantErr{},
		},
		{
			name: "test Bin",
			args: args{
				save: &bin.Model{
					Common: model.Common{Key: "test-set-Bin-1"},
					Data: &bin.Data{
						Bin: []byte("SOME BYTE SLICE"),
					},
				},
				list: &model.ListQuery{},
				del:  true,
			},
			wantErr: wantErr{},
		},
		{
			name:    "test Bin File wrong path",
			wantErr: wantErr{save: true, get: true, del: true},
			args: args{
				save: &bin.Model{
					Common: model.Common{
						Key:      "test-set-Bin-2",
						FileName: filepath.Join("some-wrong-path", "SomeFile.pdf"),
					},
				},
			},
		},
		{
			name: "test Bin File",
			args: args{
				save: &bin.Model{
					Common: model.Common{
						Key:      "test-set-Bin-3",
						FileName: filepath.Join(testDataPath, "SomeFile.pdf"),
					},
				},
				list: &model.ListQuery{},
				del:  true,
			},
			wantErr: wantErr{},
		},
		{
			name: "test text",
			args: args{
				save: &text.Model{
					Common: model.Common{Key: "test-set-text-1"},
					Data: &text.Data{
						Text: "some text\ntext some\nmultiline",
					},
				},
				list: &model.ListQuery{},
				del:  true,
			},
			wantErr: wantErr{},
		},
		{
			name: "test text File",
			args: args{
				save: &text.Model{
					Common: model.Common{
						Key:      "test-set-text-3",
						FileName: filepath.Join(testDataPath, "SomeFile.txt"),
					},
				},
				list: &model.ListQuery{},
				del:  true,
			},
			wantErr: wantErr{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.save != nil {
				t.Run(tt.name+" save", func(t *testing.T) {
					err := suite.srv.Save(tt.args.save)
					require.Equal(t, tt.wantErr.save, (err != nil),
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
					if tt.args.save != nil {
						assert.Greater(t, gotData.Total, 0)
						exist := false
						for _, i := range gotData.Items {
							if i.Key == tt.args.save.GetKey() {
								exist = true
								break
							}
						}
						assert.True(t, exist, tt.args.save.GetKey()+" not exist in list ", gotData)
					}
				})
			}

			if tt.args.save != nil {
				t.Run(tt.name+" get saved", func(t *testing.T) {
					gotItemData, err := suite.srv.Get(tt.args.save.GetKey())

					assert.Equal(t, tt.wantErr.get, (err != nil),
						fmt.Sprintf("GetStored() error = %v, wantErr %v", err, tt.wantErr.get))

					if !tt.wantErr.get {
						assert.NotNil(t, gotItemData)

						assert.Equal(t, gotItemData.Key, tt.args.save.GetKey())
						assert.Equal(t, gotItemData.Description, tt.args.save.GetDescription())

						if !reflect.DeepEqual(gotItemData.Data.GetPacked(), tt.args.save.GetPacked()) {
							t.Errorf("GetStored() gotData = %v, want %v", gotItemData.Data, tt.args.save.GetPacked())
						}
					}
				})
			}

			if tt.args.del {
				t.Run(tt.name+" delete", func(t *testing.T) {
					err := suite.srv.Delete(tt.args.save.GetKey())

					assert.Equal(t, tt.wantErr.del, (err != nil),
						fmt.Sprintf("Delete() error = %v, wantErr %v", err, tt.wantErr.del))

					_, err = suite.srv.Get(tt.args.save.GetKey())
					if err == nil || !errors.Is(err, sql.ErrNoRows) {
						t.Errorf("Delete failed.  %s steel alife, error: %v ", tt.args.save.GetKey(), err)
					}
				})
			}
		})
	}

}

func (suite *serviceStoreTestSuite) Test_GetToken() {
	t := suite.T()
	t.Run("Test GetToken again", func(t *testing.T) {

		token, err := suite.srv.GetToken()
		require.NoError(t, err)

		// clear stored encryption_key for initialize request password from user
		cfg.User.Set("encryption_key", "")
		input := []byte(strings.Join([]string{suite.pass, ""}, "\n"))

		rp, wp, err := os.Pipe()
		require.NoError(suite.T(), err)
		_, err = wp.Write(input)
		require.NoError(suite.T(), err)
		err = wp.Close()
		require.NoError(suite.T(), err)
		os.Stdin = rp
		token2, err := suite.srv.GetToken()
		require.NoError(t, err)
		require.Equal(t, token, token2)
	})
}
