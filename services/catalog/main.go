package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/Nulandmori/micorservices-pattern/pkg/env"
	"github.com/Nulandmori/micorservices-pattern/pkg/logger"
	"github.com/Nulandmori/micorservices-pattern/pkg/tracer"
	"github.com/Nulandmori/micorservices-pattern/services/catalog/grpc"
	"golang.org/x/sys/unix"
)

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	defaultPort := 8080
	port := env.GetPort(defaultPort)

	ctx, stop := signal.NotifyContext(ctx, unix.SIGTERM, unix.SIGINT)
	defer stop()

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
	clogger := l.WithName("catalog")

	errCh := make(chan error, 1)
	go func() {
		errCh <- grpc.RunServer(ctx, port, clogger.WithName("grpc"))
	}()
	select {
	case err := <-errCh:
		fmt.Println(err.Error())
		return 1
	case <-ctx.Done():
		fmt.Println("shutting down...")
		return 0
	}
}
