package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code  int         `json:"code"`
	Error error       `json:"error,omitempty"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data,omitempty"`
}

func ginContext(ctx context.Context) *gin.Context {
	switch ctx.(type) {
	case *gin.Context:
		return ctx.(*gin.Context)
	case *MicroContext:
		return ctx.(*MicroContext).Context
	default:
		return new(gin.Context)
	}
}

func Success(ctx context.Context, rsp Response) {
	ginContext(ctx).JSONP(http.StatusOK, rsp)
}

func ServerError(ctx context.Context, rsp Response) {
	ginContext(ctx).JSONP(http.StatusInternalServerError, rsp)
}

func RequestError(ctx context.Context, rsp Response) {
	ginContext(ctx).JSONP(http.StatusBadRequest, rsp)
}
