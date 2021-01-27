package distributedlock

import (
	"encoding/json"
	"os"

	"github.com/hanguangbaihuo/sparrow_cloud_go/restclient"
)

var (
	distributedlockSvc = os.Getenv("SC_SPARROW_DISTRIBUTED_LOCK_SVC")
	distributedlockApi = os.Getenv("SC_SPARROW_DISTRIBUTED_LOCK_API")
)

type LockData struct {
	UnLockData
	ExpireTime int `json:"expire_time"`
}

type UnLockData struct {
	ServiceName string `json:"service_name"`
	Secret      string `json:"secret"`
	Key         string `json:"key"`
}

type Reply struct {
	Message string `json:"message"`
	Code    *int   `json:"code"`
}

func AddLock(serviceName string, secret string, key string, expireTime int, kwargs ...map[string]interface{}) (Reply, error) {
	lockData := LockData{
		UnLockData: UnLockData{
			ServiceName: serviceName,
			Secret:      secret,
			Key:         key,
		},
		ExpireTime: expireTime,
	}
	res, err := restclient.Post(distributedlockSvc, distributedlockApi, lockData, kwargs...)
	if err != nil {
		return Reply{
			Message: string(res.Body),
			Code:    &res.Code,
		}, err
	}
	var reply Reply
	err = json.Unmarshal(res.Body, &reply)
	if err != nil {
		return Reply{
			Message: string(res.Body),
			Code:    &res.Code,
		}, err
	}
	return reply, nil
}

func RemoveLock(serviceName string, secret string, key string, kwargs ...map[string]interface{}) (Reply, error) {
	unlockData := UnLockData{
		ServiceName: serviceName,
		Secret:      secret,
		Key:         key,
	}
	res, err := restclient.Delete(distributedlockSvc, distributedlockApi, unlockData, kwargs...)
	if err != nil {
		return Reply{
			Message: string(res.Body),
			Code:    &res.Code,
		}, err
	}
	var reply Reply
	err = json.Unmarshal(res.Body, &reply)
	if err != nil {
		return Reply{
			Message: string(res.Body),
			Code:    &res.Code,
		}, err
	}
	return reply, nil
}
