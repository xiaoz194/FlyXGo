package main

import (
	"context"
	"github.com/xiaoz194/FlyXGo/src/internal/example/gin_server/config"
	"github.com/xiaoz194/FlyXGo/src/internal/example/gin_server/routes"
	"github.com/xiaoz194/FlyXGo/src/pkg/utils/dbutil"
	"github.com/xiaoz194/FlyXGo/src/pkg/utils/logutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func echoBanner() {
	logutil.LogrusObj.Info(`
	      ____  ___  _____ _    _   ___  __      ____ ___ _   _ 
  		/ ___|/ _ \|  ___| |  | | | \ \/ /     / ___|_ _| \ | |
 		| |  _| | | | |_  | |  | | | |\  /_____| |  _ | ||  \| |
 		| |_| | |_| |  _| | |__| |_| |/  \_____| |_| || || |\  |
  		\____|\___/|_|   |_____\___//_/\_\     \____|___|_| \_|

	`)
}

func init() {
	echoBanner()
	config.InitConfig()
}

func startServer() {
	// 初始化gin引擎
	r := routes.NewRouter()

	server := &http.Server{
		Addr:    ":" + config.HttpPort,
		Handler: r,
	}

	// 启动gin服务
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logutil.LogrusObj.Fatalf("Server failed to run: %v\n", err)
		}
	}()

	// 设置信号捕获
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logutil.LogrusObj.Info("shutting down server...")
	// 创建上下文，关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅的关闭gin服务
	if err := server.Shutdown(ctx); err != nil {
		logutil.LogrusObj.Fatalf("server forced to shutdown: %v", err)
	}

	// 关闭mysql连接池
	defer func() {
		sqlDB, _ := dbutil.DB.DB()
		sqlDB.Close()
	}()

}

func main() {
	// 默认server启动
	startServer()

}
