package authorization

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/hanguangbaihuo/sparrow_cloud_go/restclient"
)

type TokenData struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

func GetAppToken(svcName string, svcSecret string) (string, error) {
	appManageSvc := os.Getenv("SC_MANAGE_SVC")
	appManageApi := os.Getenv("SC_MANAGE_API")
	data := struct {
		Name   string `json:"name"`
		Secret string `json:"secret"`
	}{
		Name:   svcName,
		Secret: svcSecret,
	}
	res, err := restclient.Post(appManageSvc, appManageApi, data)
	if err != nil {
		log.Printf("get app token occur error %s\n", err)
		return "", err
	}
	if res.Code != 200 {
		log.Printf("get app token occur error, code %v, body %v\n", res.Code, string(res.Body))
		return "", errors.New(string(res.Body))
	}
	var tokenData TokenData
	err = json.Unmarshal(res.Body, &tokenData)
	if err != nil {
		log.Printf("unmarshal token occur error %v\n", err)
		return "", err
	}
	return tokenData.Token, nil
}

func GetUserToken(svcName string, svcSecret string, userID string) (string, error) {
	appManageSvc := os.Getenv("SC_MANAGE_SVC")
	appManageApi := os.Getenv("SC_MANAGE_API")
	data := struct {
		Name   string `json:"name"`
		Secret string `json:"secret"`
		UserID string `json:"uid"`
	}{
		Name:   svcName,
		Secret: svcSecret,
		UserID: userID,
	}
	res, err := restclient.Post(appManageSvc, appManageApi, data)
	if err != nil {
		log.Printf("get user token occur error %s\n", err)
		return "", err
	}
	if res.Code != 200 {
		log.Printf("get user token occur error, code %v, body %v\n", res.Code, string(res.Body))
		return "", errors.New(string(res.Body))
	}
	var tokenData TokenData
	err = json.Unmarshal(res.Body, &tokenData)
	if err != nil {
		log.Printf("unmarshal token occur error %v\n", err)
		return "", err
	}
	return tokenData.Token, nil
}
