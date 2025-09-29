package httpserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/gin-gonic/gin"
)

type Option func(*Server)

type Server struct {
	*gin.Engine

	// 外部用户中间件
	appMiddlewares []gin.HandlerFunc

	recoverMiddleware         gin.HandlerFunc
	loggerMiddleware          gin.HandlerFunc
	metricMiddleware          gin.HandlerFunc
	tracerMiddleware          gin.HandlerFunc
	contextMiddleware         gin.HandlerFunc
	responseCaptureMiddleware gin.HandlerFunc
}

func New(opts ...Option) *Server {
	s := &Server{
		Engine: gin.New(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Run() error {
	if runtime.Get() == nil {
		fmt.Printf("runtime engine is nil, Server can not be running")
	}

	addr := fmt.Sprintf(":%d", runtime.Get().ServerConfig().Port)
	fmt.Println("Server is running on ", addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      s,
		ReadTimeout:  time.Duration(runtime.Get().ServerConfig().ReadTimeoutSec) * time.Second,
		WriteTimeout: time.Duration(runtime.Get().ServerConfig().WriteTimeoutSec) * time.Second,
		IdleTimeout:  time.Duration(runtime.Get().ServerConfig().IdleTimeoutSec) * time.Second,
	}

	s.applyMiddleware()

	return server.ListenAndServe()
}

func (s *Server) applyMiddleware() {
	s.Use(s.recoverMiddleware)
	s.Use(s.loggerMiddleware)
	s.Use(s.metricMiddleware)
	s.Use(s.tracerMiddleware)
	s.Use(s.contextMiddleware)
	s.Use(s.responseCaptureMiddleware)
	s.Use(s.appMiddlewares...)
}
