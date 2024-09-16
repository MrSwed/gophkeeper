package app

import (
	"context"
	"fmt"
	pb "gophKeeper/internal/proto"
	"log"
	"math/rand"
	"net"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
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

/*
func testGRPCProto(suite AppTestSuite) {
	t := suite.T()
	t.Run("generated grpc proto not implemented", func(t *testing.T) {
		g := myGrpc.NewMetricsServer(suite.Srv(), suite.Cfg(), zap.NewNop())
		ctx := context.Background()
		_, err := g.UnimplementedMetricsServer.GetMetrics(ctx, nil)
		assert.Error(t, err, status.Errorf(codes.Unimplemented, "method GetMetrics not implemented"))
		_, err = g.UnimplementedMetricsServer.GetMetric(ctx, nil)
		assert.Error(t, err, status.Errorf(codes.Unimplemented, "method GetMetric not implemented"))
		_, err = g.UnimplementedMetricsServer.SetMetrics(ctx, nil)
		assert.Error(t, err, status.Errorf(codes.Unimplemented, "method SetMetrics not implemented"))
		_, err = g.UnimplementedMetricsServer.SetMetric(ctx, nil)
		assert.Error(t, err, status.Errorf(codes.Unimplemented, "method SetMetric not implemented"))
	})
}
*/

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
		name     string
		req      *pb.RegisterClientRequest
		wantResp *pb.ClientToken
		headers  map[string]string
		wantErr  error
	}{
		{
			name: "success register",
			req: &pb.RegisterClientRequest{
				Email:    "test@email.ru",
				Password: "11111",
			},
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
			} else {
				require.NoError(t, err)
			}
		})
	}
}
