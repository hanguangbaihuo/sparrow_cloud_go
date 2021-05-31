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
	    fmt.Println(string(res.Body), res.Code)
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
	    fmt.Println(string(res.Body), res.Code)
    }

#### IRIS中请求头传递示例

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
	    kwargs := map[string]interface{}{"timeout":5000, "headers":ctx.Request().Header}
	    res, err := restclient.Post(serviceAddr, apiPath, data, kwargs)
	    if err != nil {
	    // do something
	    }
	    fmt.Println(string(res.Body), res.Code)
    }

#### restclient中的方法

	目前该package共有5个方法，分别是Get,Post,Put,Patch,Delete方法
	每个方法的参数和返回值完全一样，是XXX(serviceAddr string, apiPath string, payload interface{}, kwargs ...map[string]interface{}) (Response, error)
	参数用途：
	serviceAddr：跨服务调用的服务地址，格式类似"sparrow-product-svc:8001"或者"127.0.0.1:8000"
	apiPath：请求服务的api，格式类似"/api/sparrow_products/products/create/"
	payload：请求服务接口所需要的数据
	kwargs：用来或者添加一些额外信息，见下方的kwargs详解
	
##### 方法参数中kwargs详解
	
	kwargs类型是map[string]interface{}
	其中参数包括以下几部分：
	timeout：建立连接和发送接收数据超时设置，不填写默认为10秒，时间单位为毫秒
	protocol：默认为"http"，构建url所用
	headers：要传递的请求头信息，数据格式是http.Header类型
	Content-Type：默认为"application/json"
	Accept： 默认为"application/json"
	Authorization：添加到请求头中的Authorization，如果设置，则请求头中的Authorization为用户设置的字符串
	token：服务内部请求认证的token，如果只有一个token字符串，则会设置为"token "+token，默认为空
	
	举例：
	kwargs := map[string]interface{}{"timeout":10000}
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
		
#### 配合链路追踪

需先在iris应用中添加 [追踪链中间件](/middleware/opentracing/README.md)
