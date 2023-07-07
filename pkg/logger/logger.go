package logger

import (
	"github.com/micro/go-micro/v2/logger"
	mzap "github.com/micro/go-plugins/logger/zap/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

type Config struct {
	Level      string   `json:"level"`
	Outputs    []string `json:"outputs"`
	ErrOutputs []string `json:"errOutputs"`
	MaxSize    int      `json:"maxSize"`
	MaxNum     int      `json:"maxNum"`
	Console    bool     `json:"console"`
}

func NewLogger(config Config) (logger.Logger, error) {
	level := getLoggerLevel(config.Level)

	newLog, err := mzap.NewLogger(
		mzap.WithCallerSkip(2),
		mzap.WithConfig(zap.Config{
			Level:       zap.NewAtomicLevelAt(level),
			Development: true,
			Encoding:    "json",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "t",
				LevelKey:       "level",
				NameKey:        "log",
				CallerKey:      "caller",
				MessageKey:     "msg",
				StacktraceKey:  "trace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.FullCallerEncoder,
			},
			OutputPaths:      config.Outputs,
			ErrorOutputPaths: config.ErrOutputs,
			//InitialFields: map[string]interface{}{
			//	"app": "test",
			//},
		}),
	)
	if err != nil {
		return nil, err
	}
	return newLog, nil
}
