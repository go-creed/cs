package auth

import (
	"sync"

	authPb "cs/app/auth-srv/proto/auth"
)

var (
	once sync.Once
	s    Service
)

type Service interface {
	GenerateToken(request *authPb.Request) (string, error)
	ParseToken(token string) (request *authPb.Request, err error)
	GetToken(request *authPb.Request) (string, error)
}

func GetService() Service {
	return s
}

func Init() {
	once.Do(func() {
		s = &service{}
	})
}
