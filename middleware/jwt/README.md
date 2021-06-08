#### JWTMiddleware ####

    > 描述：Token 验证和解析
    > 配置 JWTMiddleware 中间件需要的参数

    ```
    import (
        "github.com/hanguangbaihuo/sparrow_cloud_go/middleware/jwt"
    )

    app := iris.New()
    // 全局添加中间件
    app.Use(jwt.AutoServe)

    // 下列使用方式将会废弃
    // jwtMiddleware := jwt.DefaultJwtMiddleware("your_jwt_secret")
    // app.Use(jwtMiddleware.Serve)
    ```

#### 注意

    必须配置的环境变量：
    JWT_SECRET：jwt密钥
    PUBLIC_KEY_PATH：rsa签名公钥文件路径

    JWT中间件只会对携带jwt token的数据进行验证，
    如果token过期或者解析无效则直接返回错误
    如果没有携带token，则直接放过。
    因此，如果用户的接口需要认证，还需要在接口中添加auth中间件认证。详见:
[auth中间件](/middleware/auth/README.md)