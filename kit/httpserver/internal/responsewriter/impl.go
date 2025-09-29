package responsewriter

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ICapture interface {
	SetStatusCode(statusCode int)
	SetData(data []byte)
	SetHeader(header http.Header)
}

type Writer struct {
	gin.ResponseWriter
	capture        ICapture
	headerCaptured bool
}

func New(w gin.ResponseWriter, capture ICapture) *Writer {
	return &Writer{
		ResponseWriter: w,
		capture:        capture,
		headerCaptured: false,
	}
}

func (w *Writer) WriteHeader(statusCode int) {
	if w.capture != nil {
		w.capture.SetStatusCode(statusCode)
		if !w.headerCaptured {
			w.capture.SetHeader(w.ResponseWriter.Header().Clone())
			w.headerCaptured = true
		}
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *Writer) Write(data []byte) (int, error) {
	if w.capture != nil {
		dataCopy := make([]byte, len(data))
		copy(dataCopy, data)
		w.capture.SetData(dataCopy)
		if !w.headerCaptured {
			w.capture.SetHeader(w.ResponseWriter.Header().Clone())
			w.headerCaptured = true
		}
	}
	return w.ResponseWriter.Write(data)
}

func (w *Writer) Header() http.Header {
	return w.ResponseWriter.Header()
}
