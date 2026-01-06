package otel

import (
	"context"
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
)

// Init sets up OpenTelemetry tracer provider, meter provider, and global propagator.
// Returns a shutdown function to flush and cleanup providers.
func Init(conf *Config) (func(context.Context) error, error) {
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

	meterProvider = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)
	if err := runtime.Start(runtime.WithMeterProvider(meterProvider)); err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	shutdown := func(ctx context.Context) error {
		//TODO 收集所有错误 而不是只返回第一个
		err1 := tp.Shutdown(ctx)
		err2 := traceExporter.Shutdown(ctx)
		var err3 error
		var err4 error
		if meterProvider != nil {
			err3 = meterProvider.Shutdown(ctx)
		}
		if metricExporter != nil {
			err4 = metricExporter.Shutdown(ctx)
		}
		if err1 != nil {
			return err1
		}
		if err2 != nil {
			return err2
		}
		if err3 != nil {
			return err3
		}
		return err4
	}
	return shutdown, nil
}
