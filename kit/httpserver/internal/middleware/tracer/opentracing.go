package middleware

import (
	"time"

	"github.com/KingTrack/gin-kit/kit/internal/tracer/client/common"
	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracinglog "github.com/opentracing/opentracing-go/log"
)

// 使用自定义opentracing而不是global opentracing
func (m *Middleware) traceWithOT(c *gin.Context, traceClient opentracing.Tracer) {
	baseCtx := c.Request.Context()
	defer func() {
		c.Request = c.Request.WithContext(baseCtx)
	}()

	spanName := m.getSpanName(c)
	var span opentracing.Span
	// 提取上游的 span context
	carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
	wireCtx, err := traceClient.Extract(opentracing.HTTPHeaders, carrier)
	if err != nil {
		// 创建新的 Span
		runtime.Get().LoggerRegistry().GenLogger().Printf("opentracing middleware extract upstream span error:%v", err)
		span = traceClient.StartSpan(spanName, ext.SpanKindRPCServer)
	} else {
		// 延续现有 Span
		span = traceClient.StartSpan(spanName, opentracing.ChildOf(wireCtx), ext.SpanKindRPCServer)
		setPeerName(c, common.ExtractPeerNameWithOT(wireCtx))
	}
	defer span.Finish()

	setTraceID(common.ExtractTraceIDWithOT(span.Context()))

	// 设置 HTTP 相关标签
	ext.HTTPMethod.Set(span, c.Request.Method)
	ext.HTTPUrl.Set(span, c.Request.URL.String())
	span.SetTag("http.request.body.size", c.Request.ContentLength)
	span.SetTag("client.address", getClientIP(c))
	span.SetTag("server.address", c.Request.Host)

	// 将 span 注入到请求上下文
	spanCtx := opentracing.ContextWithSpan(c.Request.Context(), span)
	c.Request = c.Request.WithContext(spanCtx)

	// 处理请求
	c.Next()

	// 添加自定义tag
	if appCodeKey := m.config.GetAppCodeKey(); len(appCodeKey) > 0 {
		span.SetTag(appCodeKey, getAppCode(c))
	}

	if durationKey := m.config.GetRequestDurationMsKey(); len(durationKey) > 0 {
		span.SetTag(durationKey, time.Since(getStartTime(c)).Milliseconds())
	}

	// 设置响应tag
	span.SetTag("http.response.body.size", c.Writer.Size())
	ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))

	// 处理错误
	if len(c.Errors) > 0 {
		ext.Error.Set(span, true)
		span.LogFields(
			opentracinglog.String("event", "error"),
			opentracinglog.String("message", c.Errors.String()),
		)
	}
}
