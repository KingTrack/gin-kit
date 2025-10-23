package client

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"time"

	clientcontext "github.com/KingTrack/gin-kit/kit/internal/httpclient/context"
	"github.com/KingTrack/gin-kit/kit/internal/httpclient/middleware"
	runtimedatacenter "github.com/KingTrack/gin-kit/kit/runtime/datacenter"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/conf"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/request"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/response"
	"github.com/pkg/errors"
)

type HTTPClient struct {
	config     *conf.Config
	transports map[string]*http.Transport
	mu         sync.RWMutex
}

func New() *HTTPClient {
	return &HTTPClient{
		transports: make(map[string]*http.Transport, 1),
	}
}

func (c *HTTPClient) Init(ctx context.Context, config *conf.Config) error {
	var fixedURL *url.URL
	var err error
	if len(config.ProxyURL) > 0 {
		fixedURL, err = url.Parse(config.ProxyURL)
		if err != nil {
			return errors.WithMessagef(err, "http client parse proxy url failed, service_name:%s, proxy_url:%s", config.ServiceName, config.ProxyURL)
		}
	}
	_ = c.getOrAddTransport(fixedURL, config.MaxIdleConns, config.IdleConnTimeoutSec)

	c.config = config

	return nil
}

func (c *HTTPClient) Do(ctx context.Context, req *request.Request) (*response.Response, error) {
	if err := req.GetError(); err != nil {
		return nil, err
	}

	instance, err := runtimedatacenter.Get().PickInstance(c.config.ServiceName, nil, nil)
	if err != nil {
		return nil, err
	}

	clientCtx := clientcontext.New(ctx, req, instance)
	clientCtx.Use(
		middleware.Retry(&c.config.RetryerConfig),
		c.call(),
	)
	clientCtx.Next()

	return clientCtx.Resp, clientCtx.Err
}

func (c *HTTPClient) buildRequest(cc *clientcontext.Context) *http.Request {
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

func (c *HTTPClient) buildClient(cc *clientcontext.Context) *http.Client {
	timeout := time.Duration(c.config.TimeoutMs) * time.Millisecond
	if cc.Req.Timeout.Milliseconds() > 0 {
		timeout = cc.Req.Timeout
	}

	return &http.Client{
		Transport: c.getOrAddTransport(cc.Req.ProxyURL, c.config.MaxIdleConns, c.config.IdleConnTimeoutSec),
		Timeout:   timeout,
	}
}

func (c *HTTPClient) getOrAddTransport(fixedURL *url.URL, maxIdleConns int, idleConnTimeoutSec int) *http.Transport {
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

func (c *HTTPClient) call() clientcontext.HandlerFunc {
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
