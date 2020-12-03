### 发送通知消息（钉钉或微信）

### 安装

    go get github.com/hanguangbaihuo/sparrow_cloud_go/

### 注意

    必须配置环境变量:
    SC_MESSAGE_ROBOT
    SC_MESSAGE_ROBOT_API

#### 发送消息示例

    import (
        "github.com/hanguangbaihuo/sparrow_cloud_go/authorization"
	    "github.com/hanguangbaihuo/sparrow_cloud_go/robot"
    )

    func main() {
        // 1.先获取服务认证token
        // 注意修改函数参数为自己服务的名称和服务注册密钥
        // 同时配置authorization所需要的环境变量
        token, err := authorization.GetAppToken("YourServiceName", "ServiceSecret")
        if err != nil {
            // handle error
        }
        // 2.发送消息通知
        err = robot.SendMsg("test", []string{"backend"}, "dingtalk", "text", token)
        if err != nil {
            // 发送失败，查看返回错误
        }
    }

#### 发送消息函数参数说明

    SendMsg(msg string, codeList []string, channel string, msgType string, token string) error

    msg: 发送的消息内容
    codeList: 消息群code
    channel: 消息发送的渠道，可选为("wechat", "dingtalk")两种
    msgType: 消息类型，钉钉只支持("text"), 微信支持("text", "markdown")