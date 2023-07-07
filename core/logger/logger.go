package logger

import (
	micrologger "github.com/micro/go-micro/v2/logger"
	conf2 "github.com/star-table/usercenter/core/conf"
	"github.com/star-table/usercenter/pkg/logger"
)

var (
	log micrologger.Logger = nil
)

func init() {
	logConf := conf2.Cfg.Logger
	newLog, err := logger.NewLogger(logger.Config{
		Level:      logConf.Level,
		Outputs:    logConf.Outputs,
		ErrOutputs: logConf.ErrOutputs,
		MaxSize:    logConf.MaxSize,
		MaxNum:     logConf.MaxNum,
		Console:    logConf.Console,
	})
	if err != nil {
		panic(err)
	}
	log = newLog
}

func InfoF(format string, v ...interface{}) {
	log.Logf(micrologger.InfoLevel, format, v...)
}

func ErrorF(format string, v ...interface{}) {
	log.Logf(micrologger.ErrorLevel, format, v...)
}

func WarnF(format string, v ...interface{}) {
	log.Logf(micrologger.WarnLevel, format, v...)
}

func DebugF(format string, v ...interface{}) {
	log.Logf(micrologger.DebugLevel, format, v...)
}

func TraceF(format string, v ...interface{}) {
	log.Logf(micrologger.TraceLevel, format, v...)
}

func FatalF(format string, v ...interface{}) {
	log.Logf(micrologger.FatalLevel, format, v...)
}

func Info(format interface{}) {
	log.Log(micrologger.InfoLevel, format)
}

func Error(format interface{}) {
	log.Log(micrologger.ErrorLevel, format)
}

func Warn(format interface{}) {
	log.Log(micrologger.WarnLevel, format)
}

func Debug(format interface{}) {
	log.Log(micrologger.DebugLevel, format)
}

func Trace(format interface{}) {
	log.Log(micrologger.TraceLevel, format)
}

func Fatal(format interface{}) {
	log.Log(micrologger.FatalLevel, format)
}
