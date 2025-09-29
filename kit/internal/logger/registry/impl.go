package registry

import (
	"context"
	"os"

	applogger "github.com/KingTrack/gin-kit/kit/internal/logger/registry/app"
	"github.com/KingTrack/gin-kit/kit/internal/logger/registry/core"
	corelogger "github.com/KingTrack/gin-kit/kit/internal/logger/registry/core"
	"github.com/KingTrack/gin-kit/kit/types/logger/conf"
	"github.com/pkg/errors"
)

type Registry struct {
	config         *conf.Config
	appLogger      *applogger.Logger
	accessLogger   *corelogger.Logger
	genLogger      *corelogger.Logger
	businessLogger *corelogger.Logger
	crashLogger    *corelogger.Logger
	dataLogger     *corelogger.Logger
}

func New() *Registry {
	return &Registry{
		appLogger:      applogger.New(),
		accessLogger:   corelogger.New(conf.OutputAccess),
		genLogger:      corelogger.New(conf.OutputGen),
		businessLogger: corelogger.New(conf.OutputBusiness),
		crashLogger:    corelogger.New(conf.OutputCrash),
		dataLogger:     corelogger.New(conf.OutputData),
	}
}

func (l *Registry) Init(ctx context.Context, config *conf.Config) error {
	l.config = config

	if err := l.makeLogDir(); err != nil {
		return errors.WithMessage(err, "failed to create log directory")
	}

	l.appLogger.Init(config)
	l.accessLogger.Init(config)
	l.genLogger.Init(config)
	l.businessLogger.Init(config)
	l.crashLogger.Init(config)
	l.dataLogger.Init(config)

	return nil
}

func (l *Registry) makeLogDir() error {
	return os.MkdirAll(l.config.GetLogDir(), 0755)
}

func (l *Registry) AppLogger() *applogger.Logger {
	return l.appLogger
}

func (l *Registry) AccessLogger() *core.Logger {
	return l.accessLogger
}

func (l *Registry) GenLogger() *core.Logger {
	return l.genLogger
}

func (l *Registry) BusinessLogger() *core.Logger {
	return l.businessLogger
}

func (l *Registry) CrashLogger() *core.Logger {
	return l.crashLogger
}

func (l *Registry) DataLogger() *core.Logger {
	return l.dataLogger
}
