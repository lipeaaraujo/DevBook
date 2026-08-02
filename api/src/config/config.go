package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	DbConnectionString = ""
	Port               = 0
	JwtTokenSecret     = ""
	FrontendURL        = ""
)

func Load() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	Port, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		Port = 4000
	}

	DbConnectionString = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	JwtTokenSecret = os.Getenv("JWT_TOKEN_SECRET")
	FrontendURL = os.Getenv("FRONTEND_URL")
}
