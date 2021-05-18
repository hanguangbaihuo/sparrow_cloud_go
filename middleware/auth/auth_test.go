package auth

import (
	"os"
	"testing"
	"time"

	jwt "github.com/hanguangbaihuo/sparrow_cloud_go/middleware/jwt"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/httptest"
)

var (
	jwtSecret = []byte("My JWTSecret")
)

func TestNormalAuth(t *testing.T) {
	var app = iris.New()
	os.Setenv("JWT_SECRET", string(jwtSecret))

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/secured/ping", jwt.AutoServe, IsAuthenticated, handlePing)
	e := httptest.New(t, app)

	// test normal token
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Unix() + 100,
		"uid": "abc123",
	})
	tokenString, _ := token.SignedString(jwtSecret)

	e.GET("/secured/ping").WithHeader("Authorization", "Token "+tokenString).
		Expect().Status(iris.StatusOK)

	// test empty token
	e.GET("/secured/ping").Expect().Status(iris.StatusUnauthorized)
}

func TestLostJwtMiddleware(t *testing.T) {
	var app = iris.New()
	os.Setenv("JWT_SECRET", string(jwtSecret))

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/secured/ping", IsAuthenticated, handlePing)
	e := httptest.New(t, app)

	// test normal token
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Unix() + 100,
		"uid": "abc123",
	})
	tokenString, _ := token.SignedString(jwtSecret)

	e.GET("/secured/ping").WithHeader("Authorization", "Token "+tokenString).
		Expect().Status(iris.StatusUnauthorized)
}

func TestLostUIDInToken(t *testing.T) {
	var app = iris.New()
	os.Setenv("JWT_SECRET", string(jwtSecret))

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/secured/ping", jwt.AutoServe, IsAuthenticated, handlePing)
	e := httptest.New(t, app)

	// test lost uid in token
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Unix() + 100,
		// "uid": "abc123",
	})
	tokenString, _ := token.SignedString(jwtSecret)

	e.GET("/secured/ping").WithHeader("Authorization", "Token "+tokenString).
		Expect().Status(iris.StatusUnauthorized).Body().Contains("missing")
}
