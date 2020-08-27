package accesscontrol

import (
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/httptest"
)

func TestBaseAC(t *testing.T) {
	var (
		app = iris.New()
		ac  = New("ADMIN")
	)
	securedPingHandler := func(ctx context.Context) {
		ctx.JSON(iris.Map{
			"message": "ok",
		})
	}
	app.Get("/ac/ping", ac.Serve, securedPingHandler)
	e := httptest.New(t, app)
	e.GET("/ac/ping").Expect().Status(iris.StatusForbidden)
}
