package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"gophKeeper/internal/server/closer"
	"gophKeeper/internal/server/config"
	"gophKeeper/internal/server/constant"

	hgrpc "gophKeeper/internal/server/handler/grpc"
	myMigrate "gophKeeper/internal/server/migrate"
	"gophKeeper/internal/server/repository"
	"gophKeeper/internal/server/service"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func buildInfo(s string) string {
	if s == "" {
		return "N/A"
	}
	return s
}

type BuildMetadata struct {
	Version string `json:"buildVersion"`
	Date    string `json:"buildDate"`
	Commit  string `json:"buildCommit"`
}

type App struct {
	cfg        *config.Config
	log        *zap.Logger
	db         *sqlx.DB
	closer     *closer.Closer
	grpc       *grpc.Server
	srv        service.Service
	http       *http.Server
	eg         *errgroup.Group
	stop       context.CancelFunc
	lockDB     chan struct{}
	isNewStore bool
}

func RunApp(ctx context.Context, cfg *config.Config, log *zap.Logger, buildData BuildMetadata) {
	var (
		err  error
		stop context.CancelFunc
	)
	ctx, stop = signal.NotifyContext(ctx, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	if cfg == nil {
		cfg, err = config.NewConfig().Init()
		if err != nil {
			panic(err)
		}
	}
	if log == nil {
		if cfg.Debug {
			log, err = zap.NewDevelopment()
		} else {
			log, err = zap.NewProduction()
		}
		if err != nil {
			panic(err)
		}
	}

	appHandler := newApp(ctx, stop, cfg, log)

	appHandler.log.Info("Init app", zap.Any(`Build info`, map[string]string{
		`Build version`: buildInfo(buildData.Version),
		`Build date`:    buildInfo(buildData.Date),
		`Build commit`:  buildInfo(buildData.Commit)}))

	appHandler.Run(ctx)
	appHandler.Stop()
}

func newApp(ctx context.Context, stop context.CancelFunc, cfg *config.Config, log *zap.Logger) *App {
	eg, ctx := errgroup.WithContext(ctx)
	a := App{
		stop:       stop,
		cfg:        cfg,
		eg:         eg,
		isNewStore: true,
		closer:     &closer.Closer{},
		log:        log,
		lockDB:     make(chan struct{}),
	}

	a.maybeConnectDB(ctx)
	repo := repository.NewRepository(&a.cfg.StorageConfig, a.db)
	a.srv = service.New(repo, a.cfg)
	g := hgrpc.NewServer(a.srv, a.cfg, a.log)

	if a.cfg.GRPCAddress != "" {
		a.grpc = g.Handler()
	}

	return &a
}

func (a *App) maybeConnectDB(ctx context.Context) {
	if len(a.cfg.DatabaseDSN) == 0 {
		a.log.Fatal("database config is empty")
	}
	var err error
	if a.db, err = sqlx.ConnectContext(ctx, "postgres", a.cfg.DatabaseDSN); err != nil {
		a.log.Fatal("cannot connect db", zap.Error(err))
	}
	a.isNewStore = false
	a.log.Info("DB connected")
	versions, errM := myMigrate.Migrate(a.db.DB)
	switch {
	case errors.Is(errM, migrate.ErrNoChange):
		a.log.Info("DB migrate: ", zap.Any("info", errM), zap.Any("versions", versions))
	case errM == nil:
		a.log.Info("DB migrate: new applied ", zap.Any("versions", versions))
		a.isNewStore = versions[0] == 0
	default:
		a.log.Fatal("DB migrate: ", zap.Any("versions", versions), zap.Error(errM))
	}

}

func (a *App) shutdownDBStore(_ context.Context) (err error) {
	if a.db != nil {
		<-a.lockDB
		if err = a.db.Close(); err == nil {
			a.log.Info("Db Closed")
		}
	}
	return
}

func (a *App) grpcShutdown(_ context.Context) error {
	if a.grpc != nil {
		a.grpc.GracefulStop()
	}
	return nil
}

func (a *App) Run(ctx context.Context) {
	a.log.Info("Start server", zap.Any("Config", a.cfg))

	a.closer.Add("grpc", a.grpcShutdown)

	// Start HTTP server if HTTPAddress is set
	if a.cfg.HTTPAddress != "" {
		a.http = &http.Server{
			Addr:         a.cfg.HTTPAddress,
			Handler:      nil, // Use default handler
			ReadTimeout:  a.cfg.ReadTimeout,
			WriteTimeout: a.cfg.WriteTimeout,
			IdleTimeout:  a.cfg.IdleTimeout,
		}
		a.closer.Add("WEB", a.http.Shutdown)
		go func() {
			a.log.Info("Starting HTTP server", zap.String("address", a.cfg.HTTPAddress))
			if err := a.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				a.log.Error("HTTP server error", zap.Error(err))
				a.stop()
			}
		}()
	}

	close(a.lockDB)

	if a.db != nil {
		a.closer.Add("DB Close", a.shutdownDBStore)
	}

	go func() {
		listen, err := net.Listen("tcp", a.cfg.GRPCAddress)
		if err != nil {
			a.log.Error("grpc server", zap.Error(err))
		}
		if err = a.grpc.Serve(listen); err != nil {
			a.log.Error("grpc server", zap.Error(err))
			a.stop()
		}
	}()
	a.log.Info("grpc server started")
	<-ctx.Done()
}

func (a *App) Stop() {
	a.log.Info("Shutting down server gracefully")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), constant.ServerShutdownTimeout*time.Second)
	defer cancel()

	if err := a.closer.Close(shutdownCtx); err != nil {
		a.log.Error("Shutdown", zap.Error(err), zap.Any("timeout: ", constant.ServerShutdownTimeout))
	}

	a.log.Info("Server stopped")
}
