package handler

import (
	"github.com/gin-gonic/gin"

	authPb "cs/app/auth-srv/proto/auth"
	"cs/public/gin-middleware"
)

func Auth() middleware.AuthFunc {
	return func(ctx *gin.Context) (string, error) {
		auth, err := middleware.GetAuth(ctx)
		if err != nil {
			return "", err
		}
		if token, err := authClient.GetToken(ctx, &authPb.Request{
			Id: auth.UserId,
		}); err != nil {
			return "", err
		} else {
			return token.Token, nil
		}
	}
}
