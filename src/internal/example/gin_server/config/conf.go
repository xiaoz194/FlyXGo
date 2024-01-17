package config

import (
	"github.com/spf13/viper"
	"github.com/xiaoz194/GoFlux/src/pkg/utils/logutil"
)

var (
	AppMode  string
	HttpPort string
)

func InitConfig() {
	logutil.LogrusObj.Info("start init...")
	viper.AddConfigPath("src/internal/example/gin_server/config")
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		logutil.LogrusObj.Printf("Error reading config file, falling back to environment variables: %v", err)
	}
	AppMode = viper.GetString("service.AppMode")
	HttpPort = viper.GetString("service.HttpPort")
	logutil.LogrusObj.Info("init success!")
}
