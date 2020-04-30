package auth_handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/micro/go-micro/v2/logger"

	"cs/public/rsp"
	"cs/public/session"
)

func AuthWrapper() gin.HandlerFunc {
	log.Info("use auth middleware...")
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie(session.RememberMeCookieName)
		if err != nil {
			rsp.ServerError(ctx, rsp.Response{
				Error: err,
			})
			ctx.Abort()
			return
		}
		sess := session.GetSessionGin(ctx)
		if sess.ID != "" {
			if sess.Values["valid"] != nil {
				rsp.ServerError(ctx, rsp.Response{
					Error: err,
				})
				ctx.Abort()
				return
			} else {
				userId := sess.Values["userId"].(int64)
				if userId != 0 {
					//TODO 调用远程服务
				}
			}

		}
		fmt.Println(cookie)
		ctx.Next()
	}
}
