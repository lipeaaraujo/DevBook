package routes

import (
	"api/src/controllers"
	"net/http"
)

var userRoutes = []Route{
	{
		URI: "/users",
		Method: http.MethodPost,
		Handler: controllers.CreateUser,
		RequiresAuth: false,
	},
	{
		URI: "/users",
		Method: http.MethodGet,
		Handler: controllers.GetUsers,
		RequiresAuth: false,
	},
	{
		URI: "/users/{userId}",
		Method: http.MethodGet,
		Handler: controllers.GetUser,
		RequiresAuth: false,
	},
	{
		URI: "/users/{userId}",
		Method: http.MethodPut,
		Handler: controllers.UpdateUser,
		RequiresAuth: false,
	},
	{
		URI: "/users/{userId}",
		Method: http.MethodDelete,
		Handler: controllers.DeleteUser,
		RequiresAuth: false,
	},
}
