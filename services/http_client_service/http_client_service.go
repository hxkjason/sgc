package http_client_service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
		RequestBody []byte
		RetryTimes  int
		NotUseHttp2 bool
		Debug       bool
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
	var resp *http.Response
	var err error

	if rb.NotUseHttp2 {
		if err = os.Setenv("GODEBUG", "http2client=0"); err != nil {
			return nil, errors.New("set use http1 err:" + err.Error())
		}
	}

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
		request, err = http.NewRequest(rb.HttpMethod, rb.RequestUrl, bytes.NewBuffer(rb.RequestBody))
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

	retryTimes := 1
	for {
		requestTime := time.Now()
		resp, err = client.Do(request)
		if rb.Debug {
			fmt.Println("HttpClientRequestCost:", retryTimes, time.Now().Sub(requestTime))
		}
		if err == nil || (err != nil && retryTimes >= rb.RetryTimes) {
			break
		}
		resp.Close = true
		retryTimes++
	}

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
