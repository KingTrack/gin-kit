package core

import (
	"net/http"
	"sync"
	"time"
)

type Pool struct {
	pool sync.Pool
}

func NewPool() *Pool {
	return &Pool{
		pool: sync.Pool{
			New: func() interface{} {
				return &Context{}
			},
		},
	}
}

func (p *Pool) Get() *Context {
	ctx := p.pool.Get().(*Context)
	ctx.Reset()

	return ctx
}

func (p *Pool) Put(ctx *Context) {
	if ctx == nil {
		return
	}
	ctx.Cleanup()

	p.pool.Put(ctx)
}

type Context struct {
	requestStartTime time.Time
	peerName         string
	clientIP         string
	appCode          string
	statusCode       int
	requestBody      []byte
	responseBody     []byte
	responseHeader   http.Header
	used             bool
	mu               sync.RWMutex
}

func (c *Context) IsUsed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.used
}

func (c *Context) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.requestStartTime = time.Time{}
	c.peerName = ""
	c.clientIP = ""
	c.appCode = ""
	c.requestBody = nil
	c.responseBody = nil
	c.statusCode = 0
	c.responseHeader = nil
	c.used = true
}

func (c *Context) Cleanup() {
	c.Reset()

	c.mu.Lock()
	c.used = false
	c.mu.Unlock()
}

func (c *Context) SetRequestStartTime(startTime time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.requestStartTime = startTime
}

func (c *Context) GetRequestStartTime() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.requestStartTime
}

func (c *Context) SetAppCode(code string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.appCode = code
}

func (c *Context) GetAppCode() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.appCode
}

func (c *Context) SetPeerName(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.peerName = name
}

func (c *Context) GetPeerName() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.peerName
}

func (c *Context) SetClientIP(ip string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.clientIP = ip
}

func (c *Context) GetClientIP() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.clientIP
}

func (c *Context) SetRequestBody(body []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.requestBody = body
}

func (c *Context) GetRequestBody() []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.requestBody
}

func (c *Context) SetResponseBody(body []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.responseBody = body
}

func (c *Context) GetResponseBody() []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.responseBody
}

func (c *Context) SetStatusCode(code int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.statusCode = code
}

func (c *Context) GetStatusCode() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.statusCode
}

func (c *Context) SetResponseHeader(header http.Header) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.responseHeader = header
}

func (c *Context) GetResponseHeader() http.Header {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.responseHeader
}
