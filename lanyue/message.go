package lanyue

import (
	"os"

	"github.com/hanguangbaihuo/sparrow_cloud_go/restclient"
)

type Msg struct {
	ContentType string      `json:"content_type"`
	Data        interface{} `json:"data"`
	NickName    string      `json:"nickname,omitempty"`
	Title       string      `json:"title,omitempty"`
}

type InputData struct {
	ShopID     int      `json:"shop_id,omitempty"`
	MsgSender  string   `json:"msg_sender"`
	CodeType   string   `json:"code_type"`
	UserIDList []string `json:"user_id_list,omitempty"`
	Msg        Msg      `json:"msg"`
}

func SendMsg(data interface{}, codeType string, contentType string, msgSender string, kwargs ...map[string]interface{}) (restclient.Response, error) {
	lanyueMsgSvc := os.Getenv("SC_LY_MESSAGE")
	lanyueMsgApi := os.Getenv("SC_LY_MESSAGE_API")
	kwarg := make(map[string]interface{})
	if len(kwargs) > 0 {
		kwarg = kwargs[0]
	}
	shopID, _ := kwarg["shop_id"].(int)
	userIDList, _ := kwarg["user_id_list"].([]string)
	nickName, _ := kwarg["nickname"].(string)
	title, _ := kwarg["title"].(string)

	inputdata := InputData{
		ShopID:     shopID,
		MsgSender:  msgSender,
		CodeType:   codeType,
		UserIDList: userIDList,
		Msg: Msg{
			ContentType: contentType,
			Data:        data,
			NickName:    nickName,
			Title:       title,
		},
	}
	res, err := restclient.Post(lanyueMsgSvc, lanyueMsgApi, inputdata, kwargs...)
	return res, err
	// if err != nil {
	// 	log.Printf("send lanyue message occur error: %s\n", err)
	// 	return err
	// }
	// if res.Code != 200 {
	// 	log.Printf("send lanyue message occur error, code %v, body %v\n", res.Code, string(res.Body))
	// 	return errors.New(string(res.Body))
	// }
	// log.Printf("send data: %v\n", inputdata)
	// log.Printf("res is %v, %v\n", res.Code, string(res.Body))
	// return nil
}
