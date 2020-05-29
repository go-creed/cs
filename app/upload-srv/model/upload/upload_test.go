package upload

import (
	"fmt"
	"log"
	"testing"

	uploadSrv "cs/app/upload-srv/proto/upload"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func TestDB(t *testing.T) {
	open, err := gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=true&loc=Local", "root", "root", "localhost:3306", "cs"))
	if err != nil {
		log.Fatal(err)
		return
	}
	open.SingularTable(true)
	info := uploadSrv.FileMate{}

	err = open.Table("file_mate").Find(&info).Error
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(info)
	//db, err = sql.Open()

}
