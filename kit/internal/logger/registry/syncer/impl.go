package syncer

import (
	"path/filepath"
	"time"

	"github.com/KingTrack/gin-kit/kit/types/logger/conf"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func getFilename(output conf.OutputEnum, config *conf.Config) string {
	// 自定义格式：access-2025091216.log (YYYYMMDDHH)
	now := time.Now()
	fileFormat := now.Format(config.GetRotate(output).ToFileFormat())
	filename := output.ToString() + "-" + fileFormat + conf.DefaultLogFileExt
	return filepath.Join(config.GetLogDir(), filename)
}

func New(output conf.OutputEnum, config *conf.Config) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   getFilename(output, config),
		MaxSize:    102400, // 自定义文件大小
		MaxAge:     0,      // 自定义保留天数
		MaxBackups: 0,      // 自定义备份数
		LocalTime:  true,
		Compress:   false, // 自定义压缩
	})
}
