package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"cs/public/ecode"
)

type Response struct {
	Code  int         `json:"code"`
	Error error       `json:"-"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data,omitempty"`
}

type responseErr struct {
	Response
	Err string `json:"error,omitempty"`
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
	ginContext(ctx).JSON(http.StatusOK, rsp)
}

func ServerError(ctx context.Context, rsp Response) {
	errRsp(ctx, http.StatusInternalServerError, rsp)
}

func RequestError(ctx context.Context, rsp Response) {
	errRsp(ctx, http.StatusBadRequest, rsp)
}

func errRsp(ctx context.Context, code int, rsp Response) {

	switch rsp.Error.(type) {
	case *ecode.Err:
		err := rsp.Error.(*ecode.Err)
		rsp.Code = err.Code
		rsp.Msg += err.Message
		ginContext(ctx).JSONP(code, responseErr{Response: rsp, Err: err.Errord.Error()})
		return
	case *ecode.Errno:
		err := rsp.Error.(*ecode.Errno)
		rsp.Code = err.Code
		rsp.Msg = err.Error()
	default:
		rsp.Msg = rsp.Error.Error()
	}
	ginContext(ctx).JSONP(code, responseErr{Response: rsp, Err: "2"})
}
