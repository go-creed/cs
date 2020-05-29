package ext_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileExt(t *testing.T) {
	f := "upload-srv/static/file/1.txt"
	base := filepath.Base(f)
	split := strings.Split(base, ".")
	fmt.Println(split)
}
