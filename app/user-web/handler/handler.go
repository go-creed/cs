package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"

	userPb "cs/app/user-srv/proto/user"
	"cs/public/rsp"
	"cs/public/session"
)

var (
	userClient userPb.UserService
)

func Init() {
	userClient = userPb.NewUserService("go.micro.cs.service.user", client.DefaultClient)
}

func Login(ctx *gin.Context) {
	user_name := ctx.PostForm("user_name")
	password := ctx.PostForm("password")
	request := userPb.Request{}
	request.UserInfo = &userPb.UserInfo{
		UserName: user_name,
		Password: password,
	}
	login, err := userClient.Login(ctx, &request)
	if err != nil {
		rsp.ServerError(ctx, rsp.Response{
			Error: err,
		})
		return
	}

	ctx.Writer.Header().Add("set-cookie", "application/json; charset=utf-8")
	expire := time.Now().Add(30 * time.Minute)
	cookie := http.Cookie{Name: session.RememberMeCookieName, Value: login.Token, Path: "/", Expires: expire, MaxAge: 90000}
	http.SetCookie(ctx.Writer, &cookie)

	sessionGin := session.GetSessionGin(ctx)
	sessionGin.Values["userId"] = login.UserId
	sessionGin.Values["userName"] = login.UserName
	_ = sessionGin.Save(ctx.Request, ctx.Writer)

	login.UserId = 0
	rsp.Success(ctx, rsp.Response{
		Msg:  "Login Success",
		Data: login,
	})
}

func Register(ctx *gin.Context) {
	// Introducing parameter validators in the future like "validator" https://juejin.im/post/5e89dc725188257399158b5d
	user_name := ctx.PostForm("user_name")
	password := ctx.PostForm("password")
	phone := ctx.PostForm("phone")
	request := userPb.Request{}
	request.UserInfo = &userPb.UserInfo{
		UserName: user_name,
		Password: password,
		Phone:    phone,
	}
	login, err := userClient.Register(ctx, &request)
	if err != nil {
		rsp.ServerError(ctx, rsp.Response{
			Error: err,
		})
		return
	}
	rsp.Success(ctx, rsp.Response{
		Msg:  "Register Success",
		Data: login,
	})
	return
}
