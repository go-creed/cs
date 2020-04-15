package handler

import (
	"context"
	"sync"

	log "github.com/micro/go-micro/v2/logger"

	userMd "cs/app/user-srv/model/user"
	userPb "cs/app/user-srv/proto/user"
	"cs/plugin/db"
)

var (
	once        sync.Once
	userService userMd.Service
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
	})
}
