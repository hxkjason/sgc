package http_client_service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
		MessageId   string
		RetryTimes  int
		NotUseHttp2 bool
		Debug       bool
		SaveLog     int8
	}

	HttpResponse struct {
		Response     *http.Response
		ResBodyBytes []byte
	}
)

func RequestV1(rb RequestAttrs) (*HttpResponse, error) {

	// 设置超时时间
	if rb.Timeout == 0 {
		rb.Timeout = 60 * time.Second
	}

	client := &http.Client{
		Timeout: rb.Timeout,
	}

	var request *http.Request
	var resp *http.Response
	var res HttpResponse
	var err error

	if rb.NotUseHttp2 {
		if err = os.Setenv("GODEBUG", "http2client=0"); err != nil {
			return &res, errors.New("set use http1 err:" + err.Error())
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout*2)
	defer cancel()

	retryTimes := 1

	for {

		requestTime := time.Now()

		switch rb.HttpMethod {
		case http.MethodGet:
			request, err = http.NewRequestWithContext(ctx, rb.HttpMethod, rb.RequestUrl, nil)
		case http.MethodPost, http.MethodPut, http.MethodDelete:
			request, err = http.NewRequestWithContext(ctx, rb.HttpMethod, rb.RequestUrl, bytes.NewBuffer(rb.RequestBody))
		default:
			return &res, errors.New("当前请求方法[" + rb.HttpMethod + "]暂不支持")
		}
		if err != nil {
			return &res, errors.New("建立请求出错:" + err.Error())
		}

		if len(rb.QueryParams) > 0 {
			q := request.URL.Query()
			for k, v := range rb.QueryParams {
				q.Add(k, v)
			}
			request.URL.RawQuery = q.Encode()
		}

		// 设置请求头
		if len(rb.Headers) == 0 {
			request.Header.Set("Content-Type", "application/json")
		} else {
			for k, v := range rb.Headers {
				request.Header.Set(k, v)
			}
		}

		request.Close = true

		resp, err = client.Do(request)
		if rb.Debug {
			fmt.Println("HttpClientRequestCost:", retryTimes, time.Now().Sub(requestTime))
		}
		if err == nil || (err != nil && retryTimes >= rb.RetryTimes) {
			break
		}
		//resp.Close = true
		retryTimes++
	}

	res.Response = resp

	if err != nil {
		if strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
			return &res, errors.New("请求超时:" + err.Error())
		}
		return &res, errors.New("请求出错:" + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		resBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return &res, errors.New("读取响应失败:" + err.Error())
		}
		res.ResBodyBytes = resBody

		if err = json.Unmarshal(resBody, &rb.Result); err != nil {
			return &res, errors.New("解码响应出错:" + err.Error() + ",originResBody:" + string(resBody))
		}
	}

	return &res, nil
}

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
		if len(rb.QueryParams) > 0 {
			q := request.URL.Query()
			for k, v := range rb.QueryParams {
				q.Add(k, v)
			}
			request.URL.RawQuery = q.Encode()
		}
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

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resBody, errors.New("读取响应失败:" + err.Error())
	}

	if err = json.Unmarshal(resBody, &rb.Result); err != nil {
		return resBody, errors.New("解码响应出错:" + err.Error())
	}

	return resBody, nil
}
