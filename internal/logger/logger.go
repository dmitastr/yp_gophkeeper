package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ILogger interface {
	Info(string, ...zap.Field)
	Error(string, ...zap.Field)
	Fatal(string, ...zap.Field)
}

type Logger struct {
	logger *zap.Logger
}

func NewLogger() ILogger {
	config := zap.NewDevelopmentConfig()
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	zapcore.TimeEncoderOfLayout("Jan _2 15:04:05.000000000")
	encoderConfig.StacktraceKey = "" // to hide stacktrace info
	config.EncoderConfig = encoderConfig

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	return Logger{logger: logger}

}
func (l Logger) Info(message string, fields ...zap.Field) {
	l.logger.Info(message, fields...)
}

func (l Logger) Error(message string, fields ...zap.Field) {
	l.logger.Error(message, fields...)
}

func (l Logger) Fatal(message string, fields ...zap.Field) {
	l.logger.Fatal(message, fields...)
}
