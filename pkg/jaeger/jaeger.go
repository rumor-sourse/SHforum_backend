package jaeger

import (
	"SHforum_backend/internal/settings"
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"strconv"
	"time"
)

// newJaegerTraceProvider 创建一个 Jaeger Trace Provider
func newJaegerTraceProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	jaegerHttpEndpoint := settings.Conf.JaegerConfig.Host + ":" + strconv.Itoa(settings.Conf.JaegerConfig.HttpCollectorPort)
	// 创建一个使用 HTTP 协议连接本机Jaeger的 Exporter
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(jaegerHttpEndpoint),
		otlptracehttp.WithInsecure())
	if err != nil {
		return nil, err
	}
	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName(settings.Conf.Name)))
	if err != nil {
		return nil, err
	}
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // 采样
		sdktrace.WithBatcher(exp, sdktrace.WithBatchTimeout(time.Second)),
	)
	return traceProvider, nil
}

// InitTracer 初始化 Tracer
func InitTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	tp, err := newJaegerTraceProvider(ctx)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)
	return tp, nil
}
