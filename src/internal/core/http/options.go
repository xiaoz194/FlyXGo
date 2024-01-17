package http_client

import (
	"bytes"
	"context"
	"time"
)

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
