# FlyXGo

go的快速通信框架 fly旨在形容 快速开发

## 项目目录结构

```shell
├── logs
└── src
    ├── config
    ├── internal
    │   ├── core
    │   │   ├── grpc
    │   │   └── http
    │   └── example
    │       ├── gin_server
    │       │   ├── beans
    │       │   ├── cmd
    │       │   ├── config
    │       │   ├── controller
    │       │   ├── middleware
    │       │   ├── routes
    │       │   └── serializer
    │       └── http_client
    │           └── simple_demo
    └── pkg
        ├── e
        │   └── constant
        └── utils
            ├── dbutil
            └── logutil
```


## V1.x 版本日志
http客户端和grpc连接池客户端：

1）go语言基于的net/http封装的快速开发万能框架

2）grpc连接池封装与请求 优化连接池

3）提供样例example,给出基于gin框架的服务端，并提供客户端，给出使用案例

## Start
### 1 http客户端
#### 1.1 简介
core/http包提供万能接口以提供开发者快速开发,用户可以快速创建http客户端，可以调用统一的接口进行POST/GET/PUT 及 传递jsonData/io流 不同方式的请求。

FlyXHttpClient结构：
``` golang
type FlyXHttpClient struct {
	HttpClient *HttpClient
	Apis       map[string]ApiDef
}
```
包含两个子结构:
``` golang
type ApiDef struct {
	Method string
	Path   string
}

type HttpClient struct {
	Client  *http.Client
	Headers map[string]string
}
```

ApiDef 用于定义Api接口的method和path，例如：
``` golang
Apis := map[string]http_client.ApiDef{
		"test_get":  {Method: "GET", Path: fmt.Sprintf("%s/api/v1/test_get/uid/%d", config.ExampleHttpRequestUrlPrefix, uid)},
		"test_post": {Method: "POST", Path: fmt.Sprintf("%s/api/v1/test_post/uid/%d", config.ExampleHttpRequestUrlPrefix, uid)},
	}
```

HttpClient用于创建http客户端并提供一些核心方法。

FlyXHttpClient开放一个统一接口FlyXHttpClient.RetryApi供用户调用，该方法需要传递一个结构体参数RetryApiOptions
``` golang
// RetryApiOptions 创建一个结构体 来封装所有的参数
/**
参数说明：
	arg1 mode string 调用哪一个api  必传 NORMAL 通用模式，api接口  post_path  专门设置的post接口,用于直接传path的场景  postBody 传body数据流请求的post接口
	arg2 apisName string 请求的apis名称，除了 UNIQUE_PATH 模式外，均传此方法（目前是这样，不排除后续会有单独的"get"等模式） 因为这里需要通过通过ApiDef拿路径和请求方式  这里不传则置 ""
	arg3 path string 路径，如果是"post"模式，需要传path 因为它不通过ApiDef拿路径  这里不传则置 ""
	arg4 method string 请求方法，如果是"post"模式，需要传method 因为它不通过ApiDef拿路径  这里不传则置 ""
	arg5 jsonData map[string]interface{} 传递数据  这里不传则置 nil ,
	arg6 headers map[string]string 传递headers  这里不传则置 nil
	arg7 kwargs map[string]interface{} 额外的需要拼接到url的查询参数  这里不传则置 nil
	arg8 body *bytes.Buffer 有些地方不能传jsonData map[string]interface{} 只能通过body传递，这里不传则置 nil
	arg9 ctx 上下文，可以传递http请求超时时间等上下文信息，主要用来传递超时时间信息
	arg10 retryMax 失败重试最大次数 一般设置3次
	arg11 retryTimeout 失败后重试间歇时间 避免瞬间重试和减轻服务压力 一般设置5-10s 也可以结合其他策略，比如指数退避
**/

type RetryApiOptions struct {
	Mode         int
	ApisName     string
	Path         string
	Method       string
	JsonData     map[string]interface{}
	Headers      map[string]string
	Kwargs       map[string]interface{}
	Body         *bytes.Buffer
	Ctx          context.Context
	RetryMax     int
	RetryTimeout time.Duration
}
```
FlyXHttpClient.RetryApi接口提供三种模式调用，如下所示：
``` golang
const (
	NORMAL = iota
	UNIQUE_PATH
	POSTBODY
)
```
其中， NORMAL模式最常用，用户只需要定义ApiDefs，根据ApiDefs定义的name请求对应的接口
``` golang
resp, err = flyXHttpClient.api(options.ApisName, options.JsonData, options.Headers, options.Kwargs, options.Ctx)
```
UNIQUE_PATH 比较特殊，它不遵守本框架的Client定义方式，为了兼容没有继承FlyXHttpClient的客户端请求，该模式下，用户常规的传递请求Method和URL(Path)
``` golang
resp, err = flyXHttpClient.unique(options.Method, options.Path, options.JsonData, options.Headers, options.Kwargs, options.Ctx)
```
POSTBODY 用于传递字节流的特殊场景
``` golang
resp, err = flyXHttpClient.apiPassByBody(options.ApisName, options.Body, options.Headers, options.Ctx)
```
注意 上述接口均为内部调用，外部统一走FlyXHttpClient.RetryApi接口

具体逻辑可自行阅读代码，有疑问或者有改进的地方欢迎pr共同建设。

#### 1.2 使用

##### 步骤1 创建一个你的客户端 继承FlyXHttpClient
``` golang
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
```

##### 步骤2 实现每一个你定义的接口的具体逻辑
``` golang
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
```

你会发现，每一个方法实现都是相同的套路，非常简单，开发者只需要按模版写即可：
``` golang
func (exampleHttpClient *ExampleHttpClient) TestPost() (map[string]interface{}, error) {
    // 根据你需要传的参数来写
	// jsonData := xxx
	// headers := xxx
	// xxx... 
	// 传递上下文，设置http请求超时时间 需要改的是具体的超时时间，如果需要传递其他上下文信息，按需要改
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// 传递options参数 参数值根据自己的需求调整
	options := http_client.RetryApiOptions{
		Mode:         http_client.NORMAL,
		ApisName:     "xxx",
		Path:         "xxx",
		Method:       "xxx",
		JsonData:     jsonData,
		Headers:      headers,
		Kwargs:       nil,
		Body:         nil,
		Ctx:          ctx,
		RetryMax:     3,
		RetryTimeout: 5 * time.Second,
	}
	// 一下部分几乎不需要修改 模板
    // 调用RetryApi接口
	resp, err := exampleHttpClient.FlyXHttpClient.RetryApi(options)
	if err != nil {
		logutil.LogrusObj.Error("Api请求出错！错误信息：", err.Error())
		return nil, err
	}
	// 校验 CheckResponseIsOk 下面基本是一样的套路 不需要修改代码
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

```

这样写看似代码较多，实际上改动的地方很少，而且可以较好的对Client进行封装，当遇到比较复杂的场景，这样封装会比较好。当然你可以进一步的进行封装。

上面代码只个出了最基础的示例，一些更复杂的场景，如权限点校验等会在后续提供案例。

##### 步骤3 创建你的自定义Client 并调用请求
简单示例：
```shell
    var uid int64
	uid = 1
	exampleClient := simple_demo.NewExampleHttpClient(uid)
	// ==================== test get =========================
	resp, err := exampleClient.TestGet()
	if err != nil {
		logutil.LogrusObj.Fatal(err)
	}
	logutil.LogrusObj.Info(resp)
```

具体的案例，代码example目录提供了一个相对完整且简单的案例，包括一个基于gin框架的http服务端和使用core/http的http客户端