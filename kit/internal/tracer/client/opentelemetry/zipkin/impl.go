package client

import (
	"github.com/KingTrack/gin-kit/kit/types/tracer/conf"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

func New(config *conf.Config) (trace.Tracer, error) {
	// Create Zipkin exporter
	exp, err := zipkin.New(config.ReportURL) // e.g. "http://localhost:9411/api/v2/spans"
	if err != nil {
		return nil, err
	}

	// Create trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(config.ServiceName),
		)),
	)

	otel.SetTracerProvider(tp)

	tracer := tp.Tracer("")

	return tracer, nil
}
