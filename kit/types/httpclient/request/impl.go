package request

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"go.uber.org/multierr"
)

type Request struct {
	*http.Request
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
		r.Header.Add(k, v)
	}
	return r
}

func (r *Request) SetHeader(headers map[string]string) *Request {
	for k, v := range headers {
		r.Header.Set(k, v)
	}
	return r
}

func (r *Request) SetQueryValues(query map[string][]string) *Request {
	if r.URL == nil {
		r.URL = &url.URL{}
	}

	if r.URL.RawQuery == "" {
		r.URL.RawQuery = url.Values(query).Encode()
	} else {
		r.URL.RawQuery += "&" + url.Values(query).Encode()
	}
	return r
}

func (r *Request) SetURL(rawURL string) *Request {
	u, err := url.Parse(rawURL)
	if err != nil {
		r.err = multierr.Append(r.err, err)
		return r
	}
	r.URL = u
	return r
}

func (r *Request) SetMethod(method string) *Request {
	r.Method = method
	return r
}

func (r *Request) SetBodyReader(reader io.Reader) *Request {
	r.Body = io.NopCloser(reader)
	return r
}

func (r *Request) SetRequest(raw *http.Request) *Request {
	r.Request = raw
	return r
}

func (r *Request) SetMetricName(name string) *Request {
	r.metricName = name
	return r
}

func (r *Request) GetError() error {
	return r.err
}
