package conf

import "go.uber.org/zap/zapcore"

const (
	DefaultLogDir     = "./logs/"
	DefaultLogFileExt = ".log"
)

type OutputEnum string

func (e OutputEnum) ToString() string {
	return string(e)
}

func (e OutputEnum) IsAppOutput() bool {
	switch e {
	case OutputInfo, OutputDebug, OutputWarn, OutputError:
		return true
	default:
		return false
	}
}

const (
	OutputInfo     OutputEnum = "info"
	OutputDebug    OutputEnum = "debug"
	OutputWarn     OutputEnum = "warn"
	OutputError    OutputEnum = "error"
	OutputAccess   OutputEnum = "access"   // 访问日志
	OutputBusiness OutputEnum = "business" // 执行日志，sql 执行、http app 执行、redis 执行、kafka 执行、redis 执行日志
	OutputGen      OutputEnum = "gen"      // 框架日志，启动、关停、调试相关日志
	OutputCrash    OutputEnum = "crash"    // 崩溃日志
	OutputData     OutputEnum = "data"     // 数据埋点日志
)

type LevelEnum string

func (e LevelEnum) ToZapLevel() zapcore.Level {
	switch e {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func (e LevelEnum) ToString() string {
	return string(e)
}

const (
	LevelDebug LevelEnum = "debug"
	LevelInfo  LevelEnum = "info"
	LevelError LevelEnum = "error"
	LevelWarn  LevelEnum = "warn"
)

type RotateEnum string

func (e RotateEnum) ToFileFormat() string {
	if e == RotateHour {
		return "2006010215"
	}
	return "20060102"
}

const (
	RotateDay  RotateEnum = "day"
	RotateHour RotateEnum = "hour"
)

type Config struct {
	App    AppConfig  `toml:"app" json:"app" yaml:"app"`
	Core   CoreConfig `toml:"store" json:"store" yaml:"store"`
	LogDir string     `toml:"log_dir" json:"log_dir" yaml:"log_dir"`
}

type AppConfig struct {
	Level         LevelEnum  `toml:"level" json:"level" yaml:"level"`
	Rotate        RotateEnum `toml:"rotate" json:"rotate" yaml:"rotate"`
	EnableConsole bool       `toml:"enable_console" json:"enable_console" yaml:"enable_console"`
}

type CoreConfig struct {
	Rotate        RotateEnum `toml:"rotate" json:"rotate" yaml:"rotate"`
	EnableConsole bool       `toml:"enable_console" json:"enable_console" yaml:"enable_console"`
}

func (c *Config) GetLogDir() string {
	if len(c.LogDir) == 0 {
		return DefaultLogDir
	}
	return c.LogDir
}

func (c *Config) GetRotate(output OutputEnum) RotateEnum {
	if output.IsAppOutput() {
		return c.GetAppRotate()
	}
	return c.GetCoreRotate()
}

func (c *Config) GetAppLevel() LevelEnum {
	if len(c.App.Level) == 0 {
		return LevelDebug
	}
	return c.App.Level
}

func (c *Config) GetAppRotate() RotateEnum {
	if len(c.App.Rotate) == 0 {
		return RotateDay
	}
	return c.App.Rotate
}

func (c *Config) GetAppOutputs() []OutputEnum {
	switch c.App.Level {
	case LevelDebug:
		return []OutputEnum{OutputDebug, OutputInfo, OutputWarn, OutputError}
	case LevelInfo:
		return []OutputEnum{OutputInfo, OutputWarn, OutputError}
	case LevelWarn:
		return []OutputEnum{OutputWarn, OutputError}
	case LevelError:
		return []OutputEnum{OutputError}
	default:
		return []OutputEnum{OutputDebug, OutputInfo, OutputWarn, OutputError}
	}
}

func (c *Config) GetCoreRotate() RotateEnum {
	if len(c.Core.Rotate) == 0 {
		return RotateDay
	}
	return c.Core.Rotate
}
