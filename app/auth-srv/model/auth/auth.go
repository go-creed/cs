package auth

import (
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"

	authPb "cs/app/auth-srv/proto/auth"
)

var (
	once sync.Once
	s    Service
)

type service struct {
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
	return s.generateToken(token)
}

type Service interface {
	GenerateToken(request *authPb.Request) (string, error)
}

func GetService() Service {
	return s
}

func Init() {
	once.Do(func() {
		s = &service{}
	})
}
