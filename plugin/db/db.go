package db

import (
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/micro/go-micro/v2/logger"
)

var (
	db   *gorm.DB
	once sync.Once
	err  error
)

func Init() {
	once.Do(func() {
		connect()
	})
}
func connect() {
	log.Info("db init the connection start")
	db, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=true&loc=Local", "root", "root", "localhost:3306", "cs"))
	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetConnMaxLifetime(time.Second * 1000)
	if err = db.DB().Ping(); err != nil {
		log.Fatalf("db connect failure %+v", err)
	}
	log.Info("db init the connection success")
}

//func connect() {
//	log.Info("db init the connection start")
//	db, err = sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=true&loc=Local", "root", "root", "localhost:3306", "cs"))
//	if err != nil {
//		log.Fatalf("db open failure %s", err)
//	}
//	db.SetConnMaxLifetime(time.Second * 1000)
//	if err = db.Ping(); err != nil {
//		log.Errorf("db connect failure %s", err)
//	}
//	log.Info("db init the connection success")
//}

func DB() *gorm.DB {
	return db
}
