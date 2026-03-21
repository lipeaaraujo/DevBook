package main

import (
	"api/config"
	"api/src/router"
	"fmt"
	"log"
	"net/http"
)

func main()  {
	config.Load()

	fmt.Println("Running API on port", config.Port)
	r := router.Generate()

	portStr := fmt.Sprintf(":%d", config.Port)
	err := http.ListenAndServe(portStr, r)
	if err != nil {
		log.Fatal(err)
	}
}
