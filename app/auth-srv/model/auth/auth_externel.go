package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	authPb "cs/app/auth-srv/proto/auth"
	"cs/plugin/cache"
)

type service struct{}

func (s *service) GetToken(request *authPb.Request) (string, error) {
	return s.getByCache(cache.Cache(), request.Id)
}

func (s *service) ParseToken(token string) (request *authPb.Request, err error) {
	t := &Token{}
	err = s.parseToken(t, token)
	if err != nil {
		return nil, fmt.Errorf("[Auth][ParseToken] %s", err)
	}
	request = new(authPb.Request)
	request.Id = t.UserId
	request.UserName = t.UserName
	return
}

type Token struct {
	UserId   int64
	UserName string
	jwt.StandardClaims
}

func (s *service) GenerateToken(request *authPb.Request) (string, error) {
	token := Token{}
	token.UserId = request.Id
	token.UserName = request.UserName
	token.ExpiresAt = time.Now().Add(time.Hour * 24 * 7).Unix()
	generateToken, err := s.generateToken(token)
	if err != nil {
		return "", errors.WithStack(err)
	}
	// write token to redis
	if err = s.saveToCache(cache.Cache(), token.UserId, generateToken); err != nil {
		return "", errors.WithStack(err)
	}
	return generateToken, nil
}
