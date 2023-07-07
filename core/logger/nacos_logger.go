package logger

import (
	"github.com/nacos-group/nacos-sdk-go/common/logger"
)

type NacosLogger struct {
}

func (l NacosLogger) Info(args ...interface{}) {
	Info(args)
}
func (l NacosLogger) Warn(args ...interface{}) {
	Warn(args)
}
func (l NacosLogger) Error(args ...interface{}) {
	Error(args)
}
func (l NacosLogger) Debug(args ...interface{}) {
	Debug(args)
}

func (l NacosLogger) Infof(fmt string, args ...interface{}) {
	InfoF(fmt, args...)
}
func (l NacosLogger) Warnf(fmt string, args ...interface{}) {
	WarnF(fmt, args...)
}
func (l NacosLogger) Errorf(fmt string, args ...interface{}) {
	ErrorF(fmt, args...)
}
func (l NacosLogger) Debugf(fmt string, args ...interface{}) {
	DebugF(fmt, args...)
}

func InitNacosLog() {
	logger.SetLogger(NacosLogger{})
}
