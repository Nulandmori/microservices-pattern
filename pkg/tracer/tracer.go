package tracer

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
)

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	signozToken  = os.Getenv("SIGNOZ_ACCESS_TOKEN")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

func InitTracer() func(context.Context) error {

	headers := map[string]string{
		"signoz-access-token": signozToken,
	}

	secureOption := otlpgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if len(insecure) > 0 {
		secureOption = otlpgrpc.WithInsecure()
	}

	exporter, err := otlp.NewExporter(
		context.Background(),
		otlpgrpc.NewDriver(
			secureOption,
			otlpgrpc.WithEndpoint(collectorURL),
			otlpgrpc.WithHeaders(headers),
		),
	)

	if err != nil {
		log.Fatal(err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			label.String("service.name", serviceName),
			label.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Printf("Could not set resources: ", err)
	}

	otel.SetTracerProvider(
		trace.NewTracerProvider(
			trace.WithConfig(trace.Config{DefaultSampler: trace.AlwaysSample()}),
			trace.WithSpanProcessor(trace.NewBatchSpanProcessor(exporter)),
			trace.WithSyncer(exporter),
			trace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}
