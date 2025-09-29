package core

import (
	"os"

	"github.com/KingTrack/gin-kit/kit/internal/logger/registry/syncer"
	tlsstore "github.com/KingTrack/gin-kit/kit/internal/tls/store"
	"github.com/KingTrack/gin-kit/kit/types/logger/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey: "msg",
		LineEnding: zapcore.DefaultLineEnding,
	}
}

type Logger struct {
	*zap.Logger
	output conf.OutputEnum
}

func New(output conf.OutputEnum) *Logger {
	encoder := zapcore.NewConsoleEncoder(newEncoderConfig())
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &Logger{
		Logger: logger,
		output: output,
	}
}

func (l *Logger) Init(config *conf.Config) {
	l.Logger = l.createLogger(config)
}

func (l *Logger) createLogger(config *conf.Config) *zap.Logger {
	encoder := zapcore.NewConsoleEncoder(newEncoderConfig())

	// 创建 write syncer
	var writeSyncers []zapcore.WriteSyncer
	writeSyncers = append(writeSyncers, syncer.New(l.output, config))

	// 控制台输出
	if config.Core.EnableConsole {
		writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))
	}

	// 创建 store
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writeSyncers...),
		zap.DebugLevel,
	)

	// 创建 logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return logger
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(zap.String("namespace", tlsstore.GetNamespace()), zap.String("trace_id", tlsstore.GetTraceID())).
			With(fields...),
	}
}

func (l *Logger) Printf(template string, args ...interface{}) {
	l.With().Sugar().Debugf(template, args...)
}

func (l *Logger) Print(template string) {
	l.With().Sugar().Debug(template)
}
