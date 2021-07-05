package cors

import (
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/httptest"
)

func TestSigleRouterOptions(t *testing.T) {
	var app = iris.New()
	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}
	app.AllowMethods(iris.MethodOptions)

	app.Get("/cors/get", Serve, handlePing)
	e := httptest.New(t, app)
	e.OPTIONS("/cors/get").Expect().Status(iris.StatusNoContent).Header("Access-Control-Allow-Origin").Equal("*")
}

func TestAllRouterOptions(t *testing.T) {
	var app = iris.New()
	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}
	app.Use(Serve)
	app.AllowMethods(iris.MethodOptions)

	app.Get("/cors/get", handlePing)
	app.Post("/cors/post", handlePing)

	e := httptest.New(t, app)
	e.OPTIONS("/cors/get").Expect().Status(iris.StatusNoContent).Header("Access-Control-Allow-Origin").Equal("*")
	e.OPTIONS("/cors/post").Expect().Status(iris.StatusNoContent).Header("Access-Control-Allow-Origin").Equal("*")
}
