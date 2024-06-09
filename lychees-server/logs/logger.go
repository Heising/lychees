package logs

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"lychees-server/configs"
	"os"

	"github.com/charmbracelet/log"
)

var Logger *log.Logger

func init() {
	if configs.Config.Debug {
		Logger = log.NewWithOptions(os.Stderr, log.Options{
			Level:           log.DebugLevel,
			ReportCaller:    true, // 调用的文件名和函数
			ReportTimestamp: true, // 时间
		})
	} else {

		jackLogger := &lumberjack.Logger{
			Filename:   "./logs/log",
			MaxSize:    500,  // megabytes 最大500m
			MaxBackups: 3,    // 超出MaxSize后可以3个备份
			MaxAge:     28,   //days
			Compress:   true, // disabled by default 压缩备份
		}

		multiWriter := io.MultiWriter(os.Stderr, jackLogger)

		Logger = log.NewWithOptions(multiWriter, log.Options{
			Level: log.WarnLevel,
			ReportTimestamp: true, // 时间
		})
	}
}
