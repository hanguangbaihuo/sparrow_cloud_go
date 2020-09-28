package restclient

import (
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	Body string
	Code int
}

func request(method string, serviceAddr string, apiPath string, timeout int64, payload string, kwarg map[string]string) (Response, error) {
	var protocol string
	protocol, ok := kwarg["protocol"]
	if !ok {
		protocol = "http"
	}
	destURL := buildURL(protocol, serviceAddr, apiPath)
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Duration(timeout)*time.Second) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second)) //设置发送接收数据超时
				return c, nil
			},
		},
	}
	body := strings.NewReader(payload)
	req, err := http.NewRequest(strings.ToUpper(method), destURL, body)
	if err != nil {
		return Response{}, err
	}
	token, ok := kwarg["token"]
	if ok {
		req.Header.Add("Authorization", "token "+token)
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
	// todo: add opentracing

	response, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer response.Body.Close()

	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Response{}, err
	}
	// fmt.Println(string(body))
	return Response{string(resBody), response.StatusCode}, nil
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

func Get(serviceAddr string, apiPath string, payload string, kwargs ...map[string]string) (Response, error) {
	kwarg, timeout := parseTimeout(kwargs...)
	return request("GET", serviceAddr, apiPath, timeout, payload, kwarg)
}

func Post(serviceAddr string, apiPath string, payload string, kwargs ...map[string]string) (Response, error) {
	kwarg, timeout := parseTimeout(kwargs...)
	return request("POST", serviceAddr, apiPath, timeout, payload, kwarg)
}

func Put(serviceAddr string, apiPath string, payload string, kwargs ...map[string]string) (Response, error) {
	kwarg, timeout := parseTimeout(kwargs...)
	return request("PUT", serviceAddr, apiPath, timeout, payload, kwarg)
}

func Patch(serviceAddr string, apiPath string, payload string, kwargs ...map[string]string) (Response, error) {
	kwarg, timeout := parseTimeout(kwargs...)
	return request("PATCH", serviceAddr, apiPath, timeout, payload, kwarg)
}

func Delete(serviceAddr string, apiPath string, payload string, kwargs ...map[string]string) (Response, error) {
	kwarg, timeout := parseTimeout(kwargs...)
	return request("DELETE", serviceAddr, apiPath, timeout, payload, kwarg)
}
