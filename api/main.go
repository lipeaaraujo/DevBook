package main

import (
	"api/src/router"
	"fmt"
	"net/http"
)

func main()  {
	fmt.Println("Running API on port 4000")
	r := router.Generate()

	http.ListenAndServe(":4000", r)
}
