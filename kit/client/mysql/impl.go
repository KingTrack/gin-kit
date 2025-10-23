package client

import (
	"context"

	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/KingTrack/gin-kit/kit/types/mysql/unknown"
	"gorm.io/gorm"
)

type MySQL struct {
	name string
}

func New(name string) IMySQL {
	return &MySQL{name: name}
}

func (m *MySQL) Master(ctx context.Context) *gorm.DB {
	if runtime.Get().MySQLRegistry() == nil {
		return unknown.New()
	}
	return runtime.Get().MySQLRegistry().GetDB(ctx, m.name).Master()
}

func (m *MySQL) Slave(ctx context.Context) *gorm.DB {
	if runtime.Get().MySQLRegistry() == nil {
		return unknown.New()
	}
	return runtime.Get().MySQLRegistry().GetDB(ctx, m.name).Slave()
}
