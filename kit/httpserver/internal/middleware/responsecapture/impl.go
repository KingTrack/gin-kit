package middleware

import (
	"github.com/KingTrack/gin-kit/kit/httpserver/internal/responsewriter"
	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/gin-gonic/gin"
)

type Middleware struct{}

func (m *Middleware) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		if runtime.Get() == nil {
			c.Next()
			return
		}

		capture := NewCapture(c)
		writer := responsewriter.New(c.Writer, capture)
		c.Writer = writer
		c.Next()
	}
}
