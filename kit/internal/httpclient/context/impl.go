package context

import (
	"context"

	"github.com/KingTrack/gin-kit/kit/types/datacenter/instance"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/request"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/response"
)

type HandlerFunc func(c *Context)

type Context struct {
	Ctx      context.Context
	Instance *instance.Instance
	Req      *request.Request
	Resp     *response.Response
	Err      error
	handlers []HandlerFunc
	index    int
	aborted  bool
}

func New(ctx context.Context, req *request.Request, instance *instance.Instance) *Context {
	return &Context{Ctx: ctx, Req: req, Instance: instance}
}

func (c *Context) Use(hls ...HandlerFunc) {
	c.handlers = append(c.handlers, hls...)
}

func (c *Context) Next() {
	c.index++

	if c.aborted || c.index >= len(c.handlers) {
		return
	}

	h := c.handlers[c.index]
	h(c)
}

func (c *Context) Abort() {
	c.aborted = true
}
