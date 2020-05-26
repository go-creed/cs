package ecode

import (
	"errors"
	"fmt"
	"testing"
)

func TestErr_Error(t *testing.T) {
	fmt.Println(ErrInternalServer)
	err := New(ErrInternalServer, errors.New("1"))
	fmt.Println(err)
}