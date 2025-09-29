package client

import (
	"context"

	"github.com/KingTrack/gin-kit/kit/types/tracer/conf"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewSkyWalkingTracerOTel 创建 SkyWalking TracerRegistry（基于 OTLP/gRPC）
// endpoint 示例: "oap.skywalking.apache.org:11800" 或 "localhost:11800"
func New(config *conf.Config) (trace.Tracer, error) {
	// 创建 gRPC 客户端连接
	var opts []grpc.DialOption

	// 如果是 HTTPS/带 TLS 的 OAP 服务，使用 credentials
	// opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	// 否则使用 insecure（开发环境）
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 创建 OTLP/gRPC Exporter
	exporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithEndpoint(config.ReportURL),
		otlptracegrpc.WithDialOption(opts...),
	)
	if err != nil {
		return nil, err
	}

	// 创建 Resource，标识服务
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(config.ServiceName),
	)

	// 创建 TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // 生产环境可改为 ParentBased(AlwaysSample)
	)

	// 设置全局 TracerProvider
	otel.SetTracerProvider(tp)

	// 创建命名 TracerRegistry
	tracer := tp.Tracer("")

	return tracer, nil
}
