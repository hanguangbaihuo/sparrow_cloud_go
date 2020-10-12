#### JWTMiddleware ####

> 描述：Token 验证和解析
> 配置 JWTMiddleware 中间件需要的参数


```
app := iris.New()
jwtMiddleware := jwt.DefaultJwtMiddleware("your_jwt_secret")
// 自定义jwt中间件配置
// jwtMiddleware := jwt.New(jwt.Config{
//      ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
// 	    return []byte("jwt_secret"), nil
// 	    },
// 	    SigningMethod: jwt.SigningMethodHS256,
//})

// 全局添加中间件
app.Use(jwtMiddleware.Serve)
```

#### 注意

    jwt.DefaultJwtMiddleware("your_jwt_secret")
    上面的默认生成JWT中间件只会对携带jwt token的数据进行验证，
    如果token过期或者解析无效则直接返回错误
    如果没有携带token，则直接放过。
    因此，如果用户的接口需要认证，还需要在接口中添加auth中间件认证。详见: github.com/hanguangbaihuo/sparrow_cloud_go/middleware/auth