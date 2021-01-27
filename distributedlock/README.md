### 分布式锁用于加锁或解锁，防止重复提交等问题

### 安装

    go get github.com/hanguangbaihuo/sparrow_cloud_go/

### 注意

    必须配置环境变量:
    SC_SPARROW_DISTRIBUTED_LOCK_SVC
    SC_SPARROW_DISTRIBUTED_LOCK_API

### 获取锁和移除锁示例

    import (
        "fmt"
        dl "github.com/hanguangbaihuo/sparrow_cloud_go/distributedlock"
    )

    func main() {
        // 加锁
        res, err := dl.AddLock("sparrow_cloud_go", "abcde", "testuserid", 100)
        // 解锁
        // res, err := dl.RemoveLock("sparrow_cloud_go", "abcde", "testuserid")
        if err != nil {
            // 此处一般是加锁或者移除锁失败，业务层需要针对返回的错误异常情况进行处理。
        }
        if res.Code == nil {
            // 此处一般是请求参数出错或者服务器发生错误，需要打印出res.Message查看具体出错信息
            return
        }
        if *res.Code !=0 {
            // 说明此次加锁或者解锁失败
        }
        // *res.Code == 0表示加锁或者解锁成功
        fmt.Printf("code %v, body %v\n", *res.Code, res.Message)
        fmt.Printf("error %v\n", err)
    }

### 加锁函数参数说明

    AddLock(serviceName string, secret string, key string, expireTime int, kwargs ...map[string]interface{}) (Reply, error)

    serviceName: 发送方服务的名字
    secret: 发送方服务的密钥
    key: 加锁的键名
    expireTime: 键的过期时间
    kwargs: 可选参数，用于restclient包的可选参数，可以设置跨服务调用超时时间等参数，详见restclient包

### 解锁函数参数说明

    RemoveLock(serviceName string, secret string, key string, kwargs ...map[string]interface{}) (Reply, error)

    serviceName: 发送方服务的名字
    secret: 发送方服务的密钥
    key: 解锁的键名
    kwargs: 可选参数，用于restclient包的可选参数，可以设置跨服务调用超时时间等参数，详见restclient包

### 函数返回结构体说明

    type Reply struct {
        Message string `json:"message"`
        Code    *int   `json:"code"`
    }

    当加锁解锁函数返回的error不为空时：
    Message是分布式锁服务返回的具体信息，*Code是服务返回的状态码

    当error为空时：
    如果Code是nil，则说明请求出现问题，Message是具体的出错信息
    如果*Code是0，则说明加锁或者解锁成功
    如果*Code不是0，则说明加锁或者解锁失败