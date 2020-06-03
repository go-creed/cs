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

	uploadPb "cs/app/upload-srv/proto/upload"
	"cs/public/util"
)

func TestFileUpload(t *testing.T) {

	var (
		request     = req.New()
		_url        = "http://localhost:12001/file/upload"
		_url_login  = "http://localhost:12003/login"
		_file_chunk = "http://localhost:12001/file/chunk"
		_file_merge = "http://localhost:12001/file/merge"
		err         error
		fileName    string
		filePath    = "/Users/gre/Downloads/go1.14.2.darwin-amd64.pkg"
	)
	//登录账号，获取cookies
	loginParams := req.Param{
		"user_name": "zg",
		"password":  "a11111",
	}
	post, err := request.Post(_url_login, loginParams)
	cookies := post.Response().Cookies()
	fmt.Println("cookies:", cookies)

	sha256, err := util.Sha256(filePath)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("filesha_256:", sha256)

	_, fileName = filepath.Split(filePath)
	fmt.Println("file_name:", fileName)

	size, err := util.FileSize(filePath)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("file_size:", size)

	get, err := request.Get(_file_chunk, req.Param{
		"file_size":   size,
		"filesha_256": sha256,
		"file_name":   fileName,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fmt.Println(get.String()))

	var toJson struct {
		Data uploadPb.ChunkResponse `json:"data"`
	}
	if err = get.ToJSON(&toJson); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toJson.Data)

	open, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer open.Close()
	buffer := make([]byte, toJson.Data.ChunkSize)
	var capSize int64
	var i int
	for {
		if n, err := open.ReadAt(buffer, capSize); err != nil && err != io.EOF {
			t.Fatal(err)
			return
		} else if err == nil || err == io.EOF {
			param := req.Param{
				"upload_id":  toJson.Data.UploadId,
				"filesha256": toJson.Data.Filesha256,
				"index":      i + 1,
				"file_name":  toJson.Data.FileName,
			}
			capSize += int64(n)

			if post, err2 := request.Post(_url, param, buffer[:n]); err2 != nil {
				t.Fatal(err)
				return
			} else if post.Response().StatusCode >= 300 {
				fmt.Println(post.String())
				return
			} else {
				fmt.Println(post.String())
			}
			if err == io.EOF {
				break
			}
			i++
		} else {
			break
		}
	}
	resp, err := request.Post(_file_merge, request, req.Param{
		"upload_id": toJson.Data.UploadId,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp.String())
}

func TestCmd(t *testing.T) {
	var srcPath = "/Users/gre/Go/micro/src/cs/conf/upload-srv/static/file"
	var destPath = "/Users/gre/Go/micro/src/cs/conf/upload-srv/static/add/1.dmg"
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
