## 异步消息

## 安装

    go get github.com/hanguangbaihuo/sparrow_cloud_go/

### 发送异步消息

#### 注意

    必须配置环境变量:
    SC_TASK_PROXY
    SC_TASK_PROXY_API

#### 发送消息示例

    import (
        "github.com/hanguangbaihuo/sparrow_cloud_go/rabbitmq"
    )

    func main() {
        // 1.发送异步消息
        // message code是test
        // 位置参数为"hola",关键字参数为hi="hello world"
        res, err := rabbitmq.SendTask("test", []interface{}{"hola"}, map[string]interface{}{"hi": "hello world"})
        if err != nil {
            // handle error
        }
        id, ok := res.(string)
        if !ok {
            // 
        }

        // 2. 发送异步延时消息
        // message code是test
        // 位置参数为空，关键字参数为info="this is a delay 测试 message"
        // 延时时间是3600s
        res, err = rabbitmq.SendTask("test", []interface{}{}, map[string]interface{}{"info": "this is a delay 测试 message"}, 3600)
        if err != nil {
           // handle error
        }
        num, ok := res.(float64)
        if !ok {
            // 
        }

    }

#### 发送消息函数参数说明

    SendTask(msgCode string, args []interface{}, kwargs map[string]interface{}, delayTime ...int) (interface{}, error) 

    参数：
    msgCode: message_code,消息码
    args: 位置参数，发送的异步消息数据
    kwargs: 关键字参数，发送的异步消息数据
    delayTime: 可选参数，延迟时间，当发送延时消息时，需要传递该参数

    返回：
    第一个返回数据类型是接口类型，如果需要返回的task_id，需要先进行断言。
    非延时消息返回类型是string
    延时消息返回类型是float64

### 消费异步消息

#### 注意

    必须配置的环境变量
    SC_CONSUMER_RETRY_TIMES 
    SC_CONSUMER_INTERVAL_TIME
    SC_BROKER_SERVICE_HOST
    SC_BROKER_SERVICE_PORT
    SC_BROKER_USERNAME
    SC_BROKER_PASSWORD
    SC_BROKER_VIRTUAL_HOST
    SC_BACKEND_SERVICE_SVC
    SC_BACKEND_SERVICE_API

#### 消费异步消息示例

    import (
        rq "github.com/hanguangbaihuo/sparrow_cloud_go/rabbitmq"
    )

    // 1. 消费者函数
    func testConsumerFunc(args []interface{}, kwargs map[string]interface{}) error {
        // args是接收到的位置参数，kwargs是接收到的关键字参数
        log.Printf("args %v\n", args)
        log.Printf("kwargs %v\n", kwargs)
        // return nil
        return errors.New("test failure situation")
    }

    // 2. messageCode对应的消费者函数
    var funcMap = map[string]rq.Func{
        "test7": testConsumerFunc,
    }

    func main() {
        // 3. 初始化消费者, 第一个参数是消费者队列名queue， 第二个参数是第2步中的对应关系变量
        w := rq.New("test9", funcMap)
        w.Run()
    }