package user

import (
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"

	userPb "cs/app/user-srv/proto/user"
	"cs/public/secret"
)

var (
	once sync.Once
	s    *service
)

type service struct {
}

func (s *service) Register(db *gorm.DB, request *userPb.UserInfo) (err error) {
	if err = s.userIsExist(db, request); err != nil {
		return fmt.Errorf("[User][Login] 判断用户存在失败, err:%s", err)
	}
	request.Password = secret.PasswordSum256(request.Password)
	if err = s.userInsert(db, request); err != nil {
		return fmt.Errorf("[User][Login] 用户数据创建失败, err:%s", err.Error())
	}
	return nil
}

func (s *service) Login(db *gorm.DB, request *userPb.UserInfo) (err error) {
	request.Password = secret.PasswordSum256(request.Password)
	if err = s.userLogin(db, request); err != nil {
		return fmt.Errorf("[User][Login] 用户数据查找失败, err:%s", err.Error())
	}
	return
}

func GetService() (Service, error) {
	if s == nil {
		return nil, fmt.Errorf("[GetService] User's Model Init Failure")
	}
	return s, nil
}

type Service interface {
	Login(db *gorm.DB, request *userPb.UserInfo) (err error)
	Register(db *gorm.DB, request *userPb.UserInfo) error
}

func Init() {
	once.Do(func() {
		s = &service{}
	})
}
