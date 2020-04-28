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

func (s *service) ParseToken(token string) (request *authPb.Request, err error) {
	t := &Token{}
	err = s.parseToken(t, token)
	if err != nil {
		return nil, fmt.Errorf("[Auth][ParseToken] %s", err)
	}
	request = new(authPb.Request)
	request.Id = t.Id
	request.UserName = t.UserName
	return
}

type Token struct {
	Id       int64
	UserName string
	jwt.StandardClaims
}

func (s *service) GenerateToken(request *authPb.Request) (string, error) {
	token := Token{}
	token.Id = request.Id
	token.UserName = request.UserName
	token.ExpiresAt = time.Now().Add(time.Hour * 24 * 7).Unix()
	generateToken, err := s.generateToken(token)
	if err != nil {
		return "", errors.WithStack(err)
	}
	// write token to redis
	if err = s.saveToCache(cache.Cache(), s.key(token.Id), generateToken); err != nil {
		return "", errors.WithStack(err)
	}
	return generateToken, nil
}
