package jwt

// This middleware was cloned from : https://github.com/iris-contrib/middleware/tree/v12/jwt
// we need change the default TokenExtractor
// 我们从 https://github.com/iris-contrib/middleware/tree/v12/jwt 克隆这个项目,
// 我们需要修改默认的 TokenExtractor

import (
	"os"
	"testing"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/httptest"
)

var (
	jwtSecret = []byte("My JWTSecret")
)

func TestBasicJwt(t *testing.T) {
	var app = iris.New()
	os.Setenv("JWT_SECRET", string(jwtSecret))

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/secured/ping", AutoServe, handlePing)
	e := httptest.New(t, app)

	// test normal token
	token := NewTokenWithClaims(SigningMethodHS256, MapClaims{
		"exp": time.Now().Unix() + 100,
		"uid": "abc123",
	})
	tokenString, _ := token.SignedString(jwtSecret)

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

	token := NewTokenWithClaims(SigningMethodHS256, MapClaims{
		"exp": time.Now().Unix() - 100,
		"uid": "abc123",
	})
	tokenString, _ := token.SignedString(jwtSecret)

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

	token := NewTokenWithClaims(SigningMethodHS256, MapClaims{
		"exp": time.Now().Unix() + 500,
		"uid": "abc123",
	})
	tokenString, _ := token.SignedString([]byte("wrongjwtSecret"))

	e.GET("/secured/ping").WithHeader("Authorization", "Token "+tokenString).
		Expect().Status(iris.StatusUnauthorized).Body().Contains("invalid")
}
