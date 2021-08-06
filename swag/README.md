### Swag 注册文档

   将接口文档注册到swagger

#### 安装

    go get github.com/hanguangbaihuo/sparrow_cloud_go/swag

#### 使用前必知

    1. go语言中swag文档书写方法 https://github.com/swaggo/swag#general-api-info
    2. api注释说明 https://github.com/swaggo/swag#general-api-info
    3. 直接拷贝命令代码至根目录下，运行即可。注意修改配置为你的服务名称

#### 文档书写示例1

    // @Summary 接口用途，列表页接口后展示该字段
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

    // @Summary 接口用途，列表页接口后展示该字段
    // @Description.markdown processReuqest.md
    // @Router /api/sparrow_test/justtest/create [post]
    func processReuqest (iris.Context) {
        //do something
    }
    使用该种方式需要配置接口文档路径：cfg.MarkdownFilesDir，其中processReuqest.md文件必须放在该目录下

#### 文档书写示例3

    // @Summary 接口用途，列表页接口后展示该字段
    // @Description.markdown ./app/processReuqest.md@create api
    // @Router /api/sparrow_test/justtest/create [post]
    func processReuqest (iris.Context) {
        //do something
    }

    该方式可以指定目录，目录是针对cmd.go文件的相对路径。不指定目录则使用配置的公共路径cfg.MarkdownFilesDir；
    同时也可以指定文件中的路由数据，即@create api，如果不指定则为整个文件数据。markdown文件书写方式如下，一个markdown文件可以写多个路由的描述
    使用该种方式需要替换官方源码，在go.mod最后添加 replace github.com/swaggo/swag => github.com/hanguangbaihuo/swag v1.7.1

#### 示例3的markdown文件书写示例

    @create api
    # 这个某个接口的post描述
        参数：
            {
                "id":1,
                "name":"test"
            }
        返回：
            {
                "message":"ok"
                "code":0
            }
    @delete some router
    # 只要按照markdown格式写就可以
    ## 参数
        {
            "payload":"sssss",
            "method":"delete"
        }
    ## 返回
        {
            "code":-1,
            "message":"error"
        }

#### 命令文件swag_register.go

    // swag_register.go
    // 自动化注册文档，文件必须是该名称，不可变动
    package main

    import (
        "fmt"
        // 修改此处为你的服务的settings初始化配置，详见sparrow_iris_template模版工程文件中的settings文件夹
        _ "sparrow_your_service/settings"

        "github.com/hanguangbaihuo/sparrow_cloud_go/swag"
        "github.com/spf13/viper"
    )

    func main() {
        cfg := swag.DefaultConfig()
        // cfg.OutputFlag = true //如果设置为true，则会在文件根目录下生成./docs/swagger.json文档
        cfg.MarkdownFilesDir = "./dir/" //该配置是当使用Description.markdown注释时，寻找md文件的公共路径，接口文档需要放在该目录下，是相对于命令文件所在位置的相对路径
        swagcfg := swag.ServiceConfig{
            viper.GetString("SC_SCHEMA_SVC"), //此处是swagger服务的名称
            viper.GetString("SC_SCHEMA_API"), //此处为swagger服务的api接口
            "YourServiceName", //该名称需要设置为你的服务的名称
            []string{"waro163","倩倩"}, //项目的贡献者，便于前端查看接口文档联系项目维护人
        }
        err := swag.Build(cfg, swagcfg)
        if err != nil {
            fmt.Println(err)
        }
    }
    
#### 运行示例
    
    //修改服务的代理
    http_proxy=http://12.34.56.78:8888 go run -mod=vendor swag_register.go