# sparrow_cloud_go
基于 Iris 微服务框架

## 测试运行 ##

    运行所有测试:
        go test ./...
    运行单个测试:
        go ....

### Iris middleware ###

* JWT Middleware : 解析 JWT Token
* AUTH Middleware: 认证用户ID
* AccessControl Middleware: 访问控制
* Opentracing Middleware: 追踪链中间件，配合restclient使用，追踪链注入由envoy完成

#### restclient ####

> 描述：跨服务间调用