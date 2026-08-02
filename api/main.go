package main

import (
	"api/src/config"
	"api/src/middlewares"
	"api/src/router"

	"fmt"
	"log"
	"net/http"
)

func main() {
	config.Load()

	fmt.Println("Running API on port", config.Port)
	r := router.Generate()

	handler := middlewares.Cors(r)

	portStr := fmt.Sprintf(":%d", config.Port)
	err := http.ListenAndServe(portStr, handler)
	if err != nil {
		log.Fatal(err)
	}
}
