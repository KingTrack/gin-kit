package client

import (
	"context"
	"net/http"

	clientcontext "github.com/KingTrack/gin-kit/kit/internal/httpclient/context"
	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/conf"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/request"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/response"
)

type Client struct {
	*http.Client
	config *conf.Config
}

func New() *Client {
	return &Client{}
}

func (c *Client) Call(ctx context.Context, req *request.Request) (*response.Response, error) {
	if err := req.GetError(); err != nil {
		return nil, err
	}

	instance, err := runtime.Get().DatacenterRegistry().PickInstance(c.config.ServiceName, nil, nil)
	if err != nil {
		return nil, err
	}

	clientCtx := clientcontext.New(ctx, req, instance)
	clientCtx.Use(c.call())
	clientCtx.Next()

	return clientCtx.Resp, clientCtx.Err
}

func (c *Client) build(cc *clientcontext.Context) *http.Request {
	httpReq := cc.Req.Request.Clone(cc.Ctx)

	if len(httpReq.URL.Scheme) == 0 {
		httpReq.URL.Scheme = cc.Instance.Schema
	}
	if len(httpReq.URL.Host) == 0 {
		httpReq.URL.Host = cc.Instance.GetHost()
	}

	return httpReq
}

func (c *Client) call() clientcontext.HandlerFunc {
	return func(cc *clientcontext.Context) {
		httpResp, err := c.Do(c.build(cc))
		if err != nil {
			cc.Err = err
			return
		}
		cc.Resp.Response = httpResp
	}
}
