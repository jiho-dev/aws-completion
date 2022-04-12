package log

import (
	"log/syslog"

	"github.com/jiho-dev/aws-completion/config"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

var logger *logrus.Logger

func InitLogger(level logrus.Level) {
	logger = &logrus.Logger{
		//Out:   os.Stderr,
		Level: level,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg%",
		},
	}

	syslogger, err := syslog.New(syslog.LOG_INFO, config.EXEC_NAME)
	if err != nil {
		logger.Fatalln(err)
	}

	logger.SetOutput(syslogger)
}

func GetLogger() *logrus.Logger {
	return logger
}
