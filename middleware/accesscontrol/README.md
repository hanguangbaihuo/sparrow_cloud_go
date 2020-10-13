### accesscontrol 访问控制中间件

    访问控制，检查用户是否拥有访问接口的资源

#### 安装

    go get github.com/hanguangbaihuo/sparrow_cloud_go
    
#### 配置前提

	需要先配置JWT中间件，详见github.com/hanguangbaihuo/sparrow_cloud_go/middleware/jwt

#### 使用前注意

    1. 必需提前配置jwt中间件
    2. 必需初始化访问控制中间件的配置
    3. 访问控制内部包含了auth中间件认证，无需在接口中再次添加。如果需要访问控制，直接添加该中间件即可
    4. 如果设置跳过访问控制，添加了访问控制的中间件仍然需要auth认证

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
		"github.com/hanguangbaihuo/sparrow_cloud_go/middleware/accesscontrol"
	)
	
	func main() {
	    // 初始化iris app
	    app := iris.New()
	    // 配置jwt中间件
	    jwtMiddleware := jwt.DefaultJwtMiddleware("your_jwt_secret")
		app.Use(jwtMiddleware.Serve)
	    ...
        // 初始化访问控制中间件配置
        accesscontrol.InitACConf("sparrow-access-control-svc:8001", "/api/ac_i/verify/", "SparrowPromotion", false)

        // /test 接口需要用户认证并拥有admin资源才可以访问
	    app.Get("/test", accesscontrol.RequestSrc("admin"), processRequest)
	    app.Listen("8081")
    }
