package app

import (
	"context"
	"fmt"
	pb "gophKeeper/internal/proto"
	"gophKeeper/internal/server/constant"
	errs "gophKeeper/internal/server/errors"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const waitPortInterval = 100 * time.Millisecond
const waitPortConnTimeout = 50 * time.Millisecond

type AppTestSuite struct {
	suite.Suite
	ctx     context.Context
	stop    context.CancelFunc
	address string
	pgCont  *postgres.PostgresContainer
	osArgs  []string
}

func createPostgresContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	pgContainer, err := postgres.Run(ctx, "postgres:14-alpine",
		// postgres.WithInitScripts(
		// 	filepath.Join("../../../", "testdata", "server.sql"),
		// ),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		return nil, err
	}

	return pgContainer, nil
}

func waitGRPCPort(ctx context.Context, address string) error {
	if address == "" {
		return nil
	}
	ticker := time.NewTicker(waitPortInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			conn, _ := net.DialTimeout("tcp", address, waitPortConnTimeout)
			if conn != nil {
				_ = conn.Close()
				return nil
			}
		}
	}
}

func testGRPCDial(addr string, ctx context.Context, meta map[string]string) (
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

func (suite *AppTestSuite) SetupSuite() {
	var (
		err error
	)
	suite.osArgs = os.Args
	os.Args = os.Args[0:1]
	suite.ctx, suite.stop = context.WithCancel(context.Background())
	// suite.cfg = config.NewConfig()

	suite.pgCont, err = createPostgresContainer(suite.ctx)
	require.NoError(suite.T(), err)
	databaseDSN, err := suite.pgCont.ConnectionString(suite.ctx, "sslmode=disable")
	require.NoError(suite.T(), err)
	suite.address = net.JoinHostPort("", fmt.Sprintf("%d", rand.Intn(200)+30000))

	suite.T().Setenv("DATABASE_DSN", databaseDSN)
	suite.T().Setenv("GRPC_ADDRESS", suite.address)
	suite.T().Setenv("GRPC_OPERATION_TIMEOUT", "5000s")

	go RunApp(suite.ctx, nil, nil, BuildMetadata{Version: "testing..", Date: time.Now().String(), Commit: ""})
	require.NoError(suite.T(), waitGRPCPort(suite.ctx, suite.address))

	db, err := sqlx.Connect("postgres", databaseDSN)
	predefined, err := os.ReadFile(filepath.Join("../../../", "testdata", "server.sql"))
	require.NoError(suite.T(), err)
	_, err = db.Exec(string(predefined))
	require.NoError(suite.T(), err)
}

func (suite *AppTestSuite) TearDownSuite() {
	if err := suite.pgCont.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
	suite.stop()
	os.Args = suite.osArgs
}

func TestApp(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}

func (suite *AppTestSuite) TestRegisterClient() {
	t := suite.T()

	tests := []struct {
		name    string
		req     *pb.RegisterClientRequest
		headers map[string]string
		wantErr []string
	}{
		{
			name: "success register",
			req: &pb.RegisterClientRequest{
				Email:    "test1@email.ru",
				Password: "Ansddd12@!",
			},
		},
		{
			name: "not valid password",
			req: &pb.RegisterClientRequest{
				Email:    "test2@email.ru",
				Password: "11111",
			},
			wantErr: []string{"Error:Field validation for 'Password'"},
		},
		{
			name: "not valid email",
			req: &pb.RegisterClientRequest{
				Email:    "test2-email.ru",
				Password: "Ansddd12@!",
			},
			wantErr: []string{"Error:Field validation for 'Email'"},
		},
		{
			name: "empty password",
			req: &pb.RegisterClientRequest{
				Email: "test2@email.ru",
			},
			wantErr: []string{"Error:Field validation for 'Password'", "required"},
		},
		{
			name: "exist user, wrong password",
			req: &pb.RegisterClientRequest{
				Email:    "example1@example.com",
				Password: "Ansddd12@!###dddsdf",
			},
			wantErr: []string{errs.ErrorWrongAuth.Error()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, conn, callOpt, err := testGRPCDial(suite.address, ctx, tt.headers)
			require.NoError(t, err)
			defer func() { require.NoError(t, conn.Close()) }()
			client := pb.NewAuthClient(conn)
			_, err = client.RegisterClient(ctx, tt.req, callOpt...)
			if tt.wantErr != nil {
				require.Error(t, err)
				for _, e := range tt.wantErr {
					require.Contains(t, err.Error(), e)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func (suite *AppTestSuite) TestSyncUser() {
	t := suite.T()

	existCreatedAt, err := time.Parse(time.RFC3339, "2024-09-17T12:00:00+03:00")
	require.NoError(t, err)
	existUpdatedAt, err := time.Parse(time.RFC3339, "2024-09-17T12:50:00+03:00")
	require.NoError(t, err)
	newCreatedAt := time.Now()
	tests := []struct {
		name     string
		req      *pb.UserSync
		wantResp *pb.UserSync
		headers  map[string]string
		wantErr  []string
	}{
		{
			name: "no token",
			req: &pb.UserSync{
				Email: "example@example.com",
			},
			wantErr: []string{errs.ErrorNoToken.Error()},
		},
		{
			name: "not valid token",
			req: &pb.UserSync{
				Email: "example@example.com",
			},
			headers: map[string]string{
				constant.TokenKey: "not valid token",
			},
			wantErr: []string{errs.ErrorInvalidToken.Error()},
		},
		{
			name: "new client, new server",
			req: &pb.UserSync{
				Email: "example1@example.com",
			},
			wantResp: &pb.UserSync{
				Email:     "example1@example.com",
				CreatedAt: timestamppb.New(existCreatedAt),
			},
			headers: map[string]string{
				constant.TokenKey: "8ca0c5a18320fc2f264cfa95639ea27888727c6090d6f9cb0d6c5798a93fcb63",
			},
			wantErr: nil,
		},
		{
			name: "send user",
			req: &pb.UserSync{
				Email:     "example2@example.com",
				CreatedAt: timestamppb.New(existCreatedAt),
				UpdatedAt: timestamppb.New(newCreatedAt),
				PackedKey: []byte("PackedKey"),
			},
			wantResp: &pb.UserSync{
				Email:     "example2@example.com",
				CreatedAt: timestamppb.New(existCreatedAt),
				UpdatedAt: timestamppb.New(newCreatedAt),
				PackedKey: []byte("PackedKey"),
			},
			headers: map[string]string{
				constant.TokenKey: "862AB376DF9DBD090F28F9DD9A2F5F1C9F88F05D27B63AE3942B5057C6BA2688",
			},
			wantErr: nil,
		},
		{
			name: "get user",
			req: &pb.UserSync{
				Email:       "example3@example.com",
				UpdatedAt:   timestamppb.New(existUpdatedAt.Add(-1 * time.Hour)),
				PackedKey:   []byte("some existed packed data"),
				Description: "description",
			},
			wantResp: &pb.UserSync{
				Email:     "example3@example.com",
				CreatedAt: timestamppb.New(existCreatedAt),
				UpdatedAt: timestamppb.New(existUpdatedAt),
				PackedKey: []byte("predefined packed data"),
			},
			headers: map[string]string{
				constant.TokenKey: "C4B7F91016F52C039804D05E61C67A87A51BB8CD78FF04E51AB769ED8336D77E",
			},
			wantErr: nil,
		},
		{
			name: "bad sync key",
			req: &pb.UserSync{
				Description: "description",
			},

			headers: map[string]string{
				constant.TokenKey: "C4B7F91016F52C039804D05E61C67A87A51BB8CD78FF04E51AB769ED8336D77E",
			},
			wantErr: []string{errs.ErrorSyncNoKey.Error()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, conn, callOpt, err := testGRPCDial(suite.address, ctx, tt.headers)
			require.NoError(t, err)
			defer func() { require.NoError(t, conn.Close()) }()
			client := pb.NewUserClient(conn)
			data, err := client.SyncUser(ctx, tt.req, callOpt...)
			if tt.wantErr != nil {
				for _, e := range tt.wantErr {
					require.Contains(t, err.Error(), e)
				}
			} else {
				require.NoError(t, err)
			}
			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp.Email, data.Email, "Email")
				assert.Equal(t, tt.wantResp.CreatedAt.AsTime(), data.CreatedAt.AsTime(), "CreatedAt")
				if tt.wantResp.UpdatedAt != nil {
					assert.Equal(t, tt.wantResp.UpdatedAt.AsTime(), data.UpdatedAt.AsTime(), "UpdatedAt")
				}
				assert.Equal(t, tt.wantResp.Description, data.Description, "Description")
				assert.Equal(t, tt.wantResp.PackedKey, data.PackedKey, "PackedKey")
			}
		})
	}
}
