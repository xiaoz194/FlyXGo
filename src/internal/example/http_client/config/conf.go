package config

import (
	"github.com/spf13/viper"
	"github.com/xiaoz194/FlyXGo/src/pkg/utils/logutil"
	"time"
)

var (
	ExampleHttpRequestUrlPrefix string
	InsecureSkipVerify          bool
	MaxIdleConns                int
	IdleConnTimeout             time.Duration
	ResponseHeaderTimeout       time.Duration
)

func InitConfig() {
	viper.AddConfigPath("src/config")
	viper.SetConfigName("values")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		logutil.LogrusObj.Errorf("Error reading config file, falling back to environment variables: %v", err)
	}
	ExampleHttpRequestUrlPrefix = viper.GetString("httpClientServer.ExampleHttpRequestUrlPrefix")
	InsecureSkipVerify = viper.GetBool("httpClientServer.InsecureSkipVerify")
	MaxIdleConns = viper.GetInt("httpClientServer.MaxIdleConns")
	//IdleConnTimeout = viper.GetInt("httpClientServer.IdleConnTimeout")
	IdleConnTimeout = viper.GetDuration("httpClientServer.IdleConnTimeout")
	ResponseHeaderTimeout = viper.GetDuration("httpClientServer.ResponseHeaderTimeout")

}
