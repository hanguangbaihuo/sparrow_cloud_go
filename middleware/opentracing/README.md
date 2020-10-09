### opentracing中间件使用方法

#### 安装

    go get github.com/hanguangbaihuo/sparrow_cloud_go

#### iris框架中使用方法

    import "github.com/hanguangbaihuo/sparrow_cloud_go/middleware/opentracing"
    func main() {
	    // 追踪链配置：初始化Jaeger为GlobalTracer
        // 初始化失败则程序直接结束
	    closer := opentracing.InitGlobalTracer("YourServiceName")
	    defer closer.Close()
	    // 初始化iris app
	    app := iris.New()
	    // 使用opentracing中间件，从header中提取父span，并存储至中间件
	    app.Use(opentracing.Serve("YourServiceName"))
        ...
    }

#### 注意

	该中间件需要配合restclient包使用，才能实现追踪链
	restclient包的地址：github.com/hanguangbaihuo/sparrow_cloud_go/restclient

