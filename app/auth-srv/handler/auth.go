package handler

import (
	"context"
	"sync"

	log "github.com/micro/go-micro/v2/logger"

	"cs/app/auth-srv/model/auth"
	authPb "cs/app/auth-srv/proto/auth"
)

var (
	once sync.Once
	s    auth.Service
)

type Auth struct{}

func (a *Auth) GenerateToken(ctx context.Context, request *authPb.Request, response *authPb.Response) error {
	token, err := s.GenerateToken(request)
	if err != nil {
		return err
	}
	response.Token = token
	log.Infof("[Auth][GenerateToken] id:%d , token:%s", request.Id, token)
	return nil
}

func (a *Auth) ParseToken(ctx context.Context, response *authPb.Response, request *authPb.Request) error {
	panic("implement me")
}

func Init() {
	once.Do(func() {
		log.Info("[Auth][Handler] Init ...")
		s = auth.GetService()
	})
}
