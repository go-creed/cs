package db

import (
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
	err  error
)

func Init() {
	once.Do(func() {
		initMysql()
	})
}

func DB() *gorm.DB {
	return db
}
