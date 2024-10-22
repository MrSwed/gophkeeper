/*
This package provides the implementation of the gRPC server for the GophKeeper application.
It defines the server handler and interceptors for managing authentication and data services.

Main functionalities include:

- Creating a new gRPC server with logging and authentication interceptors.
- Handling user authentication and managing data services.
*/
package grpc

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "gophKeeper/internal/proto"
	"gophKeeper/internal/server/config"
	"gophKeeper/internal/server/constant"
	errs "gophKeeper/internal/server/errors"
	"gophKeeper/internal/server/service"
)

// handler holds the configuration and logger for the gRPC server.
type handler struct {
	log       *zap.Logger
	c         *config.Config
	skipToken []string
}

// Handler represents the gRPC server handler.
type Handler struct {
	s service.Service
	handler
}

// NewServer creates a new instance of the gRPC server handler.
func NewServer(s service.Service, c *config.Config, log *zap.Logger) *Handler {
	return &Handler{
		s: s,
		handler: handler{
			log:       log,
			c:         c,
			skipToken: []string{"/service.Auth/RegisterClient"},
		},
	}
}

// Handler initializes and returns a new gRPC server with interceptors.
func (h *Handler) Handler() (s *grpc.Server) {
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.FinishCall),
	}
	s = grpc.NewServer(grpc.ChainUnaryInterceptor(
		logging.UnaryServerInterceptor(h.interceptorLogger(h.log), opts...),
		h.unaryInterceptor,
	))
	pb.RegisterDataServer(s, NewDataServer(h.s, h.c, h.log))
	pb.RegisterAuthServer(s, NewAuthServer(h.s, h.c, h.log))
	pb.RegisterUserServer(s, NewUserServer(h.s, h.c, h.log))
	return
}

// skipMethod checks if a method should be skipped for authentication.
func (h *Handler) skipMethod(method string) bool {
	for _, skipMethod := range h.skipToken {
		if method == skipMethod {
			return true
		}
	}
	return false
}

// unaryInterceptor intercepts unary RPC calls for authentication.
func (h *Handler) unaryInterceptor(ctx context.Context, req interface{}, g *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if !h.skipMethod(g.FullMethod) {
		var token []byte
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if values := md.Get(pb.TokenKey); len(values) > 0 {
				token = []byte(values[0])
			}
		}
		if len(token) == 0 {
			return nil, status.Error(codes.Unauthenticated, errs.ErrorNoToken.Error())
		}
		userID, err := h.s.UserIDByToken(ctx, token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, errs.ErrorInvalidToken.Error())
		}
		ctx = context.WithValue(ctx, constant.CtxUserID, userID)
	}
	return handler(ctx, req)
}

// interceptorLogger adapts zap logger to interceptor logger.
func (h *Handler) interceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)
		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]
			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}
		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)
		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
