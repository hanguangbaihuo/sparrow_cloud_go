package jwt

// This middleware was cloned from : https://github.com/iris-contrib/middleware/tree/v12/jwt
// we need change the default TokenExtractor
// 我们从 https://github.com/iris-contrib/middleware/tree/v12/jwt 克隆这个项目,
// 我们需要修改默认的 TokenExtractor

import (
    "testing"

    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/context"
    "github.com/kataras/iris/v12/httptest"
)

type Response struct {
    Text string `json:"text"`
}

func TestBasicJwt(t *testing.T) {
    var (
        app = iris.New()
        j   = New(Config{
            ValidationKeyGetter: func(token *Token) (interface{}, error) {
                return []byte("My Secret"), nil
            },
            SigningMethod: SigningMethodHS256,
        })
    )

    securedPingHandler := func(ctx context.Context) {
        userToken := j.Get(ctx)
        var claimTestedValue string
        if claims, ok := userToken.Claims.(MapClaims); ok && userToken.Valid {
            claimTestedValue = claims["foo"].(string)
        } else {
            claimTestedValue = "Claims Failed"
        }

        response := Response{"Iauthenticated" + claimTestedValue}
        // get the *Token which contains user information using:
        // user:= j.Get(ctx) or ctx.Values().Get("jwt").(*Token)

        ctx.JSON(response)
    }

    app.Get("/secured/ping", j.Serve, securedPingHandler)
    e := httptest.New(t, app)

    e.GET("/secured/ping").Expect().Status(iris.StatusUnauthorized)

    // Create a new token object, specifying signing method and the claims
    // you would like it to contain.
    token := NewTokenWithClaims(SigningMethodHS256, MapClaims{
        "foo": "bar",
    })

    // Sign and get the complete encoded token as a string using the secret
    tokenString, _ := token.SignedString([]byte("My Secret"))

    // e.GET("/secured/ping").WithHeader("Authorization", "Bearer "+tokenString).
    e.GET("/secured/ping").WithHeader("Authorization", "Token "+tokenString).
        Expect().Status(iris.StatusOK).Body().Contains("Iauthenticated").Contains("bar")
}
