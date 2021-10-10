package auth

import (
	"github.com/deltamc/otus-social-networks-backend/models/users"
	"github.com/deltamc/otus-social-networks-backend/responses"
	"github.com/dgrijalva/jwt-go/v4"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Response struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int `json:"expires_in"`
}

func GetBearerResponse(w http.ResponseWriter, r *http.Request, user users.User)  error{

	expires, err := strconv.ParseInt(os.Getenv("JWT_TOKEN_EXPIRES_MINUTE"), 10, 32)
	if err != nil {
		return err
	}
	expirationTime := time.Now().Add(time.Duration(expires) * time.Minute)


	var jwtKey = []byte(os.Getenv("SECRET_KEY"))

	claims := &Claims{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(jwtKey)
	if err != nil{
		return err
	}
	w.Header().Add("Access-Control-Expose-Headers", "Authorization")
	w.Header().Add("Authorization", accessToken)
	data := Response{
		AccessToken: accessToken,
		TokenType: "bearer",
		ExpiresIn: int(expires) * 60,
	}

	responses.ResponseJson(w, data)

	return nil
}
