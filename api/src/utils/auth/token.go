package auth

import (
	"api/src/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(userId string) (string, error) {
	perms := jwt.MapClaims{
		"authorized": true,
		"exp": time.Now().Add(time.Hour * 6).Unix(),
		"userId": userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, perms)
	return token.SignedString([]byte(config.JwtTokenSecret))
}
