package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"

	authPb "cs/app/auth-srv/proto/auth"
	userPb "cs/app/user-srv/proto/user"
	_const "cs/public/const"
	middleware "cs/public/gin-middleware"
	"cs/public/session"
)

var (
	userClient userPb.UserService
	authClient authPb.AuthService
)

func Init() {
	userClient = userPb.NewUserService(_const.UserSrv, client.DefaultClient)
	authClient = authPb.NewAuthService(_const.AuthSrv, client.DefaultClient)
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
		middleware.ServerError(ctx, middleware.Response{
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
	middleware.ServerError(ctx, middleware.Response{
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
		middleware.ServerError(ctx, middleware.Response{
			Error: err,
		})
		return
	}
	middleware.Success(ctx, middleware.Response{
		Msg:  "Register Success",
		Data: login,
	})
	return
}

func Detail(ctx *gin.Context) {
	//
	middleware.Success(ctx, middleware.Response{
		Msg:  "Detail Success",
		Data: "...",
	})
}
