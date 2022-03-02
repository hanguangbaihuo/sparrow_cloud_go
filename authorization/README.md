### 获取访问token

### 安装

    go get github.com/hanguangbaihuo/sparrow_cloud_go/

### 注意

    必须配置环境变量:
    SC_MANAGE_SVC
    SC_MANAGE_API
    配置下面环境为true跳过缓存
    SC_SKIP_TOKEN_CACHE

#### 获取app_token

    直接返回常量字符串`{"uid":"YourServiceName"}`即可

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

#### 获取user_token

    直接返回常量字符串`{"uid":"用户的ID"}`即可

    import (
	    "github.com/hanguangbaihuo/sparrow_cloud_go/authorization"
    )

    func main() {
        // 注意修改函数参数为自己服务的名称和服务注册密钥，和用户ID
        token, err := authorization.GetUserToken("YourServiceName", "ServiceSecret", "user_id")
        if err != nil {
            // handle error
        }
    }

#### token使用

    //直接将获取到的token赋值给restclient包中函数的kwargs参数中的token，例如：
    token, err := authorization.GetAppToken("YourServiceName", "ServiceSecret")
    if err != nil {
        // handle error
    }
    kwargs := map[string]interface{}{"token":token}
    res, err := restclient.Post(serviceAddr, apiPath, data, kwargs)
    if err != nil {
    // do something
    }
    fmt.Println(string(res.Body), res.Code)