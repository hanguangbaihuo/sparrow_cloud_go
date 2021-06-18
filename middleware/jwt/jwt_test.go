package jwt

// go test -mod=vendor ./middleware/jwt/ -v
// before run test, export SC_JWT_PUBLIC_KEY_PATH="./rsa_public.pem"

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/httptest"

	"github.com/dgrijalva/jwt-go"
	"github.com/hanguangbaihuo/sparrow_cloud_go/middleware/auth"
)

var (
	jwtSecret           = []byte("My JWTSecret")
	rsaPrivateKey, _    = ioutil.ReadFile("./rsa_private.pem")
	rsaPrivateSecret, _ = jwt.ParseRSAPrivateKeyFromPEM(rsaPrivateKey)
)

func TestBasicJwt(t *testing.T) {
	var app = iris.New()
	os.Setenv("JWT_SECRET", string(jwtSecret))

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/secured/ping", AutoServe, handlePing)
	e := httptest.New(t, app)

	// test hs256 token
	token := NewTokenWithClaims(SigningMethodHS256, MapClaims{
		"exp": time.Now().Unix() + 100,
		"uid": "abc123",
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		t.Errorf("signed hs256 token error: %s\n", err)
	}

	e.GET("/secured/ping").WithHeader("Authorization", "Token "+tokenString).
		Expect().Status(iris.StatusOK)

	// test rs256 token
	token = NewTokenWithClaims(SigningMethodRS256, MapClaims{
		"exp": time.Now().Unix() + 100,
		"uid": "abc123",
	})
	tokenString, err = token.SignedString(rsaPrivateSecret)
	if err != nil {
		t.Errorf("signed rs256 token error: %s\n", err)
	}
	// t.Logf("rsa 256 token is %s\n", tokenString)

	e.GET("/secured/ping").WithHeader("Authorization", "Token "+tokenString).
		Expect().Status(iris.StatusOK)
}
func TestEmptyToken(t *testing.T) {
	var app = iris.New()
	os.Setenv("JWT_SECRET", string(jwtSecret))

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/secured/ping", AutoServe, handlePing)
	e := httptest.New(t, app)

	// test empty jwt token header
	e.GET("/secured/ping").Expect().Status(iris.StatusOK)
}

func TestExpireToken(t *testing.T) {
	var app = iris.New()
	os.Setenv("JWT_SECRET", string(jwtSecret))

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/secured/ping", AutoServe, handlePing)
	e := httptest.New(t, app)

	// test hs256 token
	token := NewTokenWithClaims(SigningMethodHS256, MapClaims{
		"exp": time.Now().Unix() - 100,
		"uid": "abc123",
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		t.Errorf("signed hs256 token error: %s\n", err)
	}

	e.GET("/secured/ping").WithHeader("Authorization", "Token "+tokenString).
		Expect().Status(iris.StatusUnauthorized).Body().Contains("expired")

	// test rs256 token
	token = NewTokenWithClaims(SigningMethodRS256, MapClaims{
		"exp": time.Now().Unix() - 200,
		"uid": "abc123",
	})
	tokenString, err = token.SignedString(rsaPrivateSecret)
	if err != nil {
		t.Errorf("signed rs256 token error: %s\n", err)
	}

	e.GET("/secured/ping").WithHeader("Authorization", "Token "+tokenString).
		Expect().Status(iris.StatusUnauthorized).Body().Contains("expired")
}

func TestInvalidToken(t *testing.T) {
	var app = iris.New()
	os.Setenv("JWT_SECRET", string(jwtSecret))

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/secured/ping", AutoServe, handlePing)
	e := httptest.New(t, app)

	// test hs256 token
	token := NewTokenWithClaims(SigningMethodHS256, MapClaims{
		"exp": time.Now().Unix() + 500,
		"uid": "abc123",
	})
	tokenString, err := token.SignedString([]byte("wrongjwtSecret"))
	if err != nil {
		t.Errorf("signed hs256 token error: %s\n", err)
	}

	e.GET("/secured/ping").WithHeader("Authorization", "Token "+tokenString).
		Expect().Status(iris.StatusUnauthorized).Body().Contains("invalid")

	// test rs256 token
	invalidRsaToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJ1aWQiOiIxMjM0YWJjIiwiZXhwIjoxNzIyMjAwMzE2LCJpYXQiOjE2MjIxOTMxMTYsImFwcF9pZCI6ImNvcmUifQ.test"

	e.GET("/secured/ping").WithHeader("Authorization", "Token "+invalidRsaToken).
		Expect().Status(iris.StatusUnauthorized)
}

// test X-Jwt-Payload header the jwt middleware generate
func TestPayloadHeader(t *testing.T) {
	var app = iris.New()
	os.Setenv("JWT_SECRET", string(jwtSecret))

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/secured/ping", AutoServe, auth.IsAuthenticated, handlePing)
	e := httptest.New(t, app)
	//test hs256 token header
	token := NewTokenWithClaims(SigningMethodHS256, MapClaims{
		"exp": time.Now().Unix() + 100,
		"uid": "abc123",
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		t.Errorf("signed hs256 token error: %s\n", err)
	}

	e.GET("/secured/ping").WithHeader("Authorization", "Token "+tokenString).
		Expect().Status(iris.StatusOK).Body().Contains("pong")

	// test rs256 token header
	token = NewTokenWithClaims(SigningMethodRS256, MapClaims{
		"exp": time.Now().Unix() + 100,
		"uid": "abc123",
	})
	tokenString, err = token.SignedString(rsaPrivateSecret)
	if err != nil {
		t.Errorf("signed rs256 token error: %s\n", err)
	}

	e.GET("/secured/ping").WithHeader("Authorization", "Token "+tokenString).
		Expect().Status(iris.StatusOK)
}
