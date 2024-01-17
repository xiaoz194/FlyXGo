package config

import (
	"github.com/spf13/viper"
	"github.com/xiaoz194/FlyXGo/src/pkg/utils/logutil"
)

var (
	ExampleHttpRequestUrlPrefix string
)

func InitConfig() {
	viper.AddConfigPath("src/config")
	viper.SetConfigName("values")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		logutil.LogrusObj.Errorf("Error reading config file, falling back to environment variables: %v", err)
	}
	ExampleHttpRequestUrlPrefix = viper.GetString("httpClientServer.ExampleHttpRequestUrlPrefix")

}
