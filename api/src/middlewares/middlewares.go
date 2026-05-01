package middlewares

import (
	"api/src/responses"
	"api/src/utils/auth"
	"errors"
	"log"
	"net/http"
	"strings"
)

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI, r.Host)
		next(w, r)
	}
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		token := ""
		if len(strings.Split(authHeader, " ")) == 2 {
			token = strings.Split(authHeader, " ")[1]
		}
		if err := auth.ValidateToken(token); err != nil {
			responses.Error(w, http.StatusUnauthorized, errors.New("Invalid token"))
			return
		}
		next(w, r)
	}
}
