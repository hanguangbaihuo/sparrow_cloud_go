package frontlock

import (
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/hanguangbaihuo/sparrow_cloud_go/restclient"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/httptest"
)

func MockSucc(serviceAddr string, apiPath string, payload interface{}, kwargs ...map[string]interface{}) (restclient.Response, error) {
	return restclient.Response{Code: 200, Body: []byte(`{"code":0,"message":"ok","result":0}`)}, nil
}

func MockUsed(serviceAddr string, apiPath string, payload interface{}, kwargs ...map[string]interface{}) (restclient.Response, error) {
	return restclient.Response{Code: 200, Body: []byte(`{"code":1,"message":"ok","result":1}`)}, nil
}

func MockFail(serviceAddr string, apiPath string, payload interface{}, kwargs ...map[string]interface{}) (restclient.Response, error) {
	return restclient.Response{Code: 400, Body: []byte(`{"code":-1,"message":"bad request"}`)}, nil
}

func TestNoLock(t *testing.T) {
	var app = iris.New()

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/api/ping", CheckLock, handlePing)
	e := httptest.New(t, app)
	// test no sc-lock header
	e.GET("/api/ping").Expect().Status(iris.StatusOK)
	// test empty sc-lock header
	e.GET("/api/ping").WithHeader("Sc-Lock", "").Expect().Status(iris.StatusOK)
}

func TestAllowLock(t *testing.T) {
	var app = iris.New()

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/api/ping", CheckLock, handlePing)
	e := httptest.New(t, app)

	// mock request, allow
	patches := gomonkey.ApplyFunc(restclient.Put, MockSucc)
	defer patches.Reset()

	// test sc-lock header
	e.GET("/api/ping").WithHeader("Sc-Lock", "abc123").Expect().Status(iris.StatusOK)
}

func TestRejectLock(t *testing.T) {
	var app = iris.New()

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/api/ping", CheckLock, handlePing)
	e := httptest.New(t, app)

	// mock request, allow
	patches := gomonkey.ApplyFunc(restclient.Put, MockUsed)
	defer patches.Reset()

	// test sc-lock header, reject request
	e.GET("/api/ping").WithHeader("Sc-Lock", "abc123").Expect().Status(iris.StatusOK).Body().Contains("code").Contains("233402")
}
