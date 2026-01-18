package logging

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

type contextKey string

const loggerKey contextKey = "easygin-logger"

func (ck contextKey) String() string {
	return string(ck)
}

func init() {
	env := os.Getenv("APP_ENV")

	var logger *zap.Logger
	var err error

	if env == "production" {
		config := zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339Nano)
		config.DisableCaller = false
		config.DisableStacktrace = false
		logger, err = config.Build(zap.AddStacktrace(zapcore.ErrorLevel))
	} else {
		config := zap.NewDevelopmentConfig()
		config.DisableCaller = false
		config.DisableStacktrace = false
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339Nano)
		logger, err = config.Build(zap.AddStacktrace(zapcore.ErrorLevel))
	}

	if err != nil {
		panic(err)
	}
	globalLogger = logger
}

func GetGlobalLogger() *zap.Logger {
	return globalLogger
}

func GetLoggerFromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return logger
	}
	return globalLogger
}

func SetLoggerToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
