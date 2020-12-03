### 获取访问token

### 安装

    go get github.com/hanguangbaihuo/sparrow_cloud_go/

#### 获取app_token

    import (
	    "github.com/hanguangbaihuo/sparrow_cloud_go/authorization"
    )

    func main() {
        // 注意修改函数参数为自己服务的名称和服务注册密钥
        token, err := authorization.GetAppToken("YourServiceName", "ServiceSecret")
        if err != nil {
            // handle error
        }
    }