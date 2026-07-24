package auth

import (
	"api/src/apierrors"
	"api/src/config"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(userId string) (string, error) {
	perms := jwt.MapClaims{
		"authorized": true,
		"exp":        time.Now().Add(time.Hour * 6).Unix(),
		"userId":     userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, perms)
	return token.SignedString([]byte(config.JwtTokenSecret))
}

func ValidateToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, returnSecretKey)
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return token, nil
	}

	return nil, errors.New("Invalid token")
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	token := ""
	if len(strings.Split(authHeader, " ")) == 2 {
		token = strings.Split(authHeader, " ")[1]
	}
	return token
}

func returnSecretKey(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected token signing method: %v", token.Header["alg"])
	}
	return []byte(config.JwtTokenSecret), nil
}

func ExtractUserId(r *http.Request) (string, error) {
	token, err := ValidateToken(r)
	if err != nil {
		return "", apierrors.Unauthorized("Invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userId, ok := claims["userId"].(string)
		if ok {
			return userId, nil
		}
	}

	return "", apierrors.Unauthorized("Invalid token")
}
