package consts

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()

	// 设置日志格式为JSON格式（更结构化）
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 设置日志级别
	Logger.SetLevel(logrus.InfoLevel)

	// 输出到标准输出
	Logger.SetOutput(os.Stdout)
}
