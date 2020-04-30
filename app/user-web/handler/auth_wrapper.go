package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	log "github.com/micro/go-micro/v2/logger"

	authPb "cs/app/auth-srv/proto/auth"
	"cs/public/rsp"
	"cs/public/session"
)

func AuthWrapper() gin.HandlerFunc {
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
		if sess.ID == "" {
			err = errors.New("session verification failed")
			rsp.ServerError(ctx, rsp.Response{
				Error: err,
			})
			ctx.Abort()
			return
		}

		if sess.Values["valid"] != nil || sess.Values["userId"].(int64) == 0 {
			rsp.ServerError(ctx, rsp.Response{
				Error: err,
			})
			ctx.Abort()
			return
		}

		userId := sess.Values["userId"].(int64)
		//TODO Call remote service to get token
		token, err := authClient.GetToken(ctx, &authPb.Request{
			Id: userId,
		})
		if err != nil {
			log.Errorf("[AuthWrapper],err: %s", err)
			rsp.ServerError(ctx, rsp.Response{
				Error: err,
			})
			ctx.Abort()
			return
		}
		if token.Token != cookie {
			err = errors.New("token不一致")
			log.Error("[AuthWrapper],err:", err)
			rsp.ServerError(ctx, rsp.Response{
				Error: err,
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
