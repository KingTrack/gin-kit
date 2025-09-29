package client

import (
	"github.com/KingTrack/gin-kit/kit/internal/logger/registry/app"
	"github.com/KingTrack/gin-kit/kit/runtime"
)

func Logger() *app.Logger {
	return runtime.Get().LoggerRegistry().AppLogger()
}
