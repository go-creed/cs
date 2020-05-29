package handler

import (
	"context"
	"sync"

	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"

	authPb "cs/app/auth-srv/proto/auth"
	userMd "cs/app/user-srv/model/user"
	userPb "cs/app/user-srv/proto/user"
	"cs/plugin/db"
	_const "cs/public/const"
)

var (
	once        sync.Once
	userService userMd.Service
	userClient  authPb.AuthService
)

type User struct{}

func (u *User) Register(ctx context.Context, request *userPb.Request, response *userPb.Response) error {
	log.Info("[User][Register]:Start...")
	err := userService.Register(db.DB(), request.UserInfo)
	if err != nil {
		log.Errorf("[User][Register]:%s", err.Error())
		return err
	}
	log.Info("[User][Register]:End...")
	return nil
}

func (u *User) Login(ctx context.Context, info *userPb.Request, response *userPb.Response) error {
	log.Info("[User][Login]:Start...")
	err := userService.Login(db.DB(), info.UserInfo)
	if err != nil {
		log.Errorf("[User][Login]:%s", err.Error())
		return err
	}
	// 生成token
	token, err := userClient.GenerateToken(ctx, &authPb.Request{
		Id:       info.UserInfo.Id,
		UserName: info.UserInfo.UserName,
	})
	if err != nil {
		log.Errorf("[User][Login]:%s", err.Error())
		return err
	}
	response.Token = token.Token
	response.UserName = info.UserInfo.UserName
	response.UserId = info.UserInfo.Id
	// 这里需要
	log.Info("[User][Login]:End...")
	return nil
}

func Init() {
	var err error
	once.Do(func() {
		userService, err = userMd.GetService()
		if err != nil {
			log.Fatal("[Upload] Handler Init Failure , %s", err)
			return
		}
		userClient = authPb.NewAuthService(_const.AuthSrv, client.DefaultClient)
	})
}
