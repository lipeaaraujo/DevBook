package middlewares

import (
	"api/src/responses"
	"api/src/utils/auth"
	"errors"
	"log"
	"net/http"
)

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI, r.Host)
		next(w, r)
	}
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := auth.ValidateToken(r); err != nil {
			responses.Error(w, http.StatusUnauthorized, errors.New("Invalid token"))
			return
		}
		next(w, r)
	}
}
