package robot

import (
	"errors"
	"log"
	"os"

	"github.com/hanguangbaihuo/sparrow_cloud_go/restclient"
)

func SendMsg(msg string, codeList []string, channel string, msgType string, token string) error {
	robotMsgSvc := os.Getenv("SC_MESSAGE_ROBOT")
	robotMsgApi := os.Getenv("SC_MESSAGE_ROBOT_API")
	data := struct {
		Msg           string   `json:"msg"`
		GroupCodeList []string `json:"group_code_list"`
		Channel       string   `json:"channel"`
		MessageType   string   `json:"message_type"`
	}{
		msg, codeList, channel, msgType,
	}
	kwargs := map[string]interface{}{"Authorization": "token " + token}
	res, err := restclient.Post(robotMsgSvc, robotMsgApi, data, kwargs)
	if err != nil {
		log.Printf("send robot message occur error: %s\n", err)
		return err
	}
	if res.Code != 200 {
		log.Printf("send robot message occur error, code %v, body %v\n", res.Code, string(res.Body))
		return errors.New(string(res.Body))
	}
	return nil
}
