package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/Nulandmori/micorservices-pattern/pkg/env"
	"github.com/Nulandmori/micorservices-pattern/pkg/firebase"
	pkggrpc "github.com/Nulandmori/micorservices-pattern/pkg/grpc"
	"github.com/Nulandmori/micorservices-pattern/pkg/grpc/client/interceptor"
	catalog "github.com/Nulandmori/micorservices-pattern/services/catalog/proto"
	"github.com/Nulandmori/micorservices-pattern/services/gateway/proto"
	"github.com/go-logr/logr"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultTLSPort = "443"
	caudience      = "https://catalog-service-y64oiofbkq-an.a.run.app"
)

func RunServer(ctx context.Context, port int, logger logr.Logger) error {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	}

	catalogServiceAddr := env.MustGetEnv("CATALOG_SERVICE_ADDR")

	if strings.Contains(catalogServiceAddr, defaultTLSPort) {
		creds := credentials.NewTLS(&tls.Config{})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	copts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			interceptor.AuthServiceUnnaryClientInterceptor(caudience),
		),
	}

	copts = append(copts, opts...)

	cconn, err := grpc.DialContext(ctx, catalogServiceAddr, copts...)
	if err != nil {
		return fmt.Errorf("failed to dial catalog grpc server: %w", err)
	}

	firebaseAuth, err := firebase.NewAuthClient()
	if err != nil {
		return fmt.Errorf("failed to initialize firebase auth client: %w", err)
	}

	svc := &server{
		catalogClient: catalog.NewCatalogServiceClient(cconn),
		firebaseAuth:  firebaseAuth,
		logger:        logger,
	}
	return pkggrpc.NewServer(port, logger, func(s *grpc.Server) {
		proto.RegisterGatewayServiceServer(s, svc)
	}).Start(ctx)
}
