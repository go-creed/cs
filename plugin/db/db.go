package db

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/micro/go-micro/v2/logger"
)

var (
	db   *sql.DB
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
	db, err = sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=true&loc=Local", "root", "root", "localhost:3306", "cs"))
	if err != nil {
		log.Fatalf("db open failure %s", err)
	}
	db.SetConnMaxLifetime(time.Second * 1000)

	if err = db.Ping(); err != nil {
		log.Errorf("db connect failure %s", err)
	}
	log.Info("db init the connection success")
}

func DB() *sql.DB {
	return db
}
