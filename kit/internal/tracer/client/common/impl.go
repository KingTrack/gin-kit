package common

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracingzipkin "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/uber/jaeger-client-go"
)

func ExtractTraceIDWithOT(spanCtx opentracing.SpanContext) string {
	if jaegerCtx, ok := spanCtx.(jaeger.SpanContext); ok {
		return jaegerCtx.TraceID().String()
	}

	if zipkinCtx, ok := spanCtx.(opentracingzipkin.SpanContext); ok {
		return zipkinCtx.TraceID.String()
	}

	return ""
}

func ExtractPeerNameWithOT(wireCtx opentracing.SpanContext) string {
	var peerName string
	wireCtx.ForeachBaggageItem(func(k, v string) bool {
		if k == string(ext.PeerService) {
			peerName = v
			return false // 找到了就停止遍历
		}
		return true // 继续遍历
	})
	return peerName
}
