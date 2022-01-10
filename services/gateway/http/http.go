package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Nulandmori/micorservices-pattern/services/gateway/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

func RunServer(ctx context.Context, port int, grpcPort int) error {
	// https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/customizing_your_gateway/
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: false,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	}
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	conn, err := grpc.DialContext(ctx, fmt.Sprintf(":%d", grpcPort), opts...)
	if err != nil {
		return fmt.Errorf("failed to dial grpc server: %w", err)
	}

	if err := proto.RegisterGatewayServiceHandlerClient(ctx, mux, proto.NewGatewayServiceClient(conn)); err != nil {
		return fmt.Errorf("failed to create a grpc client: %w", err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("failed to serve http server: %w", err)
	case <-ctx.Done():
		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown http server: %w", err)
		}
		if err := <-errCh; err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to serve http server: %w", err)
		}
		return nil
	}

}
