package db

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"

	"cs/public/config"
)

type MysqlConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Addr     string `json:"addr"`     // including ip and port
	LogMode  bool   `json:"log_mode"` // 是否开启日志
	DbName   string `json:"db_name"`
}

func initMysql() {
	var (
		c   = config.C()
		cfg = &MysqlConfig{}
		err error
	)
	if err = c.Get("mysql", cfg); err != nil {
		log.Fatal(err)
	}
	if db, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=true&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Addr,
		cfg.DbName)); err != nil {
		log.Fatal(err)
	}

	if err = db.DB().Ping(); err != nil {
		log.Fatal(err)
	}
}
