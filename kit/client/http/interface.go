package client

import (
	"context"

	"github.com/KingTrack/gin-kit/kit/types/httpclient/request"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/response"
)

type IClient interface {
	Do(ctx context.Context, req *request.Request) (*response.Response, error)
}
