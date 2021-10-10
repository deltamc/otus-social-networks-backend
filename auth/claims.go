package auth
import (
	jwt "github.com/dgrijalva/jwt-go/v4"
)

type Claims struct {
	jwt.StandardClaims
	UserId int64 `json:"user_id"`
}