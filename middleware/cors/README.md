### cors使用方法

> 描述：跨域中间件

#### 安装

    go get github.com/hanguangbaihuo/sparrow_cloud_go

#### 使用方法
	
	import (
		...
		"github.com/hanguangbaihuo/sparrow_cloud_go/middleware/cors"
	)
	
	func main() {
	    // 初始化iris app
	    app := iris.New()
	    // 配置跨域中间件，后面接口全部允许跨域
	    app.Use(cors.Serve)
        // 接口允许options预检请求方法，该行必须配置
        app.AllowMethods(iris.MethodOptions)
	    ...
        //
	    app.Get("/test", processRequest)
	    app.Listen("8081")
    }