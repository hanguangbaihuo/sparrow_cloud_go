package authorization

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"crypto/md5"

	"github.com/hanguangbaihuo/sparrow_cloud_go/cache"
	"github.com/hanguangbaihuo/sparrow_cloud_go/restclient"
)

var (
	ctx = context.Background()
)

type TokenData struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

func GetAppToken(svcName string, svcSecret string) (string, error) {
	key := getAppKey(svcSecret)
	tokenCache := cache.GetOrNil()
	// 若配置redis缓存，先从缓存中获取
	if tokenCache != nil {
		value, err := tokenCache.Get(ctx, key).Result()
		if err != nil {
			log.Printf("get app token from cache is %s, set it later\n", err)
		} else {
			return value, nil
		}
	}
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
	// 若配置redis缓存，则将结果缓存
	if tokenCache != nil {
		timeout := tokenData.ExpiresIn - 120
		if err := tokenCache.SetEX(ctx, key, tokenData.Token, time.Duration(timeout)*time.Second).Err(); err != nil {
			log.Printf("setex app token to cache err %s\n", err)
		}
	}
	return tokenData.Token, nil
}

func GetUserToken(svcName string, svcSecret string, userID string) (string, error) {
	key := getUserKey(userID)
	tokenCache := cache.GetOrNil()
	// 若配置redis缓存，先从缓存中获取
	if tokenCache != nil {
		value, err := tokenCache.Get(ctx, key).Result()
		if err != nil {
			log.Printf("get user token from cache is %s, set it later\n", err)
		} else {
			return value, nil
		}
	}
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
	// 若配置redis缓存，则将结果缓存
	if tokenCache != nil {
		timeout := tokenData.ExpiresIn - 120
		if err := tokenCache.SetEX(ctx, key, tokenData.Token, time.Duration(timeout)*time.Second).Err(); err != nil {
			log.Printf("setex user token to cache err %s\n", err)
		}
	}
	return tokenData.Token, nil
}

func getAppKey(svcSecret string) string {
	return strings.ToUpper("APP_TOKEN_" + getKey(svcSecret))
}

func getUserKey(userID string) string {
	return strings.ToUpper("USER_TOKEN_" + getKey(userID))
}

// getKey: data is svcSecret in AppToken, userID in UserToken
func getKey(data string) string {
	sign := md5.Sum([]byte(data))
	signStr := fmt.Sprintf("%x", sign)
	return signStr[:7]
}
