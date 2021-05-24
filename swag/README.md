### Swag 注册文档

   将接口文档注册到swagger

#### 安装

    go get github.com/hanguangbaihuo/sparrow_cloud_go/swag

#### 使用前必知

	1. go语言中swag文档书写方法 https://github.com/swaggo/swag#general-api-info
	2. api注释说明 https://github.com/swaggo/swag#general-api-info
	3. 直接拷贝命令代码至根目录下，运行即可。注意修改配置为你的服务名称

#### 文档书写示例1

    // @Summary 接口用途，列表页接口后展示各字段
    // @Description ### 接口格式请求格式说明<br> 
    // @Description     参数：<br> 
    // @Description         {<br> 
    // @Description             "name": "hpa_name", //hpa的名字<br> 
    // @Description             "min": 1,<br> 
    // @Description             "max": 3<br> 
    // @Description         }<br> 
    // @Description     返回：200<br>
    // @Description         {<br> 
    // @Description             "message":"ok"<br> 
    // @Description         }<br>
    // @Router /api/sparrow_test/justtest/create [post]
    func processReuqest (iris.Context) {
        //do something
    }

#### 文档书写示例2

    // @Summary 接口用途，列表页接口后展示各字段
    // @Description.markdown processReuqest.md
    // @Router /api/sparrow_test/justtest/create [post]
    func processReuqest (iris.Context) {
        //do something
    }
    使用该种方式需要配置接口文档路径：cfg.MarkdownFilesDir，其中processReuqest.md文件必须放在该目录下

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
        cfg.MarkdownFilesDir = "./dir/" //该配置是当使用Description.markdown注释时，寻找md文件的公共路径，接口文档需要放在该目录下，是相对于命令文件所在位置的相对路径
		swagcfg := swag.ServiceConfig{
			"sparrow-schema-svc.frontend:8001", //此处是swagger服务的名称
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
	http_proxy=http://12.34.56.78:8888 go run -mod=vendor cmd.go