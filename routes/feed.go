package routes

import (
	c "github.com/deltamc/otus-social-networks-backend/controllers"
	m "github.com/deltamc/otus-social-networks-backend/middlewares"
	"net/http"
)

func Feed() {
	http.HandleFunc("/feed", m.Cors(m.Get(m.Jwt(c.HandleFeed))))
}
