package util_test

import (
	"fmt"
	"testing"

	"cs/public/util"
)

func TestExecLinuxShell(t *testing.T) {
	src := "/Users/gre/Go/micro/src/cs/conf/upload-srv/static/file"
	src2 := src + "/1.dmg"
	cmd := fmt.Sprintf("cd %s && ls | sort -n | xargs cat > %s", src, src2)
	fmt.Println(cmd)
	shell, err := util.ExecLinuxShell(cmd)
	fmt.Println(err)
	fmt.Println(shell)
}

func TestVerifyFile(t *testing.T) {
	file, err := util.VerifyFile("31ba830fb9de2ef49c0f803dab6bdebba1b8f526eb85e6a79a1305ddc7c2e54a", "/Users/gre/Downloads/Firefox-latest.dmg")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(file)
}
