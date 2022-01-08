package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/Nulandmori/micorservices-pattern/pkg/grpc/interceptor"
	"github.com/go-logr/logr"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	channelz "google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/reflection"
)

var defaultNOPAuthFunc = func(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

type Server struct {
	server *grpc.Server
	port   int
}

func NewServer(port int, logger logr.Logger, register func(server *grpc.Server)) *Server {
	interceptors := []grpc.UnaryServerInterceptor{
		interceptor.NewRequestLogger(logger.WithName("request")),
		interceptor.NewAuthTokenPropagator(),
		grpc_auth.UnaryServerInterceptor(defaultNOPAuthFunc),
	}

	opts := []grpc.ServerOption{
		middleware.WithUnaryServerChain(interceptors...),
	}

	server := grpc.NewServer(opts...)

	register(server)

	reflection.Register(server)
	channelz.RegisterChannelzServiceToServer(server)

	return &Server{
		server: server,
		port:   port,
	}
}

func (s *Server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen on %d: %v", s.port, err)
	}

	errCh := make(chan error, 1)

	go func() {
		if err := s.server.Serve(listener); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("server has stopped with error: %v", err)
		}
		return nil
	case <-ctx.Done():
		s.server.GracefulStop()
		return <-errCh
	}
}
