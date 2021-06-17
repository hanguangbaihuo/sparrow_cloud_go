package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hanguangbaihuo/sparrow_cloud_go/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

const DefaultClaimsKey = "payload"

var (
	// ErrTokenMissing is the error value that it's returned when
	// a token is not found based on the token extractor.
	ErrTokenMissing = errors.New("this api requires authorization token")

	// ErrUserIDMissing is the error value when no user id found in jwt token
	ErrUserIDMissing = errors.New("authorization token missing user id")
)

// ErrorHandler is the default error handler.
// Use it to change the behavior for each error.
func ErrorHandler(ctx context.Context, err error) {
	if err == nil {
		return
	}
	ctx.StopExecution()
	ctx.StatusCode(iris.StatusUnauthorized)
	ctx.JSON(context.Map{"message": err.Error()})
}

// IsAuthenticated is authentication middleware
func IsAuthenticated(ctx context.Context) {
	payload := ctx.GetHeader("X-Jwt-Payload")
	user, err := authenticate(ctx, payload)
	if err != nil {
		ErrorHandler(ctx, err)
		return
	}
	ctx.Values().Set(DefaultUserKey, user)
	utils.LogDebugf(ctx, "[AUTH] User inf is %v\n", user)
	ctx.Next()
}

// CheckUser is a function to check token contains user id,
// return a User struct
func CheckUser(ctx context.Context) User {
	payload := ctx.GetHeader("X-Jwt-Payload")
	user, _ := authenticate(ctx, payload)
	return user
}

func authenticate(ctx context.Context, rawData string) (User, error) {
	if rawData == "" {
		utils.LogDebugf(ctx, "[AUTH] no X-Jwt-Payload header found\n")
		return User{}, ErrTokenMissing
	}
	var payload map[string]interface{}

	b64Payload, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		// base64解码出错，再次尝试进行文本解析
		utils.LogDebugf(ctx, "[AUTH] base64 decode fail: %s, try text decode...\n", err)
		err = json.Unmarshal([]byte(rawData), &payload)
		if err != nil {
			utils.LogErrorf(ctx, "[AUTH] can not decode X-Jwt-Payload: %v, error %s\n", rawData, err)
			return User{}, fmt.Errorf("unmarshal header error: %s", err)
		}
	} else {
		// base64解码成功
		err = json.Unmarshal(b64Payload, &payload)
		if err != nil {
			utils.LogDebugf(ctx, "[AUTH] unmarshal base64 data: %s to map type fail: %s\n", b64Payload, err)
			return User{}, fmt.Errorf("unmarshal base64 data error: %s", err)
		}
	}
	// 存储payload至中间件
	ctx.Values().Set(DefaultClaimsKey, payload)
	// 获取uid
	id, ok := payload["uid"].(string)
	if !ok || id == "" {
		utils.LogInfof(ctx, "[AUTH] uid not found in payload: %v\n", payload)
		return User{}, ErrUserIDMissing
	}
	return User{
		ID:              id,
		IsAuthenticated: true,
	}, nil
}
