### restclient使用方法

#### 安装

    go get github.com/hanguangbaihuo/sparrow_cloud_go/
    
#### 发送Get请求示例

	import "github.com/hanguangbaihuo/sparrow_cloud_go/restclient"
	
    func processGet(ctx iris.Context) {
	    serviceAddr := "sparrow-product-svc:8001"
	    // serviceAddr := "127.0.0.1:8001"
	    apiPath := "/api/sparrow_products/products/show/"
	    res, err := restclient.Get(serviceAddr, apiPath, nil)
	    if err != nil {
	    // do something
	    }
	    fmt.Println(res.Body, res.Code)
    }

#### 发送Post请求示例
	
	import "github.com/hanguangbaihuo/sparrow_cloud_go/restclient"
	
	type Data struct { 
		Name string  `json:"name"` 
		Price int  `json:"price"` 
		Offer int  `json:"offer"` 
	}
	
    func processPost(ctx iris.Context) {
	    serviceAddr := "sparrow-product-svc:8001"
	    apiPath := "/api/sparrow_products/products/create/"
	    data := Data{"test", 99, 0}
	    res, err := restclient.Post(serviceAddr, apiPath, data)
	    if err != nil {
	    // do something
	    }
	    fmt.Println(res.Body, res.Code)
    }

#### restclient中的方法

	目前该package共有5个方法，分别是Get,Post,Put,Patch,Delete方法
	每个方法的参数和返回值完全一样，是XXX(serviceAddr string, apiPath string, payload interface{}, kwargs ...map[string]string) (Response, error)
	参数用途：
	serviceAddr：跨服务调用的服务地址，格式类似"sparrow-product-svc:8001"或者"127.0.0.1:8000"
	apiPath：请求服务的api，格式类似"/api/sparrow_products/products/create/"
	payload：请求服务接口所需要的数据
	kwargs：用来或者添加一些额外信息，见下方的kwargs详解
	
##### 方法参数中kwargs详解
	
	kwargs类型是map[string]string
	其中参数包括以下几部分：
	timeout：建立连接和发送接收数据超时设置，默认为10s
	protocol：默认为"http"，构建url所用
	Content-Type：默认为"application/json"
	Accept： 默认为"application/json"
	operationname：追踪链中的操作名称，用来识别本次跨服务调用的用途，如果未设置，则默认为目标url
	Authorization：添加到请求头中的Authorization，如果设置，则请求头中的Authorization为用户设置的字符串；
				如果只有一个token字符串，则会设置为"token "+token，默认为空
	
	举例：
	kwargs := map[string]string{"timeout":"10","operationname":"create_product"}
	res, err := restclient.Get(serviceAddr, apiPath, nil, kwargs)

#### 方法返回Response
	
	返回的结构体如下：
	type  Response  struct {
		Body []byte
		Code int
	}
	Body是返回的数据
	Code是返回的状态码

#### 本地代理

	直接在环境变量添加http_proxy
	例如http_proxy="http://12.34.56.78:8888" go run main.go
		
#### 追踪链使用方法

	// 先在iris的应用中添加追踪链中间件,见github.com/hanguangbaihuo/sparrow_cloud_go/middleware/opentracing/
	
    import "github.com/hanguangbaihuo/sparrow_cloud_go/restclient"
    func main() {
	    // 追踪链配置：初始化Jaeger为GlobalTracer
	    closer := opentracing.InitGlobalTracer("YourServiceName")
	    defer closer.Close()
	    // 初始化iris app
	    app := iris.New()
	    // 使用opentracing中间件，从header中提取父span，并存储至中间件
	    app.Use(opentracing.Serve("YourServiceName"))
	    ...
	    app.Get("/test", processRequest)
	    app.Listen()
    }
    
    func processRequest(ctx iris.Context) {
	    serviceAddr := "sparrow-product-svc:8001"
	    apiPath := "/api/sparrow_products/products/show/"
	    res, err := restclient.Get(serviceAddr, apiPath, nil)
	    if err != nil {
	    // do something
	    }
	    fmt.Println(res.Body, res.Code)
    }