package request

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/multierr"
)

type Request struct {
	*http.Request

	Timeout  time.Duration
	ProxyURL *url.URL

	err        error
	metricName string
}

func New(ctx context.Context) *Request {
	r := &Request{
		Request: &http.Request{
			Header: make(http.Header),
		},
	}
	r.Request = r.Request.WithContext(ctx)
	return r
}

func (r *Request) AddHeader(headers map[string]string) *Request {
	for k, v := range headers {
		r.Request.Header.Add(k, v)
	}
	return r
}

func (r *Request) SetHeader(headers map[string]string) *Request {
	for k, v := range headers {
		r.Request.Header.Set(k, v)
	}
	return r
}

func (r *Request) SetQueryValues(query map[string][]string) *Request {
	if r.Request.URL == nil {
		r.Request.URL = &url.URL{}
	}

	if r.Request.URL.RawQuery == "" {
		r.Request.URL.RawQuery = url.Values(query).Encode()
	} else {
		r.Request.URL.RawQuery += "&" + url.Values(query).Encode()
	}
	return r
}

func (r *Request) SetURL(rawURL string) *Request {
	u, err := url.Parse(rawURL)
	if err != nil {
		r.err = multierr.Append(r.err, err)
		return r
	}
	r.Request.URL = u
	return r
}

func (r *Request) SetMethod(method string) *Request {
	r.Request.Method = method
	return r
}

func (r *Request) SetBodyReader(reader io.Reader) *Request {
	r.Request.Body = io.NopCloser(reader)
	return r
}

func (r *Request) SetRequest(raw *http.Request) *Request {
	r.Request = raw
	return r
}

func (r *Request) SetTimeout(timeout time.Duration) *Request {
	r.Timeout = timeout
	return r
}

func (r *Request) SetProxyURL(rawURL string) *Request {
	u, err := url.Parse(rawURL)
	if err != nil {
		r.err = multierr.Append(r.err, err)
		return r
	}

	r.ProxyURL = u
	return r
}

func (r *Request) SetMetricName(name string) *Request {
	r.metricName = name
	return r
}

func (r *Request) GetError() error {
	return r.err
}
