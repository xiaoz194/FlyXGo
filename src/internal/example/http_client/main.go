package main

import (
	"github.com/xiaoz194/FlyXGo/src/internal/example/http_client/simple_demo"
	"github.com/xiaoz194/FlyXGo/src/pkg/utils/logutil"
)

func main() {
	var uid int64
	uid = 1
	exampleClient := simple_demo.NewExampleHttpClient(uid)
	// ==================== test get =========================
	resp, err := exampleClient.TestGet()
	if err != nil {
		logutil.LogrusObj.Fatal(err)
	}
	logutil.LogrusObj.Info(resp)

	// ===================== test post =======================
	resp2, err2 := exampleClient.TestPost()
	if err2 != nil {
		logutil.LogrusObj.Fatal(err2)
	}
	logutil.LogrusObj.Info(resp2["data"])

}
