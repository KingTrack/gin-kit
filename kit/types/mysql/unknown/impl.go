package unknown

import (
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ErrUnknownMySQLDB = errors.New("unknown mysql db")
)

func NewDB() *gorm.DB {
	// 使用内存 SQLite，这是最安全的方式
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if db != nil {
		db.Error = ErrUnknownMySQLDB // 设置预设错误
	}
	return db
}
