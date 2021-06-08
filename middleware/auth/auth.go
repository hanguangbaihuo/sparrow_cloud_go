package auth

import (
	"errors"

	"github.com/hanguangbaihuo/sparrow_cloud_go/middleware/jwt"
	"github.com/hanguangbaihuo/sparrow_cloud_go/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

const DefaultClaimsKey = "claims"

var (
	// ErrTokenMissing is the error value that it's returned when
	// a token is not found based on the token extractor.
	ErrTokenMissing = errors.New("this api requires authorization token")

	// ErrUserIDMissing is the error value when no user id found in jwt token
	ErrUserIDMissing = errors.New("authorization token missing user id")

	// ErrTokenType when parse jwt token, its type is not available type
	ErrTokenType = errors.New("token type error")
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
// JWT middleware must be configured before this
// only when your api need be authenticated, you should configure this middleware, otherwise do not configure it
func IsAuthenticated(ctx context.Context) {
	token := ctx.Values().Get(jwt.DefaultContextKey)

	user, err := authenticate(ctx, token)
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
	token := ctx.Values().Get(jwt.DefaultContextKey)

	user, _ := authenticate(ctx, token)
	return user
}

func authenticate(ctx context.Context, token interface{}) (User, error) {
	if token == nil {
		utils.LogDebugf(ctx, "[AUTH] no token or jwt middleware configure error\n")
		return User{}, ErrTokenMissing
	}
	jtoken, ok := token.(*jwt.Token)
	if !ok {
		utils.LogDebugf(ctx, "[AUTH] token is not jwt Token type: %v\n", token)
		return User{}, ErrTokenType
	}
	claims, ok := jtoken.Claims.(jwt.MapClaims)
	if !ok {
		utils.LogDebugf(ctx, "[AUTH] token is not jwt MapClaim type: %v\n", jtoken.Claims)
		return User{}, ErrTokenType
	}
	// 存储token中所有数据
	data := make(map[string]interface{})
	for key, value := range claims {
		data[key] = value
	}
	ctx.Values().Set(DefaultClaimsKey, data)
	// 获取uid
	var id string
	id, ok = claims["uid"].(string)
	// if !ok {
	// 	id = claims["id"].(string)
	// }
	if !ok || id == "" {
		utils.LogInfof(ctx, "[AUTH] uid not found in Jwt Claim: %v\n", claims)
		return User{}, ErrUserIDMissing
	}
	return User{
		ID:              id,
		IsAuthenticated: true,
	}, nil
}
