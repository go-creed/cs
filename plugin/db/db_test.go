package db

import (
	"fmt"
	"testing"

	"github.com/prometheus/common/log"
)

func TestDB(t *testing.T) {
	Init()
	rows, err2 := DB().DB().Query("select userName,password from user")
	type xxx struct {
		userName string
		password string
	}
	for rows.Next() {
		var xx xxx
		//var userName string
		//var password string
		err2 := rows.Scan(&xx)
		if err2 != nil {
			log.Error(err)
		}
		//fmt.Println(userName, password)
		fmt.Println(xx)
	}

	fmt.Println(rows, err2)
	prepare, err := DB().DB().Prepare("select * from user")
	if err != nil {
		log.Errorf("test db failure %s", err)
	}
	var x interface{}
	query, err := prepare.Query(&x)
	fmt.Println(query)
}
