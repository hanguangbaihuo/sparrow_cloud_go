package auth

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/httptest"
)

var (
	jwtSecret = []byte("My JWTSecret")
	payload   = map[string]interface{}{"uid": "abc123", "exp": 1722200316, "iat": 1622193116, "app_id": "core"}
)

func TestBase64Payload(t *testing.T) {
	var app = iris.New()

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/secured/ping", IsAuthenticated, handlePing)
	e := httptest.New(t, app)

	// test base64 payload
	data, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("jsom marshal payload error: %s\n", err)
		return
	}
	base64Data := base64.StdEncoding.EncodeToString(data)

	e.GET("/secured/ping").WithHeader("X-Jwt-Payload", base64Data).
		Expect().Status(iris.StatusOK).Body().Contains("pong")
}

func TestTextPayload(t *testing.T) {
	var app = iris.New()

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/secured/ping", IsAuthenticated, handlePing)
	e := httptest.New(t, app)

	// test text payload
	data, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("jsom marshal payload error: %s\n", err)
		return
	}
	e.GET("/secured/ping").WithHeader("X-Jwt-Payload", string(data)).
		Expect().Status(iris.StatusOK).Body().Contains("pong")
}

func TestEmptyPayloadHeader(t *testing.T) {
	var app = iris.New()
	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/secured/ping", IsAuthenticated, handlePing)
	e := httptest.New(t, app)

	e.GET("/secured/ping").Expect().Status(iris.StatusUnauthorized)
}

func TestLostUIDInToken(t *testing.T) {
	var app = iris.New()
	os.Setenv("JWT_SECRET", string(jwtSecret))

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/secured/ping", IsAuthenticated, handlePing)
	e := httptest.New(t, app)

	// test text payload
	data, err := json.Marshal(map[string]interface{}{"app_id": "core"})
	if err != nil {
		t.Errorf("jsom marshal payload error: %s\n", err)
		return
	}

	e.GET("/secured/ping").WithHeader("X-Jwt-Payload", string(data)).
		Expect().Status(iris.StatusUnauthorized).Body().Contains("missing").Contains("user")
}
