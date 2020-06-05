package util

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func FileSize(src string) (int64, error) {
	cmd := fmt.Sprintf("ls -l %s | awk '{print $5}'", src)
	shell, err := ExecLinuxShell(cmd)
	if err != nil {
		return 0, err
	}
	size, err := strconv.Atoi(strings.TrimSpace(shell))
	if err != nil {
		return 0, err
	}
	return int64(size), nil
}

func VerifyFile(sha256 string, src string) (bool, error) {
	s, err := Sha256(src)
	if err != nil {
		return false, err
	}
	return s == sha256, nil
}

func Sha256(src string) (string, error) {
	cmd := fmt.Sprintf("shasum -a 256 %s", src)
	shell, err := ExecLinuxShell(cmd)
	if err != nil {
		return "", err
	}
	return shell[:64], nil
}

func MergeFile(src, dest string) (string, error) {
	cmd := fmt.Sprintf("cd %s && ls | sort -n | xargs cat > %s", src, dest)
	return ExecLinuxShell(cmd)
}

// 执行 linux shell command
func ExecLinuxShell(s string) (string, error) {
	//函数返回一个io.Writer类型的*Cmd
	cmd := exec.Command("/bin/bash", "-c", s)

	//通过bytes.Buffer将byte类型转化为string类型
	var result bytes.Buffer
	cmd.Stdout = &result

	//Run执行cmd包含的命令，并阻塞直至完成
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return result.String(), err
}
