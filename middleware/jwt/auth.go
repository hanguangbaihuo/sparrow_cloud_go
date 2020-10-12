package jwt

import (
	"github.com/dgrijalva/jwt-go"
)

func authenticate(token *jwt.Token) User {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return User{}
	}
	var id string
	id, ok = claims["uid"].(string)
	if !ok {
		id = claims["id"].(string)
	}
	if id == "" {
		return User{}
	}
	return User{
		ID:              id,
		IsAuthenticated: true,
	}
}
