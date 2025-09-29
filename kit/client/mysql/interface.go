package client

import (
	"context"

	"gorm.io/gorm"
)

type IMySQL interface {
	Master(ctx context.Context) *gorm.DB
	Slave(ctx context.Context) *gorm.DB
}
