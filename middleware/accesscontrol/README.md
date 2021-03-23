### accesscontrol 访问控制中间件

    访问控制，检查用户是否拥有访问接口的资源

#### 安装

    go get github.com/hanguangbaihuo/sparrow_cloud_go
    
#### 配置前提

    需要先配置JWT中间件，详见 [JWT Middleware](/middleware/jwt/README.md)
    接口添加auth认证，详见 [AUTH Middleware](/middleware/auth/README.md)

#### 使用前注意

    1. 必需提前配置jwt中间件
    2. 必需初始化访问控制中间件的配置
    3. 如果接口需要访问控制，则接口前必需先添加auth认证中间件，否则接口不通过
    4. 如果配置跳过访问控制，仍然需要添加auth认证，因为之后变为不跳过访问控制后，必需用到user_id，该数据只能从auth中间件获得

#### 初始化访问控制中间件配置

    func InitACConf(acAddr string, api string, serviceName string, skipAC bool)
    参数含义：
    acAddr: 访问控制服务的服务地址，例如：sparrow-access-control-svc:8001
    api: 访问控制服务的api，例如：/api/ac_i/verify/
    serviceName: 你的服务的名字，例如：SparrowPromotion
    skipAC: 是否跳过访问控制，设置为true则跳过访问控制，但是仍然需要认证；设置为false，则不跳过访问控制

#### 使用方法
	
	import (
		...
		"github.com/hanguangbaihuo/sparrow_cloud_go/middleware/jwt"
        "github.com/hanguangbaihuo/sparrow_cloud_go/middleware/auth"
		ac "github.com/hanguangbaihuo/sparrow_cloud_go/middleware/accesscontrol"
	)
	
	func main() {
	    // 初始化iris app
	    app := iris.New()
	    // 配置jwt中间件
	    jwtMiddleware := jwt.DefaultJwtMiddleware("your_jwt_secret")
		app.Use(jwtMiddleware.Serve)
	    ...
        // 初始化访问控制中间件配置
        ac.InitACConf("sparrow-access-control-svc:8001", "/api/ac_i/verify/", "SparrowPromotion", false)

        // /test 接口需要用户认证并拥有admin资源才可以访问
	    app.Get("/test", auth.IsAuthenticated, ac.RequestSrc("admin"), processRequest)
	    app.Listen("8081")
    }
