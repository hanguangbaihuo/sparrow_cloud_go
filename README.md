# sparrow_cloud_go
基于 Iris 微服务框架

### 测试运行 ###

    运行所有测试:
        go test ./...
    运行单个测试:
        go ....

### Iris middleware ###

* [JWT Middleware](/middleware/jwt/README.md) : 解析 JWT Token
* [AUTH Middleware](/middleware/auth/README.md) : 认证用户ID
* [AccessControl Middleware](/middleware/accesscontrol/README.md) : 访问控制
* [Opentracing Middleware](/middleware/opentracing/README.md): 追踪链中间件，配合restclient使用，追踪链注入由envoy完成

### restclient ###

[跨服务间请求调用](/restclient/README.md)

### authorization ###

[获取访问认证token](/authorization/README.md)

### rabbitmq ###

[发送异步消息和异步延时消息](/rabbitmq/README.md#发送异步消息)
[消费异步任务](/rabbitmq/README.md#消费异步消息)

### robot ###

[发送通知消息（钉钉和微信）](/robot/README.md)

### lanyue ###

[发送揽月app消息](/lanyue/README.md)

### distributedlock ###

[分布式锁加锁及解锁](/distributedlock/README.md)

### swag ###

[swagger文档注册](/swag/README.md)