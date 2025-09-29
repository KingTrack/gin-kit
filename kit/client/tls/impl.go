package client

import (
	"context"

	tlscontext "github.com/KingTrack/gin-kit/kit/internal/tls/context"
	"github.com/KingTrack/gin-kit/kit/internal/tracer/client/common"
	"github.com/KingTrack/gin-kit/kit/runtime"
)

func NewBackendContext(namespace string) context.Context {
	ctx := tlscontext.NewBackground(tlscontext.Value(namespace))

	traceID := startBackendTrace()
	if len(traceID) > 0 {
		ctx = tlscontext.WithTraceID(ctx, tlscontext.Value(traceID))
	}
	return ctx
}

func startBackendTrace() string {
	tracerRegistry := runtime.Get().TracerRegistry()
	if tracerRegistry.OTTracer() != nil {
		span := tracerRegistry.OTTracer().StartSpan("Backend Root")
		defer span.Finish()

		return common.ExtractTraceIDWithOT(span.Context())
	} else if tracerRegistry.OTelTracer() != nil {
		_, span := tracerRegistry.OTelTracer().Start(context.Background(), "Backend Root")
		defer span.End()

		return span.SpanContext().TraceID().String()
	}

	return ""
}
