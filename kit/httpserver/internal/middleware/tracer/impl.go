package middleware

import (
	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/KingTrack/gin-kit/kit/types/httpserver/middleware/tracer"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	config tracer.IConfig
}

func New(config tracer.IConfig) *Middleware {
	return &Middleware{config: config}
}

func (m *Middleware) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m.config == nil {
			c.Next()
			return
		}

		if t := runtime.Get().TracerRegistry().OTTracer(); t != nil {
			m.traceWithOT(c, t)
		} else if t := runtime.Get().TracerRegistry().OTelTracer(); t != nil {
			m.traceWithOTel(c, t)
		}
	}
}
