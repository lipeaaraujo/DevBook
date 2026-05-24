package router

import (
	"api/src/db"
	"api/src/login"
	"api/src/middlewares"
	"api/src/users"
	"log"

	"github.com/gorilla/mux"
)

func Generate() *mux.Router {
	r := mux.NewRouter()

	db, err := db.Connect()
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}

	userRepo := users.NewUserRepo(db)
	userService := users.NewUserService(userRepo)
	userController := users.NewUserController(userService)

	loginService := login.NewLoginService(userRepo)
	loginController := login.NewLoginController(loginService)

	routes := users.CreateUserRoutes(userController)
	routes = append(routes, login.CreateLoginRoutes(loginController)...)

	for _, route := range routes {
		configuredHandler := route.Handler
		configuredHandler = middlewares.Logger(configuredHandler)
		if route.RequiresAuth {
			configuredHandler = middlewares.Auth(configuredHandler)
		}
		r.HandleFunc(route.URI, configuredHandler).Methods(route.Method)
	}

	return r
}
