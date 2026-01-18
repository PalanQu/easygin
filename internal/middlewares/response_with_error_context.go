package middlewares

import (
	"easygin/internal/models"
	"easygin/pkg/apperror"
	"easygin/pkg/logging"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func ResponseWithErrorContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logging.GetLoggerFromContext(c.Request.Context())

		logger.Info("Request received",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
		)

		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleError(c, err)
		}
	}
}

func handleError(c *gin.Context, err error) {
	requestID := GetRequestID(c)
	logger := logging.GetLoggerFromContext(c.Request.Context())

	stackTrace := ""
	if stackTracer, ok := err.(interface{ StackTrace() errors.StackTrace }); ok {
		stackTrace = fmt.Sprintf("%+v", stackTracer.StackTrace())
	}

	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		logger.Info("Response error",
			zap.Int("status", appErr.HTTPStatusCode()),
			zap.String("error_code", string(appErr.Code())),
			zap.String("error_message", appErr.Message()),
		)

		details := map[string]interface{}{
			"error": appErr.Unwrap().Error(),
			"stack": stackTrace,
		}
		response := models.ErrorResponseWithDetails(appErr.Code(), appErr.Message(), requestID, details)
		c.JSON(appErr.HTTPStatusCode(), response)
		return
	}

	logger.Info("Response error",
		zap.Int("status", http.StatusInternalServerError),
		zap.String("error_code", "UNEXPECTED_ERROR"),
		zap.String("error_message", "An unexpected error occurred"),
	)

	details := map[string]interface{}{
		"error": err.Error(),
		"stack": stackTrace,
	}
	response := models.ErrorResponseWithDetails("UNEXPECTED_ERROR", "An unexpected error occurred", requestID, details)
	c.JSON(http.StatusInternalServerError, response)
}
