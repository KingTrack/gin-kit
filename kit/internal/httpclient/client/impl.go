package client

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"time"

	clientcontext "github.com/KingTrack/gin-kit/kit/internal/httpclient/context"
	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/conf"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/request"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/response"
	"github.com/pkg/errors"
)

type Client struct {
	config     *conf.Config
	transports map[string]*http.Transport
	mu         sync.RWMutex
}

func New() *Client {
	return &Client{
		transports: make(map[string]*http.Transport, 1),
	}
}

func (c *Client) Init(ctx context.Context, config *conf.Config) error {
	var fixedURL *url.URL
	var err error
	if len(config.ProxyURL) > 0 {
		fixedURL, err = url.Parse(config.ProxyURL)
		if err != nil {
			return errors.WithMessage(err, "httpclient parse proxy url failed")
		}
	}
	_ = c.getOrAddTransport(fixedURL, config.MaxIdleConns, config.IdleConnTimeoutSec)

	c.config = config

	return nil
}

func (c *Client) Do(ctx context.Context, req *request.Request) (*response.Response, error) {
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

func (c *Client) buildRequest(cc *clientcontext.Context) *http.Request {
	httpReq := cc.Req.Request.Clone(cc.Ctx)

	if httpReq.URL == nil {
		httpReq.URL = &url.URL{}
	}
	if len(httpReq.URL.Scheme) == 0 {
		httpReq.URL.Scheme = cc.Instance.Schema
	}
	if len(httpReq.URL.Host) == 0 {
		httpReq.URL.Host = cc.Instance.GetHost()
	}

	return httpReq
}

func (c *Client) buildClient(cc *clientcontext.Context) *http.Client {
	timeout := time.Duration(c.config.TimeoutMs) * time.Millisecond
	if cc.Req.Timeout.Milliseconds() > 0 {
		timeout = cc.Req.Timeout
	}

	return &http.Client{
		Transport: c.getOrAddTransport(cc.Req.ProxyURL, c.config.MaxIdleConns, c.config.IdleConnTimeoutSec),
		Timeout:   timeout,
	}
}

func (c *Client) getOrAddTransport(fixedURL *url.URL, maxIdleConns int, idleConnTimeoutSec int) *http.Transport {
	key := getTransportKey(fixedURL.String())

	c.mu.RLock()
	tr, ok := c.transports[key]
	c.mu.RUnlock()
	if ok {
		return tr
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if tr, ok := c.transports[key]; ok {
		return tr
	}
	tr = newTransport(fixedURL, maxIdleConns, idleConnTimeoutSec)
	c.transports[key] = tr
	return tr
}

func newTransport(fixedURL *url.URL, maxIdleConns int, idleConnTimeoutSec int) *http.Transport {
	tr := &http.Transport{
		MaxIdleConns:    maxIdleConns,
		IdleConnTimeout: time.Duration(idleConnTimeoutSec) * time.Second,
	}

	if fixedURL != nil {
		tr.Proxy = http.ProxyURL(fixedURL)
	}

	return tr
}

func getTransportKey(rawURL string) string {
	return rawURL
}

func (c *Client) call() clientcontext.HandlerFunc {
	return func(cc *clientcontext.Context) {
		httpResp, err := c.buildClient(cc).Do(c.buildRequest(cc))
		if err != nil {
			cc.Err = err
			return
		}
		cc.Resp.Response = httpResp
		cc.Next()
	}
}
