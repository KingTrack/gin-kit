package registry

import (
	"context"
	"os"
	"testing"

	tlsstore "github.com/KingTrack/gin-kit/kit/internal/tls/store"

	"github.com/KingTrack/gin-kit/kit/types/logger/conf"
	"github.com/stretchr/testify/assert"
)

func TestRegistryInitialization(t *testing.T) {
	registry := New()
	assert.NotNil(t, registry)
	assert.NotNil(t, registry.appLogger)
	assert.NotNil(t, registry.accessLogger)
	assert.NotNil(t, registry.genLogger)
	assert.NotNil(t, registry.businessLogger)
	assert.NotNil(t, registry.crashLogger)
	assert.NotNil(t, registry.dataLogger)
}

func TestRegistryGetters(t *testing.T) {
	registry := New()

	// 测试各个 getter 方法返回不同的实例
	appLogger := registry.AppLogger()
	accessLogger := registry.AccessLogger()
	genLogger := registry.GenLogger()
	businessLogger := registry.BusinessLogger()
	crashLogger := registry.CrashLogger()
	dataLogger := registry.DataLogger()

	// 验证返回的不是同一个实例
	assert.NotSame(t, appLogger, accessLogger)
	assert.NotSame(t, accessLogger, genLogger)
	assert.NotSame(t, genLogger, businessLogger)
	assert.NotSame(t, businessLogger, crashLogger)
	assert.NotSame(t, crashLogger, dataLogger)
}

func TestLoggersDefaultUsage(t *testing.T) {

	tlsstore.SetNamespace("example")

	registry := New()

	// 测试各个 logger 的默认使用（不初始化配置）
	appLogger := registry.AppLogger()
	accessLogger := registry.AccessLogger()
	genLogger := registry.GenLogger()
	businessLogger := registry.BusinessLogger()
	crashLogger := registry.CrashLogger()
	dataLogger := registry.DataLogger()

	// 测试 zap logger (app logger)
	appLogger.Info("Test message from app logger")
	appLogger.Error("Test error from app logger")

	// 测试 store logger (access logger)
	accessLogger.Print("Test message from access logger")
	accessLogger.Printf("Test formatted message: %s", "from access logger")

	// 测试 store logger (gen logger)
	genLogger.Print("Test message from gen logger")
	genLogger.Printf("Test formatted message: %s", "from gen logger")

	// 测试 store logger (business logger)
	businessLogger.Print("Test message from business logger")
	businessLogger.Printf("Test formatted message: %s", "from business logger")

	// 测试 store logger (crash logger)
	crashLogger.Print("Test message from crash logger")
	crashLogger.Printf("Test formatted message: %s", "from crash logger")

	// 测试 store logger (data logger)
	dataLogger.Print("Test message from data logger")
	dataLogger.Printf("Test formatted message: %s", "from data logger")
}

func TestRegistryInitWithConfig(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := "./logs"

	config := &conf.Config{
		LogDir: tempDir,
		App: conf.AppConfig{
			Rotate: conf.RotateHour,
		},
		Core: conf.CoreConfig{
			EnableConsole: true,
		},
	}

	registry := New()

	// 测试初始化
	err := registry.Init(context.Background(), config)
	assert.NoError(t, err)

	// 测试初始化后的 logger 使用
	appLogger := registry.AppLogger()
	accessLogger := registry.AccessLogger()
	genLogger := registry.GenLogger()

	// 测试初始化后的 logger 功能
	appLogger.Info("Test message after init")

	accessLogger.Info("Test message after init")
	accessLogger.Sugar().Infof("Test formatted message after init: %s", "test")
	genLogger.Warn("Test warning after init")

	// 验证目录创建
	_, err = os.Stat(tempDir)
	assert.NoError(t, err)

	// 清理测试目录
	//defer os.RemoveAll(tempDir)
}

func TestMultipleRegistryInstances(t *testing.T) {
	// 创建多个 registry 实例，验证它们是独立的
	registry1 := New()
	registry2 := New()

	// 验证两个实例是不同的
	assert.NotSame(t, registry1, registry2)

	// 验证它们的 logger 也是不同的
	appLogger1 := registry1.AppLogger()
	appLogger2 := registry2.AppLogger()
	assert.NotSame(t, appLogger1, appLogger2)

	accessLogger1 := registry1.AccessLogger()
	accessLogger2 := registry2.AccessLogger()
	assert.NotSame(t, accessLogger1, accessLogger2)
}
