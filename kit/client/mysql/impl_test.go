package client

import (
	"context"
	"testing"

	"github.com/KingTrack/gin-kit/kit/types/mysql/unknown"
	"github.com/stretchr/testify/assert"
)

func TestMySQL_Master(t *testing.T) {
	type record struct {
		ID string `grom:"column:id"`
	}
	assert.Equal(t, New("my.mysql").Master(context.Background()).
		Create(&record{}).Error, unknown.ErrUnknownMySQLDB)
}

func TestMySQL_Slave(t *testing.T) {
	assert.Equal(t, New("my.mysql").Slave(context.Background()).
		Select("*").Where(map[string]interface{}{"id": 1}).Error, unknown.ErrUnknownMySQLDB)
}
