package auth

import "github.com/dgrijalva/jwt-go"

var (
	ketSecret = []byte("235234")
)

func (s *service) generateToken(this jwt.Claims) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, this)
	return claims.SignedString(ketSecret)
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
