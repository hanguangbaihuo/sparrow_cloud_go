### 发送揽月app消息

### 安装

    go get github.com/hanguangbaihuo/sparrow_cloud_go/

### 注意

    必须配置环境变量:
    SC_LY_MESSAGE
    SC_LY_MESSAGE_API

#### 发送消息示例

    import (
        "github.com/hanguangbaihuo/sparrow_cloud_go/lanyue"
    )

    func main() {
        // 发送文本消息
        // 要发送的订阅消息code是sparrow_task，发送消息的服务是SparrowCloudGo，发送text消息this is a 测试 消息！
        err := lanyue.SendMsg(map[string]interface{}{"content": "this is a 测试 消息！"}, "sparrow_task", "text", "SparrowCloudGo")
        if err != nil {
            // handle error
        }

        // 发送图片消息
        // 要发送的订阅消息code是sparrow_task，发送消息的服务是SparrowCloudGo，发送image消息，图片url是https://oss.test.com/test.png
        err = lanyue.SendMsg(map[string]interface{}{"url": "https://oss.test.com/test.png"}, "sparrow_task", "image", "SparrowCloudGo")
        if err != nil {
            // handle error
        }
    }

#### 发送消息函数参数说明

    SendMsg(msg map[string]interface{}, codeType string, contentType string, msgSender string, kwargs ...map[string]interface{}) error

    msg: 发送的消息内容
    codeType: 要发送的订阅消息code
    contentType: 发送消息的类型
    msgSender: 揽月app中展示的发送消息服务的名称，一般可以取当前发送服务的名字作为该参数的值
    kwargs: 可选参数，可以添加专柜ID：shop_id，不传默认为空字符串

#### 发送msg和contentType的关系

    1.文本消息
    contentType:
        "text"
    msg:
        "content": "这是一条测试operationWarn的消息推送通知"

    2.图片消息
    contentType:
        "image"
    msg:
        "url": "https://oss.test.com/test.png"

    3.订单消息
    contentType:
        "order"
    msg:
        "order_id": 123,
        "address": "北京西单汉光百货",
        "key": "value"