package otel

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
)

var (
	metricsHandler http.Handler
	meterProvider  *sdkmetric.MeterProvider
)

// Init sets up OpenTelemetry tracer provider, meter provider, and global propagator.
// Returns a shutdown function to flush and cleanup providers.
func Init(conf *Config) (func(context.Context) error, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(conf.Endpoint),
	}
	if conf.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	} else {
		opts = append(opts, otlptracegrpc.WithDialOption(grpc.WithBlock()))
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(conf.ServiceName),
			semconv.DeploymentEnvironment(conf.Environment),
		),
	)
	if err != nil {
		return nil, err
	}

	// Metrics exporter for Prometheus scrape
	reg := prometheus.NewRegistry()
	promExporter, err := otelprom.New(otelprom.WithRegisterer(reg))
	if err != nil {
		return nil, err
	}
	meterProvider = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(promExporter),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)
	if err := runtime.Start(runtime.WithMeterProvider(meterProvider)); err != nil {
		return nil, err
	}
	metricsHandler = promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	shutdown := func(ctx context.Context) error {
		err1 := tp.Shutdown(ctx)
		err2 := exporter.Shutdown(ctx)
		var err3 error
		if meterProvider != nil {
			err3 = meterProvider.Shutdown(ctx)
		}
		if err1 != nil {
			return err1
		}
		if err2 != nil {
			return err2
		}
		return err3
	}
	return shutdown, nil
}

// MetricsHandler exposes the Prometheus scrape handler if metrics are initialized.
func MetricsHandler() http.Handler {
	return metricsHandler
}
