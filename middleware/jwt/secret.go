package jwt

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
)

var RsaPublicSecret *rsa.PublicKey

func init() {
	RsaPublicKey := os.Getenv("SC_JWT_PUBLIC_KEY")
	var err error
	RsaPublicSecret, err = jwt.ParseRSAPublicKeyFromPEM([]byte(RsaPublicKey))
	if err != nil {
		log.Printf("[JWT] convert rsa public key to rsa struct error: %s", err)
	}
}

func GetSecret(algorithm string) (interface{}, error) {
	if algorithm == "RS256" {
		return RsaPublicSecret, nil
	}
	return []byte(""), fmt.Errorf("[JWT] wrong signing method %s", algorithm)
	// if algorithm == "HS256" {
	// 	return []byte(os.Getenv("JWT_SECRET")), nil
	// } else if algorithm == "RS256" {
	// 	return RsaPublicSecret, nil
	// } else {
	// 	return []byte(""), fmt.Errorf("[JWT] wrong signing method %s", algorithm)
	// }
}
