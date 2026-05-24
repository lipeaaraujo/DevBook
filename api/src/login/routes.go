package login

import (
	"api/src/router/routes"
	"net/http"
)

func CreateLoginRoutes(loginController *LoginController) []routes.Route {
	return []routes.Route{
		{
			URI: "/login",
			Method: http.MethodPost,
			Handler: loginController.Login,
			RequiresAuth: false,
		},
	}
}
