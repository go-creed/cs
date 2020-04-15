package rsp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code  int         `json:"code"`
	Error error       `json:"error,omitempty"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data,omitempty"`
}

func Success(ctx *gin.Context, rsp Response) {
	ctx.JSONP(http.StatusOK, rsp)
}

func ServerError(ctx *gin.Context, rsp Response) {
	ctx.JSONP(http.StatusInternalServerError, rsp)
}
