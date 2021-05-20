package accesscontrol

import (
	"os"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey"
	"github.com/hanguangbaihuo/sparrow_cloud_go/middleware/auth"
	jwt "github.com/hanguangbaihuo/sparrow_cloud_go/middleware/jwt"
	"github.com/hanguangbaihuo/sparrow_cloud_go/restclient"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/httptest"
)

var (
	jwtSecret = []byte("My JWTSecret")
)

func MockAllow(serviceAddr string, apiPath string, payload interface{}, kwargs ...map[string]interface{}) (restclient.Response, error) {
	return restclient.Response{Code: 200, Body: []byte(`{"has_perm":true}`)}, nil
}

func MockReject(serviceAddr string, apiPath string, payload interface{}, kwargs ...map[string]interface{}) (restclient.Response, error) {
	return restclient.Response{Code: 403, Body: []byte(`{"has_perm":false}`)}, nil
}

func MockGetFail(serviceAddr string, apiPath string, payload interface{}, kwargs ...map[string]interface{}) (restclient.Response, error) {
	return restclient.Response{Code: 400, Body: []byte(`{"message":"wrong parameter"}`)}, nil
}

// 请求访问控制服务成功，返回200
func TestAllowAC(t *testing.T) {
	var app = iris.New()

	os.Setenv("JWT_SECRET", string(jwtSecret))
	InitACConf("ac-svc", "/api/ac", "mock-ac", false)

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/ac/ping", jwt.AutoServe, auth.IsAuthenticated, RequestSrc("mock-source"), handlePing)
	e := httptest.New(t, app)
	// mock request, allow
	patches := gomonkey.ApplyFunc(restclient.Get, MockAllow)
	defer patches.Reset()

	// test normal token
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Unix() + 100,
		"uid": "abc123",
	})
	tokenString, _ := token.SignedString(jwtSecret)

	e.GET("/ac/ping").WithHeader("Authorization", "Token "+tokenString).Expect().Status(iris.StatusOK)
}

// 访问控制服务返回403，无权限
func TestRejectAC(t *testing.T) {
	var app = iris.New()

	os.Setenv("JWT_SECRET", string(jwtSecret))
	InitACConf("ac-svc", "/api/ac", "mock-ac", false)

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/ac/ping", jwt.AutoServe, auth.IsAuthenticated, RequestSrc("mock-source"), handlePing)
	e := httptest.New(t, app)
	// mock request, reject
	patches := gomonkey.ApplyFunc(restclient.Get, MockReject)
	defer patches.Reset()

	// test normal token
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Unix() + 100,
		"uid": "abc123",
	})
	tokenString, _ := token.SignedString(jwtSecret)

	e.GET("/ac/ping").WithHeader("Authorization", "Token "+tokenString).Expect().Status(iris.StatusForbidden)
}

// 访问控制服务出错例如400、500，无权限
func TestFailAC(t *testing.T) {
	var app = iris.New()

	os.Setenv("JWT_SECRET", string(jwtSecret))
	InitACConf("ac-svc", "/api/ac", "mock-ac", false)

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/ac/ping", jwt.AutoServe, auth.IsAuthenticated, RequestSrc("mock-source"), handlePing)
	e := httptest.New(t, app)
	// mock request, fail
	patches := gomonkey.ApplyFunc(restclient.Get, MockGetFail)
	defer patches.Reset()

	// test normal token
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Unix() + 100,
		"uid": "abc123",
	})
	tokenString, _ := token.SignedString(jwtSecret)

	e.GET("/ac/ping").WithHeader("Authorization", "Token "+tokenString).Expect().Status(iris.StatusForbidden)
}
