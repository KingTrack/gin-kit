package client

import (
	"github.com/KingTrack/gin-kit/kit/types/tracer/conf"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

func New(config *conf.Config) (opentracing.Tracer, error) {
	reporter := zipkinhttp.NewReporter(config.ReportURL) // e.g. "http://localhost:9411/api/v2/spans"

	// create zipkin recorder + tracer
	localEndpoint, err := zipkin.NewEndpoint(config.ServiceName, "0.0.0.0:0")
	if err != nil {
		return nil, err
	}

	zipkinTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(localEndpoint))
	if err != nil {
		return nil, err
	}

	// wrap for opentracing
	otTracer := zipkinot.Wrap(zipkinTracer)

	opentracing.SetGlobalTracer(otTracer)

	return otTracer, nil
}
