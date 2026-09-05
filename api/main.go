package main

import (
	"api/src/config"
	"api/src/middlewares"
	"api/src/router"

	"fmt"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config.Load()

	fmt.Println("Running API on port", config.Port)
	r := router.Generate()

	m, err := migrate.New(
		"file://migrations",
		config.DbConnectionUrl,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}

	handler := middlewares.Cors(r)

	portStr := fmt.Sprintf(":%d", config.Port)
	if err := http.ListenAndServe(portStr, handler); err != nil {
		log.Fatal(err)
	}
}
