package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

func extraPeerNameWithOTel(wireCtx context.Context) string {
	if bge := baggage.FromContext(wireCtx); bge.Len() > 0 {
		return bge.Member(string(semconv.PeerServiceKey)).Value()
	}
	return ""
}

// otel.GetTextMapPropagator()
// otel.SetTextMapPropagator() 自定义传播
func (m *Middleware) traceWithOTel(c *gin.Context, tracerClient trace.Tracer) {
	baseCtx := c.Request.Context()
	defer func() {
		c.Request = c.Request.WithContext(baseCtx)
	}()

	spanName := m.getSpanName(c)
	// 从 HTTP Header 提取 Trace Context
	propagator := otel.GetTextMapPropagator()
	wireCtx := propagator.Extract(baseCtx, propagation.HeaderCarrier(c.Request.Header))

	// 设置peerName
	setPeerName(c, extraPeerNameWithOTel(wireCtx))

	// 构建标准的 HTTP 服务器 span 属性
	reqAttrs := []attribute.KeyValue{
		semconv.HTTPRequestMethodKey.String(c.Request.Method),
		semconv.HTTPRouteKey.String(c.FullPath()),
		semconv.HTTPRequestBodySizeKey.Int64(c.Request.ContentLength),
		semconv.ClientAddressKey.String(getClientIP(c)),
		semconv.ServerAddressKey.String(c.Request.Host),
	}

	opts := []trace.SpanStartOption{
		trace.WithAttributes(reqAttrs...),
		trace.WithSpanKind(trace.SpanKindServer),
	}

	// 启动 span
	spanCtx, span := tracerClient.Start(wireCtx, spanName, opts...)
	defer span.End()

	setTraceID(span.SpanContext().TraceID().String())

	// 将 span context 传递给后续处理
	c.Request = c.Request.WithContext(spanCtx)

	// 处理请求
	c.Next()

	// 构建响应属性
	respAttrs := []attribute.KeyValue{
		semconv.HTTPResponseBodySizeKey.Int(c.Writer.Size()),
		semconv.HTTPResponseStatusCodeKey.Int(c.Writer.Status()),
	}

	// 添加自定义属性
	if appCodeKey := m.config.GetAppCodeKey(); len(appCodeKey) > 0 {
		respAttrs = append(respAttrs, attribute.String(appCodeKey, getAppCode(c)))
	}

	if durationKey := m.config.GetRequestDurationMsKey(); len(durationKey) > 0 {
		respAttrs = append(respAttrs, attribute.Int64(durationKey, time.Since(getStartTime(c)).Milliseconds()))
	}

	// 设置响应属性
	span.SetAttributes(respAttrs...)

	// 设置响应属性（标准 HTTP 语义约定）
	if len(c.Errors) > 0 {
		span.SetStatus(codes.Error, c.Errors.String())
		for _, err := range c.Errors {
			span.RecordError(err.Err)
		}
	} else {
		span.SetStatus(codes.Ok, "")
	}
}
