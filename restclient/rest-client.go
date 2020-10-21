package restclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hanguangbaihuo/sparrow_cloud_go/middleware/opentracing"
)

// Response is the response data of external request url
type Response struct {
	Body []byte
	Code int
}

func request(method string, serviceAddr string, apiPath string, timeout int64, payload interface{}, kwarg map[string]string) (Response, error) {
	var protocol string
	protocol, ok := kwarg["protocol"]
	if !ok {
		protocol = "http"
	}
	destURL := buildURL(protocol, serviceAddr, apiPath)
	proxyURL := getEnvProxy()
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Duration(timeout)*time.Millisecond) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond)) //设置发送接收数据超时
				return c, nil
			},
			Proxy: http.ProxyURL(proxyURL),
		},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("marshal payload occur error: %s\n", err)
		return Response{}, err
	}
	body := bytes.NewReader(data)
	req, err := http.NewRequest(strings.ToUpper(method), destURL, body)
	if err != nil {
		return Response{}, err
	}
	token, ok := kwarg["Authorization"]
	if ok {
		validToken := checkAuthorization(token)
		req.Header.Set("Authorization", validToken)
	}
	contentType, ok := kwarg["Content-Type"]
	if ok {
		req.Header.Set("Content-Type", contentType)
	} else {
		req.Header.Set("Content-Type", "application/json")
	}
	accept, ok := kwarg["Accept"]
	if ok {
		req.Header.Set("Accept", accept)
	} else {
		req.Header.Set("Accept", "application/json")
	}
	// add opentracing b3 headers
	if opentracing.OpentracingInf != nil {
		for key, value := range opentracing.OpentracingInf {
			if len(value) > 0 {
				req.Header.Set(key, value[0])
			}
		}
		// log.Printf("send %s header is %#v\n", destURL, req.Header)
	}

	response, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer response.Body.Close()

	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Response{}, err
	}
	// log.Printf("%d, %s", response.StatusCode, string(resBody))
	return Response{resBody, response.StatusCode}, nil
}

func parseTimeout(kwargs ...map[string]string) (map[string]string, int64) {
	kwarg := make(map[string]string)
	if len(kwargs) > 0 {
		kwarg = kwargs[0]
	}
	var timeout int64
	timeout = 10
	timeoutStr, ok := kwarg["timeout"]
	if ok {
		timeout, _ = strconv.ParseInt(timeoutStr, 10, 64)
	}
	return kwarg, timeout
}

func getEnvProxy() *url.URL {
	var httpProxy string
	httpProxy = os.Getenv("http_proxy")
	if httpProxy == "" {
		httpProxy = os.Getenv("HTTP_PROXY")
	}
	if httpProxy != "" {
		proxyURL, err := url.Parse(httpProxy)
		if err != nil {
			log.Printf("parse http proxy: %s occur error: %s\n", httpProxy, err)
			return nil
		}
		return proxyURL
	}
	return nil
}

func checkAuthorization(authcode string) string {
	rawCode := strings.TrimSpace(authcode)
	codes := strings.Split(rawCode, " ")
	if len(codes) > 2 { // invalid Authorization header
		return ""
	} else if len(codes) < 2 {
		if rawCode != "" { // only token without token key
			return "token " + rawCode
		}
		// if authcode is empty, len(code) also will return 1
		return ""
	} else {
		return rawCode
	}
}

func Get(serviceAddr string, apiPath string, payload interface{}, kwargs ...map[string]string) (Response, error) {
	kwarg, timeout := parseTimeout(kwargs...)
	return request("GET", serviceAddr, apiPath, timeout, payload, kwarg)
}

func Post(serviceAddr string, apiPath string, payload interface{}, kwargs ...map[string]string) (Response, error) {
	kwarg, timeout := parseTimeout(kwargs...)
	return request("POST", serviceAddr, apiPath, timeout, payload, kwarg)
}

func Put(serviceAddr string, apiPath string, payload interface{}, kwargs ...map[string]string) (Response, error) {
	kwarg, timeout := parseTimeout(kwargs...)
	return request("PUT", serviceAddr, apiPath, timeout, payload, kwarg)
}

func Patch(serviceAddr string, apiPath string, payload interface{}, kwargs ...map[string]string) (Response, error) {
	kwarg, timeout := parseTimeout(kwargs...)
	return request("PATCH", serviceAddr, apiPath, timeout, payload, kwarg)
}

func Delete(serviceAddr string, apiPath string, payload interface{}, kwargs ...map[string]string) (Response, error) {
	kwarg, timeout := parseTimeout(kwargs...)
	return request("DELETE", serviceAddr, apiPath, timeout, payload, kwarg)
}
