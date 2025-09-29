package registry

import (
	"context"

	"github.com/KingTrack/gin-kit/kit/types/tracer/conf"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
)

type Registry struct {
	config     *conf.Config
	otTracer   opentracing.Tracer
	oTelTracer trace.Tracer
}

func New() *Registry {
	return &Registry{}
}

func (t *Registry) Init(ctx context.Context, config *conf.Config) error {
	t.config = config

	if !config.Enabled {
		return nil
	}

	switch config.Proto {
	case conf.ProtoOpenTracing:
		client, err := NewOTTracer(config)
		if err != nil {
			return err
		}
		t.otTracer = client
	case conf.ProtoOpenTelemetry:
		client, err := NewOTelTracer(config)
		if err != nil {
			return err
		}
		t.oTelTracer = client
	default:
		return errors.Errorf("unsupported proto:%s", config.Proto)
	}

	return nil
}

func (t *Registry) OTTracer() opentracing.Tracer {
	return t.otTracer
}

func (t *Registry) OTelTracer() trace.Tracer {
	return t.oTelTracer
}
