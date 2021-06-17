package jwt

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
)

var RsaPublicSecret *rsa.PublicKey

func init() {
	RsaPublicKeyPath := os.Getenv("SC_JWT_PUBLIC_KEY_PATH")
	RsaPublicKey, err := ioutil.ReadFile(RsaPublicKeyPath)
	if err != nil {
		log.Printf("[JWT] read rsa public file err: %s", err)
		return
	}
	RsaPublicSecret, err = jwt.ParseRSAPublicKeyFromPEM(RsaPublicKey)
	if err != nil {
		log.Printf("[JWT] convert rsa public key to rsa struct error: %s", err)
	}
}

func GetSecret(algorithm string) (interface{}, error) {
	if algorithm == "HS256" {
		return []byte(os.Getenv("JWT_SECRET")), nil
	} else if algorithm == "RS256" {
		return RsaPublicSecret, nil
	} else {
		return []byte(""), fmt.Errorf("[JWT] wrong signing method %s", algorithm)
	}
}
