package middleware

import (
	"net/http"

	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/gin-gonic/gin"
)

type Capture struct {
	ctx *gin.Context
}

func NewCapture(c *gin.Context) *Capture {
	return &Capture{
		ctx: c,
	}
}

func (c *Capture) SetStatusCode(statusCode int) {
	if runtime.Get().ContextRegistry() == nil {
		return
	}

	cc := runtime.Get().ContextRegistry().Load(c.ctx)
	if cc == nil {
		return
	}

	cc.SetStatusCode(statusCode)
}

func (c *Capture) SetData(data []byte) {
	if runtime.Get().ContextRegistry() == nil {
		return
	}

	cc := runtime.Get().ContextRegistry().Load(c.ctx)
	if cc == nil {
		return
	}

	cc.SetRequestBody(data)
}

func (c *Capture) SetHeader(header http.Header) {
	if runtime.Get().ContextRegistry() == nil {
		return
	}

	cc := runtime.Get().ContextRegistry().Load(c.ctx)
	if cc == nil {
		return
	}

	cc.SetResponseHeader(header)
}
