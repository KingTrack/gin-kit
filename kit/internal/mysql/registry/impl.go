package registry

import (
	"context"
	"fmt"
	"sync"

	tlscontext "github.com/KingTrack/gin-kit/kit/internal/tls/context"
	"github.com/KingTrack/gin-kit/kit/types/mysql/conf"
	"github.com/KingTrack/gin-kit/kit/types/mysql/unknown"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Registry struct {
	mu  sync.RWMutex
	dbs map[tlscontext.ResourceName]*DB
}

func New() *Registry {
	return &Registry{
		dbs: make(map[tlscontext.ResourceName]*DB, 2),
	}
}

func (r *Registry) Init(ctx context.Context, configs []conf.Config) error {
	for _, v := range configs {
		config := v
		masterConfig, slaveConfigs, err := v.Parse()
		if err != nil {
			return err
		}

		masterDB, err := newDB(ctx, masterConfig)
		if err != nil {
			return err
		}

		var slaveDBs []*gorm.DB
		for _, vv := range slaveConfigs {
			slaveConfig := vv
			slaveDB, err := newDB(ctx, slaveConfig)
			if err != nil {
				return err
			}
			slaveDBs = append(slaveDBs, slaveDB)
		}

		db := &DB{
			name:       config.Name,
			master:     masterDB,
			slaves:     slaveDBs,
			slaveIndex: 0,
		}

		r.addOrUpdate(ctx, db)
	}
	return nil
}

func (r *Registry) addOrUpdate(ctx context.Context, db *DB) {
	resourceName := tlscontext.GetResourceName(ctx, db.name)

	r.mu.Lock()
	defer r.mu.Unlock()

	if existing, ok := r.dbs[resourceName]; ok {
		// 更新已有 DB
		updateDB(existing, db)
	} else {
		// 新增 DB
		r.dbs[resourceName] = db
	}
}

func updateDB(existing, db *DB) {
	existing.mu.Lock()
	defer existing.mu.Unlock()

	// 根据需要更新配置
	existing.master = db.master
	existing.slaves = db.slaves
	existing.slaveIndex = db.slaveIndex
}

func newDB(ctx context.Context, config *conf.DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if config.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	}

	// 可选：ping 测试连接
	ctxTimeout, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()
	err = sqlDB.PingContext(ctxTimeout)
	if err != nil {
		return nil, fmt.Errorf("ping DB %s failed: %v", config.DSN, err)
	}

	return db, nil
}

func (r *Registry) GetDB(ctx context.Context, name string) *DB {
	resourceName := tlscontext.GetResourceName(ctx, name)

	r.mu.RLock()
	defer r.mu.RUnlock()

	if db, ok := r.dbs[resourceName]; ok {
		return db
	}

	unknownDB := unknown.NewDB()
	return &DB{
		name:       name,
		master:     unknownDB,
		slaves:     []*gorm.DB{unknownDB},
		slaveIndex: 0,
	}
}
