### auth使用方法

> 描述：用户身份认证

#### 安装

    go get github.com/hanguangbaihuo/sparrow_cloud_go
    
#### 配置前提

	需要先配置JWT中间件，详见 [JWT Middleware](/middleware/jwt/README.md)

#### 使用方法
	
	import (
		...
		"github.com/hanguangbaihuo/sparrow_cloud_go/middleware/jwt"
		"github.com/hanguangbaihuo/sparrow_cloud_go/middleware/auth"
	)
	
	func main() {
	    // 初始化iris app
	    app := iris.New()
	    // 配置jwt中间件
	    jwtMiddleware := jwt.DefaultJwtMiddleware("your_jwt_secret")
		app.Use(jwtMiddleware.Serve)
	    ...
        // /test 接口需要认证才可以进行
	    app.Get("/test", auth.IsAuthenticated, processRequest)
	    app.Listen("8081")
    }

#### 获取auth User

    user := ctx.Values().Get(auth.DefaultUserKey).(auth.User)
    // user := ctx.Values().Get("user")
    fmt.Println(user.ID, user.IsAuthenticated)


#### 获取token中的其他数据

	claimInf := ctx.Values().Get(auth.DefaultClaimsKey).(map[string]interface{})
	fmt.Println(claimInf["app_id"])

#### 无需auth中间件获取User

	user := auth.CheckUser(ctx)
	fmt.Println(user.ID, user.IsAuthenticated)