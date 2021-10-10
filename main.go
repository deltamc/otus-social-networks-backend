package main

import (
	"github.com/deltamc/otus-social-networks-backend/routes"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	routes.Public()
	routes.Auth()

	log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), nil))
}