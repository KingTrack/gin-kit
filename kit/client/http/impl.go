package client

import (
	"context"

	"github.com/KingTrack/gin-kit/kit/types/httpclient/request"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/response"
)

type Client struct {
	name string
}

func New(name string) IClient {
	return &Client{name: name}
}

func (c *Client) Call(ctx context.Context, req *request.Request) (*response.Response, error) {
	return nil, nil
}
