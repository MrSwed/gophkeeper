package app

import (
	"context"
	"fmt"
	pb "gophKeeper/internal/proto"
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
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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
				Email:    "example@example.com",
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
				for _, e := range tt.wantErr {
					require.Contains(t, err.Error(), e)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
