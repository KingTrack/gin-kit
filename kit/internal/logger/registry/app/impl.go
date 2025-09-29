package app

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
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

type Logger struct {
	*zap.Logger
}

func New() *Logger {
	encoder := zapcore.NewConsoleEncoder(newEncoderConfig())
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &Logger{
		Logger: logger,
	}
}

func (l *Logger) Init(config *conf.Config) {
	l.Logger = l.createLogger(config)
}

func (l *Logger) createLogger(config *conf.Config) *zap.Logger {
	// 创建 lineencoder 配置
	encoder := zapcore.NewConsoleEncoder(newEncoderConfig())

	// 创建 write syncer
	var writeSyncers []zapcore.WriteSyncer
	for _, v := range config.GetAppOutputs() {
		t := v
		writeSyncers = append(writeSyncers, syncer.New(t, config))
	}

	// 控制台输出
	if config.App.EnableConsole {
		writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))
	}

	// 创建
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writeSyncers...),
		config.GetAppLevel().ToZapLevel(),
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

func (l *Logger) Info(template string) {
	l.With().Sugar().Info(template)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.With().Sugar().Infof(template, args...)
}

func (l *Logger) Debug(template string) {
	l.With().Sugar().Debug(template)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.With().Sugar().Debugf(template, args...)
}

func (l *Logger) Warn(template string) {
	l.With().Sugar().Warn(template)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.With().Sugar().Warnf(template, args...)
}

func (l *Logger) Error(template string) {
	l.With().Sugar().Error(template)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.With().Sugar().Errorf(template, args...)
}

func (l *Logger) Fatal(template string) {
	l.With().Sugar().Fatal(template)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.With().Sugar().Fatalf(template, args...)
}

func (l *Logger) Panic(template string) {
	l.With().Sugar().Panic(template)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.With().Sugar().Panicf(template, args...)
}
