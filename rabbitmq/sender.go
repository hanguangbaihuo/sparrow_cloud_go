package rabbitmq

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/hanguangbaihuo/sparrow_cloud_go/restclient"
)

type ParentOptions struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

type InputData struct {
	Code         string                 `json:"code"`
	Args         []interface{}          `json:"args"`
	Kwargs       map[string]interface{} `json:"kwargs"`
	DeliveryMode string                 `json:"delivery_mode"`
	DelayTime    int                    `json:"delay_time"`
	ParentOpt    ParentOptions          `json:"parent_options"`
}

func SendTask(msgCode string, args []interface{}, kwargs map[string]interface{}, delayTime ...int) (interface{}, error) {
	taskProxySvc := os.Getenv("SC_TASK_PROXY")
	taskProxyApi := os.Getenv("SC_TASK_PROXY_API")
	var parentOpt ParentOptions
	parentOptions := os.Getenv("SPARROW_TASK_PARENT_OPTIONS")
	if parentOptions != "" {
		err := json.Unmarshal([]byte(parentOptions), &parentOpt)
		if err != nil {
			log.Printf("json unmarshal parentOptions: %s occur error: %s", parentOptions, err)
		}
	}
	var dt int
	if len(delayTime) > 0 {
		dt = delayTime[0]
	}
	data := InputData{
		Code:         msgCode,
		Args:         args,
		Kwargs:       kwargs,
		DeliveryMode: "persistent",
		DelayTime:    dt,
		ParentOpt:    parentOpt,
	}
	res, err := restclient.Post(taskProxySvc, taskProxyApi, data)
	if err != nil {
		log.Printf("send rabbitmq message occur error: %s\n", err)
		return "", err
	}
	if res.Code != 200 {
		log.Printf("send rabbitmq message occur error, code %v, body %v\n", res.Code, string(res.Body))
		return "", errors.New(string(res.Body))
	}
	taskData := struct {
		TaskID interface{} `json:"task_id"`
	}{}
	err = json.Unmarshal(res.Body, &taskData)
	if err != nil {
		log.Printf("json unmarshal task id: %s occur error %s", string(res.Body), err)
		return "", err
	}
	return taskData.TaskID, nil
}
