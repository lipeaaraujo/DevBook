package routes

import (
	"api/src/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	URI string
	Method string
	Handler func(http.ResponseWriter, *http.Request)
	RequiresAuth bool
}

func Configure(r *mux.Router) *mux.Router {
	routes := userRoutes
	routes = append(routes, loginRoute)

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
