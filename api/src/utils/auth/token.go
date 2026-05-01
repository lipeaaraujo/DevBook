package auth

import (
	"api/src/config"
	"errors"
	"fmt"
	"log"
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

func ValidateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, returnSecretKey)
	log.Println(config.JwtTokenSecret)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("Invalid token")
}

func returnSecretKey(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected token signing method: %v", token.Header["alg"])
	}
	return []byte(config.JwtTokenSecret), nil
}
