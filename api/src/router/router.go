package router

import (
	"api/src/db"
	"api/src/login"
	"api/src/middlewares"
	"api/src/posts"
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

	postRepo := posts.NewPostRepo(db)
	postService := posts.NewPostService(postRepo)
	postController := posts.NewPostController(postService)

	routes := users.CreateUserRoutes(userController)
	routes = append(routes, login.CreateLoginRoutes(loginController)...)
	routes = append(routes, posts.CreatePostRoutes(postController)...)

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
