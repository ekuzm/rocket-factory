package tracer

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

type Config interface {
	ServiceName() string
	Endpoint() string
	ServiceVersion() string
	DeploymentEnvironment() string
	Compressor() string
	RetryEnabled() bool
	RetryInitialInterval() time.Duration
	RetryMaxInterval() time.Duration
	RetryMaxElapsedTime() time.Duration
	Timeout() time.Duration
}

func Init(ctx context.Context, cfg Config) error {
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(cfg.Endpoint()),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithTimeout(cfg.Timeout()),
		otlptracegrpc.WithCompressor(cfg.Compressor()),
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
			Enabled:         cfg.RetryEnabled(),
			InitialInterval: cfg.RetryInitialInterval(),
			MaxInterval:     cfg.RetryMaxInterval(),
			MaxElapsedTime:  cfg.RetryMaxElapsedTime(),
		}),
	)
	if err != nil {
		return fmt.Errorf("create exporter: %w", err)
	}

	attribute, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName()),
			semconv.ServiceVersion(cfg.ServiceVersion()),
			attribute.String("environment", cfg.DeploymentEnvironment()),
		),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithContainer(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		return fmt.Errorf("create attribute resource: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(attribute),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1.0))),
	)

	otel.SetTracerProvider(provider)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return nil
}

func Shutdown(ctx context.Context) error {
	provider := otel.GetTracerProvider()
	if provider == nil {
		return nil
	}

	tracerProvider, ok := provider.(*sdktrace.TracerProvider)
	if !ok {
		return nil
	}

	err := tracerProvider.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown tracer provider: %w", err)
	}

	return nil
}
