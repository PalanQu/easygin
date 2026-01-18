package middlewares

import (
	"easygin/pkg/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logging.GetGlobalLogger()

		// Get request ID from context
		requestID := GetRequestID(c)

		// Add request ID to logger fields
		contextLogger := logger.With(
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)

		requestCtx := c.Request.Context()
		ctxWithLogger := logging.SetLoggerToContext(requestCtx, contextLogger)
		c.Request = c.Request.WithContext(ctxWithLogger)
		c.Next()
	}
}
