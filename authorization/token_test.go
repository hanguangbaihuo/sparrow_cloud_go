package authorization

import (
	"os"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/hanguangbaihuo/sparrow_cloud_go/restclient"
	"github.com/stretchr/testify/assert"
)

func MockPostSucc(serviceAddr string, apiPath string, payload interface{}, kwargs ...map[string]interface{}) (restclient.Response, error) {
	return restclient.Response{Code: 200, Body: []byte(`{"token":"abc123","expires_in":7200}`)}, nil
}

func MockPostFail(serviceAddr string, apiPath string, payload interface{}, kwargs ...map[string]interface{}) (restclient.Response, error) {
	return restclient.Response{Code: 400, Body: []byte(`{"message":"lost parameter"}`)}, nil
}

func TestAppTokenSucc(t *testing.T) {
	os.Setenv("SC_MANAGE_SVC", "manage-svc:8001")
	os.Setenv("SC_MANAGE_API", "/api/get_token")

	patches := gomonkey.ApplyFunc(restclient.Post, MockPostSucc)
	defer patches.Reset()
	token, err := GetAppToken("MockTest", "MockSecret")
	assert.Nil(t, err)
	assert.Equal(t, "abc123", token, "token not equal")
}

func TestAppTokenFail(t *testing.T) {
	os.Setenv("SC_MANAGE_SVC", "manage-svc:8001")
	os.Setenv("SC_MANAGE_API", "/api/get_token")

	patches := gomonkey.ApplyFunc(restclient.Post, MockPostFail)
	defer patches.Reset()
	token, err := GetAppToken("MockTest", "MockSecret")
	assert.NotNil(t, err)
	assert.Equal(t, "", token, "token not empty")
}

func TestUserTokenSucc(t *testing.T) {
	os.Setenv("SC_MANAGE_SVC", "manage-svc:8001")
	os.Setenv("SC_MANAGE_API", "/api/get_token")

	patches := gomonkey.ApplyFunc(restclient.Post, MockPostSucc)
	defer patches.Reset()
	token, err := GetUserToken("MockTest", "MockSecret", "MockUserID")
	assert.Nil(t, err)
	assert.Equal(t, "abc123", token, "token not equal")
}

func TestUserTokenFail(t *testing.T) {
	os.Setenv("SC_MANAGE_SVC", "manage-svc:8001")
	os.Setenv("SC_MANAGE_API", "/api/get_token")

	patches := gomonkey.ApplyFunc(restclient.Post, MockPostFail)
	defer patches.Reset()
	token, err := GetUserToken("MockTest", "MockSecret", "MockUserID")
	assert.NotNil(t, err)
	assert.Equal(t, "", token, "token not empty")
}
