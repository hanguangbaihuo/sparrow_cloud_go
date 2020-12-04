# sparrow_cloud_go
基于 Iris 微服务框架

### 测试运行 ###

    运行所有测试:
        go test ./...
    运行单个测试:
        go ....

### Iris middleware ###

* JWT Middleware : 解析 JWT Token [readme](/middleware/jwt/README.md)
* AUTH Middleware: 认证用户ID   [readme](/middleware/auth/README.md)
* AccessControl Middleware: 访问控制 [readme](/middleware/accesscontrol/README.md)
* Opentracing Middleware: 追踪链中间件，配合restclient使用，追踪链注入由envoy完成 [readme](/middleware/opentracing/README.md)

### restclient ###

    跨服务间请求调用 [readme](/restclient/README.md)

### authorization ###

    获取访问认证token [readme](/authorization/README.md)

### rabbitmq ###

    发送异步消息和异步延时消息 [readme](/rabbitmq/README.md)

### robot ###

    发送通知消息（钉钉和微信）[readme](/robot/README.md)

### swag ###

    swagger文档注册 [readme](/swag/README.md)