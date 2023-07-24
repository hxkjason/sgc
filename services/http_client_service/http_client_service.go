package http_client_service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type (
	RequestAttrs struct {
		RequestUrl  string
		HttpMethod  string
		Params      interface{}
		QueryParams map[string]string
		Headers     map[string]string
		Timeout     time.Duration
		Result      interface{}
		RequestBody string
	}
)

func Request(rb RequestAttrs) ([]byte, error) {

	// 设置超时时间
	if rb.Timeout == 0 {
		rb.Timeout = 60 * time.Second
	}

	client := &http.Client{
		Timeout: rb.Timeout,
	}

	var request *http.Request
	var err error
	requestTime := time.Now()

	switch rb.HttpMethod {
	case http.MethodGet:
		request, err = http.NewRequest(rb.HttpMethod, rb.RequestUrl, nil)
		if err != nil {
			return nil, errors.New("参数编码失败:" + err.Error())
		}
		q := request.URL.Query()
		for k, v := range rb.QueryParams {
			q.Add(k, v)
		}
		request.URL.RawQuery = q.Encode()

	case http.MethodPost, http.MethodPut, http.MethodDelete:
		paramBytes, err := json.Marshal(rb.Params)
		if err != nil {
			return nil, errors.New("参数编码失败:" + err.Error())
		}
		request, err = http.NewRequest(rb.HttpMethod, rb.RequestUrl, bytes.NewBuffer(paramBytes))
	default:
		return nil, errors.New("当前请求方法[" + rb.HttpMethod + "]暂不支持")
	}

	if err != nil {
		return nil, errors.New("建立请求出错:" + err.Error())
	}

	// 设置请求头
	if len(rb.Headers) == 0 {
		request.Header.Set("Content-Type", "application/json")
	} else {
		for k, v := range rb.Headers {
			request.Header.Set(k, v)
		}
	}

	resp, err := client.Do(request)
	fmt.Println("requestTime:", time.Now().Sub(requestTime).String())

	if err != nil {
		if strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
			return nil, errors.New("请求超时:" + err.Error())
		}
		return nil, errors.New("请求出错:" + err.Error())
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resBody, errors.New("读取响应失败:" + err.Error())
	}

	if err = json.Unmarshal(resBody, &rb.Result); err != nil {
		return resBody, errors.New("解码响应出错:" + err.Error())
	}

	return resBody, nil
}
