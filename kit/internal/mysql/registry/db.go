package registry

import (
	"sync"

	"gorm.io/gorm"
)

type DB struct {
	name       string
	master     *gorm.DB
	slaves     []*gorm.DB
	mu         sync.RWMutex
	slaveIndex int64
}

func (d *DB) Master() *gorm.DB {
	return d.master
}

func (d *DB) Slave() *gorm.DB {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if len(d.slaves) == 0 {
		return d.master
	}

	// 轮询选择 slave
	db := d.slaves[d.slaveIndex]
	d.slaveIndex = (d.slaveIndex + 1) % int64(len(d.slaves))
	return db
}
