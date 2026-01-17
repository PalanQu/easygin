package logging

import (
	"context"
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
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339Nano)
	l, err := config.Build()
	if err != nil {
		panic(err)
	}
	globalLogger = l
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
