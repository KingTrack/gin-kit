package middleware

import (
	"time"

	"github.com/KingTrack/gin-kit/kit/internal/tls/store"

	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/gin-gonic/gin"
)

func (m *Middleware) getSpanName(c *gin.Context) string {
	spanName := m.config.GetSpanName(c.Request.Method, c.Request.URL.Path)
	if len(spanName) == 0 {
		spanName = c.Request.Method + " " + c.Request.URL.Path
	}
	return spanName
}

func getStartTime(c *gin.Context) time.Time {
	cc := runtime.Get().ContextRegistry().Load(c)
	if cc == nil {
		return time.Now()
	}
	return cc.GetRequestStartTime()
}

func getAppCode(c *gin.Context) string {
	cc := runtime.Get().ContextRegistry().Load(c)
	if cc == nil {
		return ""
	}
	return cc.GetAppCode()
}

func getClientIP(c *gin.Context) string {
	cc := runtime.Get().ContextRegistry().Load(c)
	if cc == nil {
		return ""
	}
	return cc.GetClientIP()
}

func setPeerName(c *gin.Context, peerName string) {
	cc := runtime.Get().ContextRegistry().Load(c)
	if cc == nil {
		return
	}

	cc.SetPeerName(peerName)
}

func setTraceID(traceID string) {
	if len(traceID) > 0 {
		store.SetTraceID(traceID)
	}
}
