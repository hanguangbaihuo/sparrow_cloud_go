package jwt

import (
	"fmt"
	"io/ioutil"
	"os"
)

var RsaPublicSecret []byte

func init() {
	RsaPublicKeyPath := os.Getenv("PUBLIC_KEY_PATH")
	var err error
	RsaPublicSecret, err = ioutil.ReadFile(RsaPublicKeyPath)
	if err != nil {
		panic(fmt.Sprintf("[JWT] read rsa public file err: %s", err))
	}
}

func GetSecret(algorithm string) ([]byte, error) {
	if algorithm == "HS256" {
		return []byte(os.Getenv("JWT_SECRET")), nil
	} else if algorithm == "RS256" {
		return RsaPublicSecret, nil
	} else {
		return []byte(""), fmt.Errorf("[JWT] wrong signing method %s", algorithm)
	}
}
