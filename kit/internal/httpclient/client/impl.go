package client

import (
	"context"
	"net/http"

	"github.com/KingTrack/gin-kit/kit/types/httpclient/request"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/response"

	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/conf"
)

type Client struct {
	*http.Client
	config *conf.Config
}

func New() *Client {

}

func (c *Client) Call(ctx context.Context, req *request.Request) (*response.Response, error) {
	if err := req.GetError(); err != nil {
		return nil, err
	}

	instance, err := runtime.Get().DatacenterRegistry().PickInstance(c.config.ServiceName, nil, nil)
	if err != nil {
		return nil, err
	}

	// 支持插件流程

	return nil, nil
}
