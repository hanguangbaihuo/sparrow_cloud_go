### Swag 注册文档

   将接口文档注册到swagger

#### 安装

    go get github.com/hanguangbaihuo/sparrow_cloud_go/swag

#### 使用前必知

	1. go语言中swag文档书写方法 https://github.com/swaggo/swag#general-api-info
	2. api注释说明 https://github.com/swaggo/swag#general-api-info
	3. 直接拷贝命令代码至根目录下，运行即可。注意修改配置为你的服务名称

#### 文档书写示例

    // Process Request func godoc
    // @Summary 接口用途，用来处理外部请求并返回
    // @Description ### 接口格式请求格式说明   
    // @Description     参数：
    // @Description         {
    // @Description             "name": "hpa_name", //hpa的名字
    // @Description             "min": 1,
    // @Description             "max": 3
    // @Description         }
    // @Description     返回：200
    // @Description         {
    // @Description             "message":"ok"
    // @Description         }
    // @Accept  json
    // @Produce  json
    // @Param name body string true "hpa name"
    // @Success 200 string {"message":"ok"}
    // @Failure 400 string {"message":"error"}
    // @Failure 404 string {"message":"error"}
    // @Failure 500 string {"message":"error"}
    // @Router /api/sparrow_test/justtest/create [post]
    func processReuqest (iris.Context) {
        //do something
    }

#### 命令文件

	// cmd.go
    package main

	import (
		"fmt"
		"github.com/hanguangbaihuo/sparrow_cloud_go/swag"
	)

	func main() {
		cfg := swag.DefaultConfig()
		// cfg.OutputFlag = true //如果设置为true，则会在文件根目录下生成./docs/swagger.json文档
		swagcfg := swag.ServiceConfig{
			"sparrow-schema-svc:8001", //此处是swagger服务的名称
			"/api/schema_i/register/", //此处为swagger服务的api接口
			"YourServiceName", //该名称需要设置为你的服务的名称
		}
		err := swag.Build(cfg, swagcfg)
		if err != nil {
			fmt.Println(err)
		}
	}
	
#### 运行示例
	
	//修改服务的代理
	http_proxy=http://12.34.56.78:8888 go run --mod=vendor cmd.go