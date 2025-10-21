package client

import (
	"context"
	"net/http"

	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/request"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/response"
	"github.com/pkg/errors"
)

type Client struct {
	err  error
	name string
}

func New(ctx context.Context, name string, httpClient ...*http.Client) IClient {
	if len(httpClient) > 0 {
		if err := runtime.Get().HTTPClientRegistry().AddClient(name, httpClient[0]); err != nil {
			return &Client{
				err:  err,
				name: name,
			}
		}
	}
	return &Client{name: name}
}

func (c *Client) Do(ctx context.Context, req *request.Request) (*response.Response, error) {
	if c.err != nil {
		return nil, c.err
	}

	if c := runtime.Get().HTTPClientRegistry().Get(c.name); c != nil {
		return c.Do(ctx, req)
	}

	return nil, errors.New("http client is nil")
}
