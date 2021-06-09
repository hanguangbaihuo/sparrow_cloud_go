package accesscontrol

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/hanguangbaihuo/sparrow_cloud_go/middleware/auth"
	"github.com/hanguangbaihuo/sparrow_cloud_go/restclient"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/httptest"
)

var (
	payload = map[string]interface{}{"uid": "abc123", "exp": 1722200316, "iat": 1622193116, "app_id": "core"}
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

	os.Setenv("SC_ACCESS_CONTROL_SVC", "ac-svc:8001")
	os.Setenv("SC_ACCESS_CONTROL_API", "/api/ac")
	InitACConf("mock-ac-svc", false)

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/ac/ping", auth.IsAuthenticated, RequestSrc("mock-source"), handlePing)
	e := httptest.New(t, app)
	// mock request, allow
	patches := gomonkey.ApplyFunc(restclient.Get, MockAllow)
	defer patches.Reset()

	data, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("jsom marshal payload error: %s\n", err)
		return
	}

	e.GET("/ac/ping").WithHeader("X-Jwt-Payload", string(data)).Expect().Status(iris.StatusOK)
}

// 访问控制服务返回403，无权限
func TestRejectAC(t *testing.T) {
	var app = iris.New()

	os.Setenv("SC_ACCESS_CONTROL_SVC", "ac-svc:8001")
	os.Setenv("SC_ACCESS_CONTROL_API", "/api/ac")
	InitACConf("mock-ac", false)

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/ac/ping", auth.IsAuthenticated, RequestSrc("mock-source"), handlePing)
	e := httptest.New(t, app)
	// mock request, reject
	patches := gomonkey.ApplyFunc(restclient.Get, MockReject)
	defer patches.Reset()

	data, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("jsom marshal payload error: %s\n", err)
		return
	}
	e.GET("/ac/ping").WithHeader("X-Jwt-Payload", string(data)).Expect().Status(iris.StatusForbidden)
}

// 访问控制服务出错例如400、500，无权限
func TestFailAC(t *testing.T) {
	var app = iris.New()

	os.Setenv("SC_ACCESS_CONTROL_SVC", "ac-svc:8001")
	os.Setenv("SC_ACCESS_CONTROL_API", "/api/ac")
	InitACConf("mock-ac", false)

	handlePing := func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "pong"})
	}

	app.Get("/ac/ping", auth.IsAuthenticated, RequestSrc("mock-source"), handlePing)
	e := httptest.New(t, app)
	// mock request, fail
	patches := gomonkey.ApplyFunc(restclient.Get, MockGetFail)
	defer patches.Reset()

	data, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("jsom marshal payload error: %s\n", err)
		return
	}
	e.GET("/ac/ping").WithHeader("X-Jwt-Payload", string(data)).Expect().Status(iris.StatusForbidden)
}
