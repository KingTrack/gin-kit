package client

import (
	"context"

	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/request"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/response"
	"github.com/pkg/errors"
)

type Client struct {
	name string
}

func New(ctx context.Context, name string) IClient {
	return &Client{name: name}
}

func (c *Client) Do(ctx context.Context, req *request.Request) (*response.Response, error) {
	if c := runtime.Get().HTTPClientRegistry().Get(c.name); c != nil {
		return c.Do(ctx, req)
	}

	return nil, errors.New("http client is nil")
}
