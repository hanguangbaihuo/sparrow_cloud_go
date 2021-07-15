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
    ```

#### 注意

    可选配置的环境变量：
    SC_JWT_PUBLIC_KEY：rsa签名公钥文件数据

    JWT中间件只会对携带jwt token的数据进行验证，
    如果token过期或者解析无效则直接返回错误
    如果没有携带token，则直接放过。
    因此，如果用户的接口需要认证，还需要在接口中添加auth中间件认证。详见:
[auth中间件](/middleware/auth/README.md)