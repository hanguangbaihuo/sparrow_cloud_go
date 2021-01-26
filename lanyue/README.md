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
        err := lanyue.SendMsg("this is a 测试 消息！", "sparrow_task", "text", "SparrowCloudGo")
        if err != nil {
            // handle error
        }

        // 发送图片消息
        // 要发送的订阅消息code是sparrow_task，发送消息的服务是SparrowCloudGo，发送image消息，图片url是https://oss.test.com/test.png
        err = lanyue.SendMsg("https://oss.test.com/test.png", "sparrow_task", "image", "SparrowCloudGo", map[string]interface{}{"shop_id":2,"user_id_list":[]string{"abcdef"},"nickname":"waro","title":"image card test"})
        if err != nil {
            // handle error
        }
    }

#### 发送消息函数参数说明

    SendMsg(data interface{}, codeType string, contentType string, msgSender string, kwargs ...map[string]interface{}) (restclient.Response, error)

    data: 发送的消息内容
    codeType: 要发送的订阅消息code
    contentType: 发送消息的类型,目前支持"text","image","markdown","card_text","card_image".
    msgSender: 揽月app中展示的发送消息服务的名称，一般可以取当前发送服务的名字作为该参数的值
    kwargs: 可选参数，可以添加shop_id(整型),user_id_list(字符串切片类型),nickname(字符串类型),title(字符串类型)

#### 发送数据格式和contentType的关系

    1.文本消息、图片消息、markdown消息
    contentType:
        "text"/"image"/"markdown"
    数据格式为：
        {
            "shop_id": 8, #int型，可为none
            "msg_sender": "测试wenwen",
            "code_type": "bowen_test",
            "user_id_list":["111","222"],   # 可为空列表
            "msg": {
                "content_type": "markdown",
                "data": "# xxx",
                "nickname": "宇智波悟空" #非必传，
            }
        }

    2.文本卡片、图片卡片
    contentType:
        "card_text"/"card_image"
    数据格式为：
        {
            "shop_id": 8,
            "msg_sender": "测试wenwen",
            "code_type": "bowen_test",
            "user_id_list":["111","222"],
            "msg": {
                "content_type": "card_image",
                "data": "www.xxxxxxxx.comn",
                "title": "通知"
            }
        }