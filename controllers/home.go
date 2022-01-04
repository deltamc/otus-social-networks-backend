package controllers

import (
	"github.com/deltamc/otus-social-networks-backend/responses"
	"net/http"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	responses.ResponseJson(w, map[string]string{
		"text":"hello worold!",
	})
}
