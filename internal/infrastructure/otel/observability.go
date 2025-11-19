package otel

import (
	"context"
	"net/url"
	"time"

	"gochat/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

// Initialize sets up tracing according to Telemetry config. Metrics reuse existing prometheus implementation.
func Initialize(ctx context.Context, serviceName string, conf *config.TelemetryConfig) (func(context.Context) error, error) {
	if conf == nil || !conf.Enabled || !conf.TraceEnabled {
		return nil, nil
	}

	u, err := url.Parse(conf.OTLPEndpoint)
	if err != nil {
		return nil, err
	}

	exp, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(u.Host), otlptracehttp.WithInsecure())
	if err != nil {
		return nil, err
	}

	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(conf.SampleRatio))

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(conf.ServiceVersion),
			attribute.String("deployment.environment", conf.Environment),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exp, sdktrace.WithBatchTimeout(5*time.Second)),
	)
	otel.SetTracerProvider(tp)
	return tp.Shutdown, nil
}
