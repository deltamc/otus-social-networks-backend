package middlewares

import "net/http"


func Get(h handler) handler {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h(w, r)
			return
		}
		http.Error(w, "get only", http.StatusMethodNotAllowed)
	}
}