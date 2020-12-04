### 发送异步消息

### 安装

    go get github.com/hanguangbaihuo/sparrow_cloud_go/

### 注意

    必须配置环境变量:
    SC_TASK_PROXY
    SC_TASK_PROXY_API

#### 发送消息示例

    import (
        "github.com/hanguangbaihuo/sparrow_cloud_go/rabbitmq"
    )

    func main() {
        // 1.发送异步消息
        // 位置参数为"test",关键字参数为hi="hello world"
        res, err := rabbitmq.SendTask("test", []interface{}{"test"}, map[string]interface{}{"hi": "hello world"})
        if err != nil {
            // handle error
        }
        id, ok := res.(string)
        if !ok {
            // 
        }

        // 2. 发送异步延时消息
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
    delayTime: 可选参数，延迟时间，

    返回：
    第一个返回数据类型是接口类型，如果需要返回的task_id，需要先进行断言。
    非延时消息返回类型是string
    延时消息返回类型是float64