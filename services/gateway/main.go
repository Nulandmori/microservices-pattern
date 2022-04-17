package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Nulandmori/micorservices-pattern/pkg/env"
	"github.com/Nulandmori/micorservices-pattern/pkg/logger"
	"github.com/Nulandmori/micorservices-pattern/pkg/run"
	"github.com/Nulandmori/micorservices-pattern/pkg/tracer"
	"github.com/Nulandmori/micorservices-pattern/services/gateway/grpc"
	"github.com/Nulandmori/micorservices-pattern/services/gateway/http"
)

// run server main
func main() {
	run.Run(server)
}

func server(ctx context.Context) int {
	grpcPort := 9090
	defaultPort := 8080
	httpPort := env.GetPort(defaultPort)

    cleanup := tracer.InitTracer()
    defer cleanup(ctx)

	l, err := logger.New()
	if err != nil {
		_, ferr := fmt.Fprintf(os.Stderr, "failed to create logger: %s", err)
		if ferr != nil {
			panic(fmt.Sprintf("failed to write log:`%s` original error is:`%s`", ferr, err))
		}
		return 1
	}
	glogger := l.WithName("gateway")

	grpcErrCh := make(chan error, 1)
	go func() {
		grpcErrCh <- grpc.RunServer(ctx, grpcPort, glogger.WithName("grpc"))
	}()

	httpErrCh := make(chan error, 1)
	go func() {
		httpErrCh <- http.RunServer(ctx, httpPort, grpcPort)
	}()

	select {
	case err := <-grpcErrCh:
		glogger.Error(err, "failed to serve grpc server")
		return 1
	case err := <-httpErrCh:
		glogger.Error(err, "failed to serve http server")
		return 1
	case <-ctx.Done():
		glogger.Info("shutting down...")
		return 0
	}
}
