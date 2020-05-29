package ext_test

import (
	"fmt"
	"testing"

	"cs/public/ext"
)

func TestExcel(t *testing.T) {
	path := "upload-srv/static/file/1.txt"
	name, ext, valid := ext.Excel().Valid(path)
	fmt.Println(name, ext, valid)
}
