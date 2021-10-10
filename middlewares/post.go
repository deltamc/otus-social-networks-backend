package middlewares

import (
	"encoding/json"
	"fmt"
	"github.com/deltamc/otus-social-networks-backend/responses"
	"net/http"
	"net/url"
)


func Post(h handler) handler {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if r.Header.Get("Content-Type") == "application/json" {
				var a map[string]interface{}
				decoder := json.NewDecoder(r.Body)
				decoder.DisallowUnknownFields()
				err := decoder.Decode(&a)
				if err != nil {
					responses.Response500(w, err)
					return
				}
				data := url.Values{}
				for k,v := range a {
					data.Set(k, fmt.Sprintf("%v", v))
				}
				r.Form = data
			}

			h(w, r)
			return
		}
		http.Error(w, "post only", http.StatusMethodNotAllowed)

	}
}