package auth

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"
)

var (
	ketSecret = []byte("235234")
)

const (
	expire = 3 * 3600 * 24
)

func (s *service) key(key int64) string {
	return fmt.Sprintf("jwt_%d", key)
}

func (s *service) generateToken(this jwt.Claims) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, this)
	return claims.SignedString(ketSecret)
}

func (s *service) saveToCache(conn redis.Conn, key, token string) (err error) {
	_, err = conn.Do("setex", key, expire, token)
	return err
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
