package httpserver

import (
	contextmiddleware "github.com/KingTrack/gin-kit/kit/httpserver/internal/middleware/context"
	loggermiddleware "github.com/KingTrack/gin-kit/kit/httpserver/internal/middleware/logger"
	metricmiddleware "github.com/KingTrack/gin-kit/kit/httpserver/internal/middleware/metric"
	recovermiddleware "github.com/KingTrack/gin-kit/kit/httpserver/internal/middleware/recover"
	responsecapturemiddleware "github.com/KingTrack/gin-kit/kit/httpserver/internal/middleware/responsecapture"
	tracermiddleware "github.com/KingTrack/gin-kit/kit/httpserver/internal/middleware/tracer"
	"github.com/gin-gonic/gin"
)

func WithRecovery(middleware *recovermiddleware.Middleware) Option {
	return func(s *Server) {
		s.recoverMiddleware = middleware.Build()
	}
}

func WithLogger(middleware *loggermiddleware.Middleware) Option {
	return func(s *Server) {
		s.tracerMiddleware = middleware.Build()
	}
}

func WithTracer(middleware *tracermiddleware.Middleware) Option {
	return func(s *Server) {
		s.tracerMiddleware = middleware.Build()
	}
}

func WithMetric(middleware *metricmiddleware.Middleware) Option {
	return func(s *Server) {
		s.metricMiddleware = middleware.Build()
	}
}

func WithContext(middleware *contextmiddleware.Middleware) Option {
	return func(s *Server) {
		s.contextMiddleware = middleware.Build()
	}
}

func WithResponseCapture(middleware *responsecapturemiddleware.Middleware) Option {
	return func(s *Server) {
		s.responseCaptureMiddleware = middleware.Build()
	}
}

func WithAppMiddleware(middlewares ...gin.HandlerFunc) Option {
	return func(s *Server) {
		s.appMiddlewares = append(s.appMiddlewares, middlewares...)
	}
}
