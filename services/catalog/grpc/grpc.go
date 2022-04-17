package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"

	"github.com/Nulandmori/micorservices-pattern/pkg/env"
	pkggrpc "github.com/Nulandmori/micorservices-pattern/pkg/grpc"
	"github.com/Nulandmori/micorservices-pattern/pkg/grpc/client/interceptor"
	"github.com/Nulandmori/micorservices-pattern/services/catalog/proto"
	customer "github.com/Nulandmori/micorservices-pattern/services/customer/proto"
	item "github.com/Nulandmori/micorservices-pattern/services/item/proto"
	"github.com/go-logr/logr"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultTLSPort = "443"
	iaudience      = "https://item-service-y64oiofbkq-an.a.run.app"
	caudience      = "https://customer-service-y64oiofbkq-an.a.run.app"
)

func RunServer(ctx context.Context, port int, logger logr.Logger) error {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	}

	itemServiceAddr := env.MustGetEnv("ITEM_SERVICE_ADDR")
	customerServiceAddr := env.MustGetEnv("CUSTOMER_SERVICE_ADDR")

	if strings.Contains(itemServiceAddr, defaultTLSPort) && strings.Contains(customerServiceAddr, defaultTLSPort) {
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			return fmt.Errorf("failed to get cert pool: %w", err)
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	iopts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(

			interceptor.AuthServiceUnnaryClientInterceptor(iaudience),
		),
	}
	iopts = append(iopts, opts...)

	copts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(

			interceptor.AuthServiceUnnaryClientInterceptor(caudience),
		),
	}
	copts = append(copts, opts...)

	iconn, err := grpc.DialContext(ctx, itemServiceAddr, iopts...)
	if err != nil {
		return fmt.Errorf("failed to dial item grpc server: %w", err)
	}

	cconn, err := grpc.DialContext(ctx, customerServiceAddr, copts...)
	if err != nil {
		return fmt.Errorf("failed to dial customer grpc server: %w", err)
	}

	itemClient := item.NewItemServiceClient(iconn)
	customerClient := customer.NewCustomerServiceClient(cconn)

	svc := &server{
		itemClient:     itemClient,
		customerClient: customerClient,
	}
	return pkggrpc.NewServer(port, logger, func(s *grpc.Server) {
		proto.RegisterCatalogServiceServer(s, svc)
	}).Start(ctx)
}
