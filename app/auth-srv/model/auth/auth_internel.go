package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
)

var (
	ketSecret = []byte("235234")
)

const (
	expire = 3 * 3600 * 24 * time.Second
)

func (s *service) key(key int64) string {
	return fmt.Sprintf("jwt_%d", key)
}

func (s *service) generateToken(this jwt.Claims) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, this)
	return claims.SignedString(ketSecret)
}

func (s *service) saveToCache(rd *redis.Client, userId int64, token string) (err error) {
	return rd.Set(s.key(userId), token, expire).Err()
}

func (s *service) getByCache(rd *redis.Client, userId int64) (token string, err error) {
	err = rd.Get(s.key(userId)).Scan(&token)
	return
}

func (s *service) parseToken(this jwt.Claims, token string) error {
	claims, err := jwt.ParseWithClaims(token, this, func(token *jwt.Token) (interface{}, error) {
		return ketSecret, nil
	})
	if err != nil {
		return err
	}
	return claims.Claims.Valid()
}

func (s *service) saveTokenToRedis() {

}
