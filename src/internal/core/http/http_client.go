package http_client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/xiaoz194/FlyXGo/src/internal/example/http_client/config"
	"github.com/xiaoz194/FlyXGo/src/pkg/utils/logutil"
	"io"
	"net/http"
	"time"
)

type ApiDef struct {
	Method string
	Path   string
}

type HttpClient struct {
	Client  *http.Client
	Headers map[string]string
}

// NewHTTPClient 的构造函数
func NewHTTPClient() *HttpClient {
	// 创建支持 TLS 的客户端
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify, // 是否跳过证书验证，true 跳过，不推荐在生产环境中使用 false不跳过
		},
		Proxy:                 http.ProxyFromEnvironment, // 使用系统代理
		MaxIdleConns:          config.MaxIdleConns,
		IdleConnTimeout:       config.IdleConnTimeout * time.Second,
		ResponseHeaderTimeout: config.ResponseHeaderTimeout * time.Second,
	}
	client := &http.Client{Transport: tr}
	return &HttpClient{
		Client:  client,
		Headers: make(map[string]string),
	}
}

func (c *HttpClient) GetContentTypeFromHeaders() string {
	return c.Headers["Content-Type"]
}

func (c *HttpClient) SetContentTypeWithHeaders(contentType string) {
	c.Headers["Content-Type"] = contentType
}

// DoRequest 封装 发送 HTTP 请求的基础方法
func (c *HttpClient) DoRequest(method string, path string, jsonData map[string]interface{}, bodyData *bytes.Buffer, headers map[string]string, ctx context.Context) (*http.Response, error) {
	logutil.LogrusObj.Infof("do request,request path is:%s", path)

	var body []byte
	var req *http.Request
	var err error
	var resp *http.Response
	// 传入的数据是bodyData 走该分支逻辑
	if bodyData != nil {
		return c.do(req, resp, err, c.Client, method, path, bodyData, headers, ctx)
	} else if jsonData != nil { // 如果传入的数据是jsonData(map)，走下面分支
		body, err = json.Marshal(jsonData)
		if err != nil {
			logutil.LogrusObj.Error("json marshal err:", err)
			return nil, err
		}
		//req.Header.Set("Content-Type", "application/json")
		c.SetContentTypeWithHeaders("application/json")
		for k, v := range headers {
			c.Headers[k] = v
		}
		return c.do(req, resp, err, c.Client, method, path, bytes.NewBuffer(body), c.Headers, ctx)
	} else {
		for k, v := range headers {
			c.Headers[k] = v
		}
		return c.do(req, resp, err, c.Client, method, path, bytes.NewBuffer(body), c.Headers, ctx)
	}
}

func (c *HttpClient) do(req *http.Request, resp *http.Response, err error, client *http.Client, method string, path string, bodyData io.Reader,
	headers map[string]string, ctx context.Context) (*http.Response, error) {
	req, err = http.NewRequest(method, path, bodyData)
	// 传递上下文 如传递http请求超时时间
	req.WithContext(ctx)
	if err != nil {
		logutil.LogrusObj.Error("new request err:", err)
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err = client.Do(req)
	if err != nil {
		logutil.LogrusObj.Error("http request failed:", err)
		return nil, err
	}
	return resp, nil
}

// GetJsonData 数据解析函数
func (c *HttpClient) GetJsonData(response *http.Response) (map[string]interface{}, error) {
	logutil.LogrusObj.Info("response:", response)
	logutil.LogrusObj.Info("body:", response.Body)
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		logutil.LogrusObj.Info(err)
		return nil, err
	}
	logutil.LogrusObj.Info("[bodyBytes]:", string(bodyBytes))
	var jsonData map[string]interface{}
	err = json.Unmarshal(bodyBytes, &jsonData)
	if err != nil {
		logutil.LogrusObj.Info(err)
		return nil, err
	}
	return jsonData, nil
}

// RaiseForStatus 错误处理函数
func (c *HttpClient) RaiseForStatus(response *http.Response) error {
	if response.StatusCode >= 400 && response.StatusCode < 500 {
		httpError := fmt.Errorf("请求错误，请稍后重试，Response status:%s", string(response.StatusCode))
		logutil.LogrusObj.Error(httpError)
		return httpError
	} else if response.StatusCode >= 500 {
		httpError := fmt.Errorf("服务端异常，请稍后重试，Response status:%s", string(response.StatusCode))
		logutil.LogrusObj.Error(httpError)
		return httpError
	}
	return nil
}

// CheckResponseIsOk 判断响应是否正常
func (c *HttpClient) CheckResponseIsOk(response *http.Response) bool {
	status := c.RaiseForStatus(response)
	if status != nil {
		return false
	}
	return true
}

// ThrowErrorMsg 如果请求目标接口异常，调用此方法跑出错误信息
func (c *HttpClient) ThrowErrorMsg(response *http.Response) (string, error) {
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		logutil.LogrusObj.Info(err)
		return "", err
	}
	return string(bodyBytes), nil
}
