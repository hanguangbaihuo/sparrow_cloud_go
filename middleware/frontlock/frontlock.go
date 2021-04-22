package frontlock

import (
	"encoding/json"
	"os"

	"github.com/hanguangbaihuo/sparrow_cloud_go/restclient"
	"github.com/hanguangbaihuo/sparrow_cloud_go/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

var (
	distributeLockSvc      = os.Getenv("SC_SPARROW_DISTRIBUTED_LOCK_SVC")
	distributeLockFrontApi = os.Getenv("SC_SPARROW_DISTRIBUTED_LOCK_FRONT_API")
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  int    `json:"result"`
}

func CheckLock(ctx context.Context) {
	key := ctx.GetHeader("Sc-Lock")
	if key == "" {
		ctx.Next()
		return
	}
	data := struct {
		Key string `json:"key"`
		Opt string `json:"opt"`
	}{
		Key: key,
		Opt: "incr",
	}
	res, err := restclient.Put(distributeLockSvc, distributeLockFrontApi, data)
	if err != nil {
		utils.LogWarnf(ctx, "[FrontLock] check front lock failed: %s\n", err.Error())
		ctx.Next()
		return
	}
	if res.Code != 200 {
		utils.LogErrorf(ctx, "[FrontLock] lock server response status code %d, body %s\n", res.Code, res.Code)
		ctx.Next()
		return
	}
	var response Response
	err = json.Unmarshal(res.Body, &response)
	if err != nil {
		utils.LogError(ctx, "[FrontLock] unmarshal response occur error: %s\n", err.Error())
		ctx.Next()
		return
	}
	if response.Code != 0 {
		ctx.JSON(iris.Map{"message": "重复提交，本次操作被禁止", "code": 233402})
		return
	}
	ctx.Next()
}

func UpdateLock(ctx context.Context) {
	key := ctx.GetHeader("Sc-Lock")
	if key == "" {
		return
	}
	if ctx.GetStatusCode() >= 200 && ctx.GetStatusCode() < 300 {
		data := struct {
			Key string `json:"key"`
		}{
			Key: key,
		}
		_, _ = restclient.Delete(distributeLockSvc, distributeLockFrontApi, data)
	} else {
		data := struct {
			Key string `json:"key"`
			Opt string `json:"opt"`
		}{
			Key: key,
			Opt: "reset",
		}
		_, _ = restclient.Put(distributeLockSvc, distributeLockFrontApi, data)
	}
}
