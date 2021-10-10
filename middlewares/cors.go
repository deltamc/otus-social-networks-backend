package middlewares

import "net/http"

func Cors(h handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {

		//w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")


		if r.Method == http.MethodOptions {
			w.WriteHeader(200)
			return
		}
		h(w, r)
	}
}
