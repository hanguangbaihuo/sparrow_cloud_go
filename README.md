# sparrow_cloud_go
基于 Iris 微服务框架

## 测试运行 ##

    运行所有测试:
        go test ./...
    运行单个测试:
        go ....

### Iris middleware ###

* JWT Middleware : 解析 JWT Token


#### JWTMiddleware ####

> 描述：Token 验证和解析
> 配置 JWTMiddleware 中间件需要的参数

```
注册中间件
app := iris.New()

jwt_middleware := jwt.New(jwt.Config{
    ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
        return mySecret, nil
    },

    // Extract by the "token" url.
    // There are plenty of options.
    // The default jwt's behavior to extract a token value is by
    // the `Authorization: Bearer $TOKEN` header.
    Extractor: jwt.FromParameter("token"),
    // When set, the middleware verifies that tokens are
    // signed with the specific signing algorithm
    // If the signing method is not constant the `jwt.Config.ValidationKeyGetter` callback
    // can be used to implement additional checks
    // Important to avoid security issues described here:
    // https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
    SigningMethod: jwt.SigningMethodHS256,
})
app.Use(jwt_middleware)
```




