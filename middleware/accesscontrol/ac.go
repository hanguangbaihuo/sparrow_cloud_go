package accesscontrol

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hanguangbaihuo/sparrow_cloud_go/middleware/auth"
	"github.com/hanguangbaihuo/sparrow_cloud_go/restclient"
	"github.com/hanguangbaihuo/sparrow_cloud_go/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

var (
	// ErrResrouceMissing 未提供资源
	ErrResrouceMissing = errors.New("required resource not found")
	// ErrAuthMissing no auth
	ErrAuthMissing = errors.New("api needs authentication, ask developer to add auth middleware for this api")
	// ErrNoPermission user don't have resource
	ErrNoPermission = errors.New("you don't have permission to access the api")
)

// AccessControllConf is accesscontrol middleware configuration.
var AccessControllConf Config

// ErrorHandler is the default error handler.
// Use it to change the behavior for each error.
func ErrorHandler(ctx context.Context, err error) {
	if err == nil {
		return
	}

	ctx.StopExecution()
	ctx.StatusCode(iris.StatusForbidden)
	ctx.JSON(context.Map{"message": err.Error()})
}

// InitACConf constructs a new global access control configuration.
func InitACConf(acAddr string, api string, serviceName string, skipAC bool) {

	AccessControllConf = Config{acAddr, api, serviceName, skipAC}
}

// RequestSrc is access control middleware
// auth middleware must be configured before this middleware
func RequestSrc(resourceName string) func(context.Context) {
	return func(ctx context.Context) {
		// auth must be configured before run ac middleware
		// it means you must configure `auth.IsAuthenticated` before this middleware

		// TODO: check AccessControllConf had been initialized
		if AccessControllConf.AccessControlService == "" || AccessControllConf.APIPath == "" || AccessControllConf.ServiceName == "" {
			ErrorHandler(ctx, errors.New("[ERROR] Please init accesscontrol middleware configuration"))
			return
		}
		// make sure auth middleware had been configured for the route
		user, ok := ctx.Values().Get(auth.DefaultUserKey).(auth.User)
		if !ok {
			ErrorHandler(ctx, ErrAuthMissing)
			return
		}
		// skip accesscontroll
		if AccessControllConf.SkipAccessContorl {
			utils.LogDebugf(ctx, "[AC] Skip the access control\n")
			ctx.Next()
			return
		}

		apiPath := fmt.Sprintf(AccessControllConf.APIPath+apiParam, user.ID, AccessControllConf.ServiceName, resourceName)
		utils.LogDebugf(ctx, "[AC] access control service: %v, api path is %v\n", AccessControllConf.AccessControlService, apiPath)
		res, err := restclient.Get(AccessControllConf.AccessControlService, apiPath, nil)
		if err != nil {
			ErrorHandler(ctx, err)
			return
		}
		if res.Code != 200 && res.Code != 403 {
			utils.LogDebugf(ctx, "[AC] response code %v, body %v\n", res.Code, string(res.Body))
			ErrorHandler(ctx, errors.New(string(res.Body)))
			return
		}
		var acResponse ACResponse
		err = json.Unmarshal(res.Body, &acResponse)
		if err != nil {
			utils.LogInfof(ctx, "[AC] unmarsha response error: %v\n", err)
			ErrorHandler(ctx, err)
			return
		}
		if !acResponse.HasPerm {
			ErrorHandler(ctx, ErrNoPermission)
			return
		}
		// If everything ok then call next.
		ctx.Next()
	}
}
