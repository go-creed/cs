package handler_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/imroc/req"
)

func TestFile(t *testing.T) {
	//base := filepath.Base("Users/gre/Downloads/Firefox-latest.dmg")
	//fmt.Println(base)
	split, file := filepath.Split("Users/gre/Downloads/Firefox-latest.dmg")
	fmt.Println(split, file)

}
func TestFileUpload(t *testing.T) {
	var (
		request = req.New()
		_url    = "http://localhost:12001/file/upload"
	)
	header := req.Header{
		"Cookie":
		"remember-me-token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjYsIlVzZXJOYW1lIjoiemciLCJleHAiOjE1OTE1ODQ0NTB9.KjlE-gvTwjI7GACNBlXscy6Y9oUnvTRUkCpB7VYDSxQ;" +
			"session-x-9f9d0332-db92-42b6-a952-62c11186b787=MTU5MDk3OTY1MHxEdi1CQkFFQ180SUFBUkFCRUFBQVBfLUNBQUlHYzNSeWFXNW5EQW9BQ0hWelpYSk9ZVzFsQm5OMGNtbHVad3dFQUFKNlp3WnpkSEpwYm1jTUNBQUdkWE5sY2tsa0JXbHVkRFkwQkFJQURBPT18fqkk9VMO3QaRjOQemKR3GArQh3xtuIWvWz9gVFVU7dE=; Path=/; Domain=localhost; Expires=Wed, 01 Jul 2020 02:47:30 GMT;",
	}
	open, err := os.Open("/Users/gre/Downloads/Firefox-latest.dmg")
	if err != nil {
		t.Fatal(err)
		return
	}
	defer open.Close()
	buffer := make([]byte, 5242880)
	var capSize int64
	var i int
	for {
		if n, err := open.ReadAt(buffer, capSize); err != nil && err != io.EOF {
			t.Fatal(err)
			return
		} else if err == nil || err == io.EOF {
			param := req.Param{
				"upload_id":  "CHUNK_6_16144ef60a40d8b8",
				"filesha256": "31ba830fb9de2ef49c0f803dab6bdebba1b8f526eb85e6a79a1305ddc7c2e54a",
				"index":      i + 1,
				"file_name":  "Firefox-latest.dmg",
			}
			capSize += int64(n)
			if post, err2 := request.Post(_url, param, header, buffer[:n]); err2 != nil {
				t.Fatal(err)
				break
			} else {
				fmt.Println(post.Response())
			}
			if err == io.EOF {
				break
			}
			i++
			//return
		} else {
			break
		}
	}
	//}

}

func TestCmd(t *testing.T) {
	var srcPath = "/Users/gre/Go/micro/src/cs/app/upload-srv/static/file"
	var destPath = "/Users/gre/Go/micro/src/cs/app/upload-srv/static/add/1.dmg"
	cmd := fmt.Sprintf("cd %s && ls | sort -n | xargs cat > %s", srcPath, destPath)
	cmd = fmt.Sprintf("cd %s && ls -tr | xargs cat > %s", srcPath, destPath)
	fmt.Println(cmd)
	command := exec.Command("/bin/zsh", "-c", cmd)
	var result bytes.Buffer
	command.Stdout = &result
	err := command.Run()
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(result.String())
}
