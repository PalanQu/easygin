package middlewares

import (
	"easygin/pkg/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logging.GetGlobalLogger()
		loggerContextData := map[string]string{}

		// add your custom fields here
		loggerContextData["foo"] = "bar"

		fields := make([]zap.Field, 0, len(loggerContextData))
		for k, v := range loggerContextData {
			fields = append(fields, zap.String(k, v))
		}
		contextLogger := logger.With(fields...)
		requestCtx := c.Request.Context()
		ctxWithLogger := logging.SetLoggerToContext(requestCtx, contextLogger)
		c.Request = c.Request.WithContext(ctxWithLogger)
		c.Next()
	}
}
