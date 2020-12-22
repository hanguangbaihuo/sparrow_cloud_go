package lanyue

import (
	"errors"
	"log"
	"os"

	"github.com/hanguangbaihuo/sparrow_cloud_go/restclient"
)

type Msg struct {
	ContentType string                 `json:"content_type"`
	Data        map[string]interface{} `json:"data"`
}

type InputData struct {
	Msg       Msg    `json:"msg"`
	CodeType  string `json:"code_type"`
	MsgSender string `json:"msg_sender"`
	ShopID    string `json:"shop_id"`
}

func SendMsg(msg map[string]interface{}, codeType string, contentType string, msgSender string, kwargs ...map[string]interface{}) error {
	lanyueMsgSvc := os.Getenv("SC_LY_MESSAGE")
	lanyueMsgApi := os.Getenv("SC_LY_MESSAGE_API")
	kwarg := make(map[string]interface{})
	if len(kwargs) > 0 {
		kwarg = kwargs[0]
	}
	var shopID string
	shopID, _ = kwarg["shop_id"].(string)
	data := InputData{
		Msg: Msg{
			ContentType: contentType,
			Data:        msg,
		},
		CodeType:  codeType,
		MsgSender: msgSender,
		ShopID:    shopID,
	}
	res, err := restclient.Post(lanyueMsgSvc, lanyueMsgApi, data, kwargs...)
	if err != nil {
		log.Printf("send lanyue message occur error: %s\n", err)
		return err
	}
	if res.Code != 200 {
		log.Printf("send lanyue message occur error, code %v, body %v\n", res.Code, string(res.Body))
		return errors.New(string(res.Body))
	}
	return nil
}
