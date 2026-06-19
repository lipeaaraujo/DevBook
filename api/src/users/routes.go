package users

import (
	"api/src/router/routes"
	"net/http"
)

func CreateUserRoutes(userController *UserController) []routes.Route {
	return []routes.Route{
		{
			URI: "/users",
			Method: http.MethodPost,
			Handler: userController.CreateUser,
			RequiresAuth: false,
		},
		{
			URI: "/users",
			Method: http.MethodGet,
			Handler: userController.GetUsers,
			RequiresAuth: true,
		},
		{
			URI: "/users/{userId}",
			Method: http.MethodGet,
			Handler: userController.GetUser,
			RequiresAuth: true,
		},
		{
			URI: "/users/{userId}",
			Method: http.MethodPut,
			Handler: userController.UpdateUser,
			RequiresAuth: true,
		},
		{
			URI: "/users/{userId}",
			Method: http.MethodDelete,
			Handler: userController.DeleteUser,
			RequiresAuth: true,
		},
		{
			URI: "/users/{userId}/follow",
			Method: http.MethodPost,
			Handler: userController.FollowUser,
			RequiresAuth: true,
		},
	}
}
