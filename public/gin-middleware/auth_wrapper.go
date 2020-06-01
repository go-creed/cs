package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	log "github.com/micro/go-micro/v2/logger"

	"cs/public/session"
)

const (
	auth_ = "AUTH"
)

type auth struct {
	UserId int64 `json:"user_id"`
}

func GetAuth(ctx *gin.Context) (auth, error) {
	if get, ok := ctx.Get(auth_); !ok {
		return auth{}, errors.New("the userId is not obtained in the context")
	} else {
		return get.(auth), nil
	}
}

type AuthFunc func(ctx *gin.Context) (string, error)

func AuthWrapper(authFunc AuthFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie(session.RememberMeCookieName)
		if err != nil {
			ServerError(ctx, Response{
				Error: err,
			})
			ctx.Abort()
			return
		}
		sess := session.GetSessionGin(ctx)
		if sess.ID == "" {
			ServerError(ctx, Response{
				Error: errors.New("session verification failed"),
			})
			ctx.Abort()
			return
		}

		if sess.Values["userId"] == nil || sess.Values["userId"].(int64) == 0 {
			ServerError(ctx, Response{
				Error: err,
			})
			ctx.Abort()
			return
		}

		userId := sess.Values["userId"].(int64)
		//TODO Call remote service to get token
		auth := auth{UserId: userId}
		ctx.Set(auth_, auth)
		if token, err := authFunc(ctx); err != nil {
			log.Errorf("[AuthWrapper],err: %s", err)
			ServerError(ctx, Response{
				Error: err,
			})
			ctx.Abort()
			return
		} else if token != cookie {
			err = errors.New("token不一致")
			log.Error("[AuthWrapper],err:", err)
			ServerError(ctx, Response{
				Error: err,
			})
			ctx.Abort()
			return
		}
		ctx.Set("userId", userId)
		ctx.Next()
	}
}
