package client

import (
	"context"

	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ErrEngineNotInitialized = errors.New("runtime engine not initialized")
	ErrRegistryNotAvailable = errors.New("mysql registry not available")
	ErrDBNotFound           = errors.New("mysql db not found")
)

func getErrorDB(err error) *gorm.DB {
	// 使用内存 SQLite，这是最安全的方式
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if db != nil {
		db.Error = err // 设置预设错误
	}
	return db
}

type MySQL struct {
	name string
}

func New(name string) IMySQL {
	return &MySQL{name: name}
}

func (m *MySQL) Master(ctx context.Context) *gorm.DB {
	db, ok := m.isErrorDB(ctx)
	if ok {
		return db
	}
	return runtime.Get().MySQLRegistry().GetDB(ctx, m.name).Master()
}

func (m *MySQL) isErrorDB(ctx context.Context) (*gorm.DB, bool) {
	if runtime.Get() == nil {
		return getErrorDB(ErrEngineNotInitialized), false
	}
	if runtime.Get().MySQLRegistry() == nil {
		return getErrorDB(ErrRegistryNotAvailable), false
	}
	if runtime.Get().MySQLRegistry().GetDB(ctx, m.name) == nil {
		return getErrorDB(ErrDBNotFound), false
	}
	return nil, false
}

func (m *MySQL) Slave(ctx context.Context) *gorm.DB {
	db, ok := m.isErrorDB(ctx)
	if ok {
		return db
	}
	return runtime.Get().MySQLRegistry().GetDB(ctx, m.name).Slave()
}
