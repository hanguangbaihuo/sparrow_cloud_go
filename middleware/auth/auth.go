package auth

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	myjwt "github.com/hanguangbaihuo/sparrow_cloud_go/middleware/jwt"
	"github.com/hanguangbaihuo/sparrow_cloud_go/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

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
// JWT middleware must be configured before this
// only when your api need be authenticated, you should configure this middleware, otherwise do not configure it
func IsAuthenticated(ctx context.Context) {
	token := ctx.Values().Get(myjwt.DefaultContextKey)
	if token == nil {
		ErrorHandler(ctx, ErrTokenMissing)
		return
	}
	user, err := authenticate(ctx, token.(*jwt.Token))
	if err != nil {
		ErrorHandler(ctx, err)
		return
	}
	ctx.Values().Set(DefaultUserKey, user)
	utils.LogDebugf(ctx, "User inf is %v\n", user)
	ctx.Next()
}

func authenticate(ctx context.Context, token *jwt.Token) (User, error) {
	if token == nil {
		return User{}, ErrTokenMissing
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		utils.LogDebugf(ctx, "token is not jwt MapClaim type: %v", token.Claims)
		return User{}, ErrTokenMissing
	}
	var id string
	id, ok = claims["uid"].(string)
	if !ok {
		id = claims["id"].(string)
	}
	if id == "" {
		utils.LogInfof(ctx, "uid not found in Jwt Claim: %v\n", claims)
		return User{}, ErrUserIDMissing
	}
	return User{
		ID:              id,
		IsAuthenticated: true,
	}, nil
}
