package engine

import (
	"github.com/KingTrack/gin-kit/kit/conf"
)

type Option func(e *Engine) error

func WithNamespace(config *conf.Namespace) Option {
	return func(e *Engine) error {
		if config == nil {
			return nil
		}
		return e.initResource(config)
	}
}
