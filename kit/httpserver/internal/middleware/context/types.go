package middleware

import (
	"net/http"
)

type IContextConfig interface {
	ParseNamespace(header http.Header) string
	ParseClientIP(header http.Header, uri string) string
}
