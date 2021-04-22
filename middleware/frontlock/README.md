### frontlock使用方法

> 描述：防前端重复提交锁

#### 安装

    go get github.com/hanguangbaihuo/sparrow_cloud_go

#### 使用方法

    import (
        ...
        fl "github.com/hanguangbaihuo/sparrow_cloud_go/middleware/frontlock"
        "github.com/kataras/iris/v12"
    )

    func main(){
        app := iris.New()
        app.Logger().SetLevel("debug")

        app.Get("api/ping", func(ctx iris.Context) {
            ctx.JSON(iris.Map{"message": "ping"})
        })
        // 之后的所有API添加检查前端锁中间件，如果锁已经被使用，直接返回不再执行业务逻辑
        app.Use(fl.CheckLock)
        // 根据业务层返回的状态码对前端锁进行更新操作
        app.Done(fl.UpdateLock)
        // 强制执行更新操作
        app.SetExecutionRules(iris.ExecutionRules{
            Done: iris.ExecutionOptions{Force: true},
        })

        app.Get("api/ok", func(ctx iris.Context) {
            // ctx.StatusCode(iris.StatusNotFound)
            ctx.JSON(iris.Map{"message": "ok"})
        })
        app.Listen("0.0.0.0:8001")
    }

#### 注意

    如果使用该中间件，下面3行是必须要写的，缺一不可。
    除非如果你已经通过其他方法可以在执行业务逻辑代码后，执行Done中的中间件。否则不要更改
    app.Use(fl.CheckLock)

    app.Done(fl.UpdateLock)

    app.SetExecutionRules(iris.ExecutionRules{
        Done: iris.ExecutionOptions{Force: true},
    })