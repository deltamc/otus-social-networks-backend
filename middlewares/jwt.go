package middlewares

import (
	"github.com/deltamc/otus-social-networks-backend/auth"
	"github.com/deltamc/otus-social-networks-backend/models/users"
	"github.com/deltamc/otus-social-networks-backend/responses"
	"github.com/dgrijalva/jwt-go/v4"
	"net/http"
	"os"
	"strings"
)


func Jwt(h handlerAuth) handler {

	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if headerParts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var jwtKey = []byte(os.Getenv("SECRET_KEY"))

		tknStr := headerParts[1]
		claims := &auth.Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := users.GetUserById(claims.UserId)

		if err != nil {
			responses.Response500(w, err)
		}



		h(w, r, user)
	}
}