package logutil

import (
	"fmt"
	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

/*
日志工具
*/

var LogrusObj *logrus.Logger

type MyFormatter struct{}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	msg := fmt.Sprintf("%s\n", entry.Message)
	_, file, line, _ := runtime.Caller(7)
	filename := filepath.Base(file)
	lineStr := strconv.Itoa(line)
	return []byte(fmt.Sprintf("-[%s]-[%s]-%s:%s-%s", entry.Time.Format("2006-01-02 15:04:05"), entry.Level, filename, lineStr, msg)), nil
}

func init() {
	if LogrusObj != nil {
		src, _ := setOutputFile()
		writers := []io.Writer{
			src,
			os.Stdout,
		}
		fileAndStdoutWriter := io.MultiWriter(writers...)
		//设置输出
		//LogrusObj.Out = src
		LogrusObj.SetOutput(fileAndStdoutWriter)
		return
	}
	//实例化
	logger := logrus.New()
	src, _ := setOutputFile()
	writers := []io.Writer{
		src,
		os.Stdout,
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	//设置输出
	//logger.Out = src
	logger.SetOutput(fileAndStdoutWriter)
	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	//打开日志输出文件和行号功能
	logger.SetReportCaller(true)
	//设置日志格式
	logger.SetFormatter(&MyFormatter{})
	logger.Info("now initialize logger client logrus...")
	//集成sentry-hook
	sentryDsn := os.Getenv("SENTRY_DSN")
	if sentryDsn != "" {
		logger.Info("检测到sentryDsn环境变量,捕获到的异常日志将上报数据到sentry...")
		hook, err := logrus_sentry.NewSentryHook(sentryDsn, []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel})
		if err != nil {
			logger.Error(err.Error())
		}
		hook.Timeout = 20 * time.Second
		hook.StacktraceConfiguration.Enable = true
		logger.Hooks.Add(hook)
	}
	LogrusObj = logger
	logger.Info("initialize logger client logrus successful!")
}

func setOutputFile() (*os.File, error) {
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(logFilePath, 0777); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	logFileName := now.Format("2006-01-02") + ".log"
	//日志文件
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	//写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return src, nil
}
