package response

import (
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	*http.Response
}

func (r *Response) JSON(v any) error {
	if r.Response == nil || r.Response.Body == nil {
		return io.EOF
	}

	defer func() { _ = r.Response.Body.Close() }()
	return json.NewDecoder(r.Response.Body).Decode(v)
}

// 需要外部显示调用r.Close()
func (r *Response) Stream() io.ReadCloser {
	if r.Response == nil || r.Response.Body == nil {
		return io.NopCloser(nil)
	}
	return r.Response.Body
}

// 需要外部显示调用r.Close()
func (r *Response) Copy(w io.Writer) (int64, error) {
	if r.Response == nil || r.Response.Body == nil {
		return 0, io.EOF
	}
	return io.Copy(w, r.Body)
}

func (r *Response) Close() error {
	if r.Response != nil && r.Response.Body != nil {
		return r.Response.Body.Close()
	}
	return nil
}

func (r *Response) StatusCode() int {
	if r.Response == nil {
		return 0
	}
	return r.Response.StatusCode
}

func (r *Response) Status() string {
	if r.Response == nil {
		return ""
	}
	return r.Response.Status
}

func (r *Response) Header() http.Header {
	if r.Response == nil {
		return http.Header{}
	}
	return r.Response.Header
}

func (r *Response) ContentType() string {
	if r.Response == nil {
		return ""
	}
	return r.Response.Header.Get("Content-Type")
}

func (r *Response) ContentLength() int64 {
	if r.Response == nil {
		return 0
	}
	return r.Response.ContentLength
}
