package controllers

import (
	"github.com/deltamc/otus-social-networks-backend/responses"
	"net/http"
	"os"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	responses.ResponseJson(w, map[string]string{
		"text":"hello worold!",
		"service": os.Getenv("SERVICE_NAME"),
	})
}
