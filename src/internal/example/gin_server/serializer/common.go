package serializer

// Response 基础序列化器
type Response struct {
	Code   int         `json:"code"` //default 0
	Status int         `json:"status_code"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
}
