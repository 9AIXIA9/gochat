package otel

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var (
	meterProvider *sdkmetric.MeterProvider
	enabled       atomic.Bool
)

// Init sets up OpenTelemetry tracer provider, meter provider, and global propagator.
// Returns a shutdown function to flush and cleanup providers.
func Init(conf *Config) (func(context.Context) error, error) {
	enabled.Store(false)

	if conf == nil || !conf.Enabled {
		return func(context.Context) error { return nil }, nil
	}
	conf.ApplyDefaults()

	if err := pingEndpoint(conf.Endpoint, 2*time.Second); err != nil {
		return nil, fmt.Errorf(
			"otel endpoint unreachable (%s): verify the OTEL endpoint configuration and ensure the OTEL collector is running and reachable: %w",
			conf.Endpoint,
			err,
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(conf.Endpoint),
	}
	metricOpts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(conf.Endpoint),
	}
	if conf.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
		metricOpts = append(metricOpts, otlpmetricgrpc.WithInsecure())
	} else {
		// 使用 WithTimeout 设置连接超时
		opts = append(opts, otlptracegrpc.WithTimeout(5*time.Second))
		metricOpts = append(metricOpts, otlpmetricgrpc.WithTimeout(5*time.Second))
	}

	traceExporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	metricExporter, err := otlpmetricgrpc.New(ctx, metricOpts...)
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

	readerOpts := []sdkmetric.PeriodicReaderOption{
		sdkmetric.WithInterval(conf.MetricExportInterval),
		sdkmetric.WithTimeout(conf.MetricExportTimeout),
	}

	meterProvider = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter, readerOpts...)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)
	if err := runtime.Start(runtime.WithMeterProvider(meterProvider)); err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(conf.TraceSampleRatio))),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	enabled.Store(true)

	shutdown := func(ctx context.Context) error {
		enabled.Store(false)
		var errs []error
		if err := tp.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}

		if meterProvider != nil {
			if err := meterProvider.Shutdown(ctx); err != nil {
				errs = append(errs, err)
			}
		}

		if len(errs) > 0 {
			err := errors.Join(errs...)
			return err
		}
		return nil
	}
	return shutdown, nil
}

func pingEndpoint(endpoint string, timeout time.Duration) error {
	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.Dial("tcp", endpoint)
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}
