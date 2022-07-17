package routes

import (
	c "github.com/deltamc/otus-social-networks-backend/controllers"
	m "github.com/deltamc/otus-social-networks-backend/middlewares"
	"net/http"
)

func Public() {
	http.HandleFunc("/", m.Cors(m.Get(c.HandleHome)))
	http.HandleFunc("/test", m.Cors(m.Get(c.HandleHome)))
	http.HandleFunc("/sign-up", m.Cors(m.Post(c.HandleSignUp)))
	http.HandleFunc("/sign-in", m.Cors(m.Post(c.HandleSignIn)))
	http.HandleFunc("/users", m.Cors(m.Get(c.HandleUsers)))
	http.HandleFunc("/gen-user", m.Cors(m.Post(c.HandleGenUser)))
}
