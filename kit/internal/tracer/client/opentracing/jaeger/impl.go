package client

import (
	"github.com/KingTrack/gin-kit/kit/types/tracer/conf"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

func New(config *conf.Config) (opentracing.Tracer, error) {
	cfg := jaegercfg.Configuration{
		ServiceName: config.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: config.ReportURL, // e.g. "localhost:6831"
		},
	}

	tracer, _, err := cfg.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)

	return tracer, nil
}
