package http_client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/xiaoz194/FlyXGo/src/pkg/utils/logutil"
	"net/http"
	"time"
)

type FlyXHttpClient struct {
	HttpClient *HttpClient
	Apis       map[string]ApiDef
}

const (
	NORMAL = iota
	UNIQUE_PATH
	POSTBODY
)

/*
RetryApi 带重试机制的Api方法 外部均调用此方法!!
*/
func (flyXHttpClient *FlyXHttpClient) RetryApi(options RetryApiOptions) (*http.Response, error) {
	var resp *http.Response
	var err error
	for i := 1; i <= options.RetryMax; i++ {
		switch options.Mode {
		case NORMAL:
			resp, err = flyXHttpClient.api(options.ApisName, options.JsonData, options.Headers, options.Kwargs, options.Ctx)
		case UNIQUE_PATH:
			resp, err = flyXHttpClient.unique(options.Method, options.Path, options.JsonData, options.Headers, options.Kwargs, options.Ctx)
		case POSTBODY:
			resp, err = flyXHttpClient.apiPassByBody(options.ApisName, options.Body, options.Headers, options.Ctx)
		}
		logutil.LogrusObj.Info("resp:", resp)
		if err == nil {
			return resp, nil
		} else {
			logutil.LogrusObj.Info("retry: ", i, " err: ", err.Error())
			time.Sleep(options.RetryTimeout)
		}
	}
	return nil, err
}

// private func, this func pass by jsonData
func (flyXHttpClient *FlyXHttpClient) api(apisName string, jsonData map[string]interface{}, headers map[string]string, kwargs map[string]interface{}, ctx context.Context) (*http.Response, error) {
	apiDef, ok := flyXHttpClient.Apis[apisName]
	if !ok {
		return nil, fmt.Errorf("api: %s not exist!", apisName)
	}
	path := apiDef.Path
	if kwargs != nil {
		var urlParams string
		for key, value := range kwargs {
			urlParams += fmt.Sprintf("%s=%s&", key, value)
		}
		path = apiDef.Path + "?" + urlParams
		path = path[:len(path)-1]
	}
	logutil.LogrusObj.Info("[method]:", apiDef.Method)
	logutil.LogrusObj.Info("[path]:", apiDef.Path)
	return flyXHttpClient.HttpClient.DoRequest(apiDef.Method, path, jsonData, nil, headers, ctx)
}

// private func, this func pass by custom
func (flyXHttpClient *FlyXHttpClient) unique(method string, path string, jsonData map[string]interface{}, headers map[string]string, kwargs map[string]interface{}, ctx context.Context) (*http.Response, error) {
	if kwargs != nil {
		var urlParams string
		for key, value := range kwargs {
			urlParams += fmt.Sprintf("%s=%s&", key, value)
		}
		path = path + "?" + urlParams
		path = path[:len(path)-1]
	}
	return flyXHttpClient.HttpClient.DoRequest(method, path, jsonData, nil, headers, ctx)
}

// private func, this function pass by body(io.Reader) not by jsonData！
func (flyXHttpClient *FlyXHttpClient) apiPassByBody(apisName string, body *bytes.Buffer, headers map[string]string, ctx context.Context) (*http.Response, error) {
	apiDef, ok := flyXHttpClient.Apis[apisName]
	if !ok {
		return nil, fmt.Errorf("api: %s not exist!", apisName)
	}
	path := apiDef.Path
	logutil.LogrusObj.Info("[method]:", apiDef.Method)
	logutil.LogrusObj.Info("[path]:", apiDef.Path)
	return flyXHttpClient.HttpClient.DoRequest(apiDef.Method, path, nil, body, headers, ctx)

}
