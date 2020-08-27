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
        return []byte("jwt_secret"), nil
    },
    SigningMethod: jwt.SigningMethodHS256,
})
// 全局添加中间件
app.Use(jwt_middleware.Serve)
```




