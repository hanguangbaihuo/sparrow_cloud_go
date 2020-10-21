### opentracing中间件使用方法

#### 安装

    go get github.com/hanguangbaihuo/sparrow_cloud_go

#### iris框架中使用方法

    import "github.com/hanguangbaihuo/sparrow_cloud_go/middleware/opentracing"

    func main() {
	    // 初始化iris app
	    app := iris.New()
	    // 使用opentracing中间件，从header中b3 headers，用于之后注入restclient请求头
	    app.Use(opentracing.Serve)
        ...
    }

### iris框架中zipkin追踪链使用方法【不再使用】

	import "github.com/hanguangbaihuo/sparrow_cloud_go/middleware/opentracing"

    func main() {
	    // 追踪链配置：初始化GlobalTracer
		// 参数true是是否开启Debug模式，true会输出各种追踪链的log，生产环境请设置为false
		opentracing.InitZipkinTracer(true)

	    // 初始化iris app
	    app := iris.New()
	    // 使用opentracing中间件，从header中提取父span，并存储至全局
	    app.Use(opentracing.ZipkinServe("YourServiceName"))
        ...
    }

#### 注意

	该中间件需要配合restclient包使用，才能实现链路追踪
	restclient包的地址：github.com/hanguangbaihuo/sparrow_cloud_go/restclient

