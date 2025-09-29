package registry

import (
	otlejaegerclient "github.com/KingTrack/gin-kit/kit/internal/tracer/client/opentelemetry/jaeger"
	otleskywalkingclient "github.com/KingTrack/gin-kit/kit/internal/tracer/client/opentelemetry/skywalking"
	otlezipkinclient "github.com/KingTrack/gin-kit/kit/internal/tracer/client/opentelemetry/zipkin"
	otjaegerclient "github.com/KingTrack/gin-kit/kit/internal/tracer/client/opentracing/jaeger"
	otzipkinclient "github.com/KingTrack/gin-kit/kit/internal/tracer/client/opentracing/zipkin"
	"github.com/KingTrack/gin-kit/kit/types/tracer/conf"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func NewOTTracer(config *conf.Config) (opentracing.Tracer, error) {

	switch config.BackendName {
	case conf.BackendJaeger:
		return otjaegerclient.New(config)
	case conf.BackendZipkin:
		return otzipkinclient.New(config)
	default:
		return nil, errors.Errorf("unsupported backend: %s", config.BackendName)
	}
}

func NewOTelTracer(config *conf.Config) (trace.Tracer, error) {
	// 设定传播方式
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	switch config.BackendName {
	case conf.BackendJaeger:
		return otlejaegerclient.New(config)
	case conf.BackendZipkin:
		return otlezipkinclient.New(config)
	case conf.BackendSkywalking:
		return otleskywalkingclient.New(config)
	default:
		return nil, errors.Errorf("unsupported backend:%s", config.BackendName)
	}
}
