package server

import (
	"github.com/KingTrack/gin-kit/kit/httpserver"
)

var (
	srv *httpserver.Server
)

func Run() {
	srv.Run()
}
