### 初始化redis，获取Cache

### 安装

    go get github.com/hanguangbaihuo/sparrow_cloud_go/

### 注意

    使用Cache时，必须要先初始化，否则获取时出错
    使用的redis包是github.com/go-redis/redis/v8中的Client
    github地址：https://github.com/go-redis/redis
    包使用帮助文档地址：https://pkg.go.dev/github.com/go-redis/redis/v8#pkg-index

### 初始化及简单使用

    import (
        "context"
	    "github.com/hanguangbaihuo/sparrow_cloud_go/cache"
    )

    func main() {
        // 初始化redis，第一个参数为redis的链接地址，第二个为认证密码，第三个为所选的db
        c := cache.InitCache("127.0.0.1:6379", "password", 0)
        // 在使用完成后关闭redis
        defer c.Close()
    }

    // 在其他的文件或者方法中使用cache。具体各种操作的使用方法详见redis包中的Client帮助文档
    func otherfunction() {
        var ctx = context.Background()
        c := cache.Get()
        err := c.Set(ctx, "key", "value", 0).Err()
        if err != nil {
            // handle error or ignore error
        }
        val, err := c.Get(ctx, "key").Result()
        if err != nil {
            // do something
        }
    }

### 初始化方法

    简单用法
    InitCache(redisAddr, redisPasswd string, redisDb int)
    redisAddr：redis的链接地址
    redisPasswd：redis的认证密码
    redisDb：选用的redis数据库
    返回redis的Client链接符*redis.Client

    自定义用法
    InitCustomCache(opt redis.Options) *redis.Client
    opt: 自定义一些redis的参数

### 获取Cache链接符

    Get()
    在没有初始化时报错
    初始化后，调用该方法返回redis的Client链接符*redis.Client

    GetOrNil()
    在没有初始化时返回nil
    初始化后，调用该方法返回redis的Client链接符*redis.Client