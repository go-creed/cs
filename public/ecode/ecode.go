package ecode

import (
	"fmt"
	"github.com/pkg/errors"
)

// 基础定义错误码
type Errno struct {
	Code    int
	Message string
}

//继承Error接口
func (err Errno) Error() string {
	return err.Message
}

// 包装错误结构，用于在常用错误之上添加具体的错误信息
type Err struct {
	Code    int    // 错误码
	Message string // 展示给用户看的
	Errord  error  // 保存内部错误信息
}

//继承Error接口
func (err *Err) Error() string {
	return fmt.Sprintf("Err - message: %s, error: %s", err.Message, err.Errord)
}

// 使用 错误码 和 error 创建新的 错误
func New(errno *Errno, err error) *Err {
	if err == nil {
		err = errors.New("尚未添加错误描述")
	}
	return &Err{
		Code:    errno.Code,
		Message: errno.Message,
		Errord:  err,
	}
}

// 解码错误信息, 获取 Code 和 Message
// 用的应该不多，正常直接调用New(errno *Errno, err error)
func DecodeErr(err error) (int, string) {
	if err == nil {
		return Success.Code, Success.Message
	}
	switch typed := err.(type) {
	case *Err:
		if typed.Errord != nil {
			typed.Message = typed.Message + " 具体内容是 " + typed.Errord.Error()
		}
		return typed.Code, typed.Message
	case *Errno:
		return typed.Code, typed.Message
	default:
	}
	return ErrInternalServer.Code, err.Error()
}
