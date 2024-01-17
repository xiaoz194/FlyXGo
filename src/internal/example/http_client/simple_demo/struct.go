package simple_demo

import (
	"context"
	"fmt"
	"github.com/xiaoz194/FlyXGo/src/internal/core/http"
	"github.com/xiaoz194/FlyXGo/src/internal/example/http_client/config"
	"github.com/xiaoz194/FlyXGo/src/pkg/utils/logutil"
	"time"
)

type ExampleHttpClient struct {
	FlyXHttpClient *http_client.FlyXHttpClient
}

// NewExampleHttpClient
func NewExampleHttpClient(uid int64) *ExampleHttpClient {
	Apis := map[string]http_client.ApiDef{
		"test_get":  {Method: "GET", Path: fmt.Sprintf("%s/api/v1/test_get/uid/%d", config.ExampleHttpRequestUrlPrefix, uid)},
		"test_post": {Method: "POST", Path: fmt.Sprintf("%s/api/v1/test_post/uid/%d", config.ExampleHttpRequestUrlPrefix, uid)},
	}
	return &ExampleHttpClient{
		&http_client.FlyXHttpClient{
			HttpClient: http_client.NewHTTPClient(),
			Apis:       Apis,
		},
	}
}

func (exampleHttpClient *ExampleHttpClient) TestGet() (map[string]interface{}, error) {
	kwargs := map[string]interface{}{
		"name": "sora",
	}
	// 传递上下文，设置http请求超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	options := http_client.RetryApiOptions{
		Mode:         http_client.NORMAL,
		ApisName:     "test_get",
		Path:         "",
		Method:       "",
		JsonData:     nil,
		Headers:      nil,
		Kwargs:       kwargs,
		Body:         nil,
		Ctx:          ctx,
		RetryMax:     3,
		RetryTimeout: 5 * time.Second,
	}
	resp, err := exampleHttpClient.FlyXHttpClient.RetryApi(options)
	if err != nil {
		logutil.LogrusObj.Error("Api请求出错！错误信息：", err.Error())
		return nil, err
	}
	if exampleHttpClient.FlyXHttpClient.HttpClient.CheckResponseIsOk(resp) {
		data, err := exampleHttpClient.FlyXHttpClient.HttpClient.GetJsonData(resp)
		if err != nil {
			logutil.LogrusObj.Error("GetJsonData数据解析函数出错！错误信息：", err.Error())
			return nil, err
		}
		return data, nil
	} else {
		errMsg, err := exampleHttpClient.FlyXHttpClient.HttpClient.ThrowErrorMsg(resp)
		if err != nil {
			return nil, err
		}
		logutil.LogrusObj.Error("请求目标接口响应异常，请稍后重试，响应内容：", errMsg)
		return nil, fmt.Errorf(errMsg)
	}
}

func (exampleHttpClient *ExampleHttpClient) TestPost() (map[string]interface{}, error) {
	jsonData := map[string]interface{}{
		"name":         "sora",
		"display_name": "jajaja",
		"description":  "no description here...",
	}
	headers := map[string]string{
		"XXX-XXX": "hahah",
	}
	// 传递上下文，设置http请求超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	options := http_client.RetryApiOptions{
		Mode:         http_client.NORMAL,
		ApisName:     "test_post",
		Path:         "",
		Method:       "",
		JsonData:     jsonData,
		Headers:      headers,
		Kwargs:       nil,
		Body:         nil,
		Ctx:          ctx,
		RetryMax:     3,
		RetryTimeout: 5 * time.Second,
	}

	resp, err := exampleHttpClient.FlyXHttpClient.RetryApi(options)
	if err != nil {
		logutil.LogrusObj.Error("Api请求出错！错误信息：", err.Error())
		return nil, err
	}
	if exampleHttpClient.FlyXHttpClient.HttpClient.CheckResponseIsOk(resp) {
		data, err := exampleHttpClient.FlyXHttpClient.HttpClient.GetJsonData(resp)
		if err != nil {
			logutil.LogrusObj.Error("GetJsonData数据解析函数出错！错误信息：", err.Error())
			return nil, err
		}
		return data, nil
	} else {
		errMsg, err := exampleHttpClient.FlyXHttpClient.HttpClient.ThrowErrorMsg(resp)
		if err != nil {
			return nil, err
		}
		logutil.LogrusObj.Error("请求目标接口响应异常，请稍后重试，响应内容：", errMsg)
		return nil, fmt.Errorf(errMsg)
	}
}
