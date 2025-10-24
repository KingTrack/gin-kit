package unknown

import (
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ErrUnknownMySQL = errors.New("unknown mysql")
)

func New() *gorm.DB {
	// 使用内存 SQLite，这是最安全的方式
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if db != nil {
		db.Error = ErrUnknownMySQL // 设置预设错误
	}
	return db
}
