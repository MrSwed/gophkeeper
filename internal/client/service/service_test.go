package service

import (
	"database/sql"
	"errors"
	"fmt"
	errs "gophKeeper/internal/client/errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

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
	db                         *sqlx.DB
	srv                        Service
	oldStdin, stdin, stdinPipe *os.File
	user                       string
	userBak                    string
	pass                       string
}

var testDataPath string = filepath.Join("..", "..", "..", "testdata")

func (s *serviceStoreTestSuite) SetupSuite() {
	s.user = "test-" + time.Now().Format("20060102150405")

	storePath := filepath.Join(s.T().TempDir(), cfg.AppName, s.user)
	dbFile := filepath.Join(storePath, "store.db")
	profiles := cfg.Glob.GetStringMap("profiles")
	profiles[s.user] = cfg.NewGlobProfileItem(storePath)
	cfg.Glob.Set("profiles", profiles)

	s.userBak = cfg.Glob.GetString("profile")
	cfg.Glob.Set("profile", s.user)
	err := cfg.UserLoad()
	require.NoError(s.T(), err)

	s.db, err = sqlx.Open("sqlite3", dbFile)
	require.NoError(s.T(), err)

	_, err = clMigrate.Migrate(s.db.DB)
	switch {
	case errors.Is(err, migrate.ErrNoChange):
	default:
		require.NoError(s.T(), err)
	}

	r := storage.NewStorage(s.db, storePath)
	s.srv = NewService(r)

	s.stdin, s.stdinPipe, err = os.Pipe()
	require.NoError(s.T(), err)
	s.oldStdin, os.Stdin = os.Stdin, s.stdin
	s.pass = "SomeUserPassword"
	s.input(s.pass, s.pass)
	_, err = s.srv.GetToken()
	require.NoError(s.T(), err)
}

func (s *serviceStoreTestSuite) input(str ...string) {
	input := []byte(strings.Join(str, "\n") + "\n")
	_, err := s.stdinPipe.Write(input)
	require.NoError(s.T(), err)
}

func (s *serviceStoreTestSuite) TearDownSuite() {
	err := s.db.Close()
	require.NoError(s.T(), err)
	// restore config
	profiles := cfg.Glob.GetStringMap("profiles")
	delete(profiles, s.user)
	cfg.Glob.Set("profiles", profiles)
	cfg.Glob.Set("profile", s.userBak)

	require.NoError(s.T(), os.RemoveAll(s.T().TempDir()))

	// restore stdin
	os.Stdin = s.oldStdin
	err = s.stdinPipe.Close()
	require.NoError(s.T(), err)
	err = s.stdin.Close()
	require.NoError(s.T(), err)
}

func TestService(t *testing.T) {
	suite.Run(t, new(serviceStoreTestSuite))
}

func (s *serviceStoreTestSuite) Test_service() {
	t := s.T()
	s.input(s.pass, s.pass)

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
					err := s.srv.Save(tt.args.save)
					require.Equal(t, tt.wantErr.save, (err != nil),
						fmt.Sprintf("SaveStore() error = %v, wantErr %v", err, tt.wantErr.save))
				})
			}

			if tt.args.list != nil {
				t.Run(tt.name+" list", func(t *testing.T) {
					gotData, err := s.srv.List(*tt.args.list)
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
					gotItemData, err := s.srv.Get(tt.args.save.GetKey())

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
					err := s.srv.Delete(tt.args.save.GetKey())

					assert.Equal(t, tt.wantErr.del, (err != nil),
						fmt.Sprintf("Delete() error = %v, wantErr %v", err, tt.wantErr.del))

					_, err = s.srv.Get(tt.args.save.GetKey())
					if err == nil || !errors.Is(err, sql.ErrNoRows) {
						t.Errorf("Delete failed.  %s steel alife, error: %v ", tt.args.save.GetKey(), err)
					}
				})
			}
		})
	}

}

func (s *serviceStoreTestSuite) Test_GetToken() {
	t := s.T()
	t.Run("Test GetToken again", func(t *testing.T) {

		token, err := s.srv.GetToken()
		require.NoError(t, err)

		// clear stored encryption_key for initialize request password from user
		cfg.User.Set("encryption_key", "")
		// auth again
		s.input(s.pass)
		token2, err := s.srv.GetToken()
		require.NoError(t, err)
		require.Equal(t, token, token2)
	})

}

func (s *serviceStoreTestSuite) Test_WrongPass() {
	t := s.T()
	t.Run("Test wrong pass", func(t *testing.T) {
		token, err := s.srv.GetToken()
		defer cfg.User.Set("encryption_key", token)
		require.NoError(t, err, "check current token error ")
		testKey := "testKeyForWrongPass"
		err = s.srv.Save(
			&text.Model{
				Common: model.Common{
					Key: testKey,
				},
				Data: &text.Data{
					Text: "some text here",
				},
			},
		)
		require.NoError(t, err, "save test data error ")

		_, err = s.srv.Get(testKey)
		require.NoError(t, err, "Get test data error ")

		// clear stored encryption_key for initialize request password from user
		cfg.User.Set("encryption_key", "")
		s.input("someWrongPass")

		cfg.Glob.Set("debug", false)
		_, err = s.srv.Get(testKey)
		require.Equal(t, true, errors.Is(err, errs.ErrPassword), fmt.Sprintf("Actual: %v. Expected: %v", err, errs.ErrPassword))

		cfg.Glob.Set("debug", true)
		s.input("someWrongPass")
		_, err = s.srv.Get(testKey)
		require.Contains(t, err.Error(), `unpad error`)
	})
}

func (s *serviceStoreTestSuite) Test_WrongPassConfirm() {
	t := s.T()
	t.Run("Test wrong pass confirm", func(t *testing.T) {
		token := cfg.User.GetString("packed_key")
		defer cfg.User.Set("encryption_key", token)
		packed := cfg.User.GetString("packed_key")
		defer cfg.User.Set("packed_key", packed)

		cfg.User.Set("packed_key", "")
		cfg.User.Set("encryption_key", "")
		// auth again
		s.input(s.pass, "wrongConfirm", "", "", "")
		_, err := s.srv.GetToken()
		require.Equal(s.T(), true, errors.Is(err, errs.ErrPasswordConfirm))
	})
}
