package middlewares

import (
	"easygin/internal/models"
	"easygin/pkg/apperror"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func ResponseWithErrorContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleError(c, err)
		}
	}
}

func handleError(c *gin.Context, err error) {
	requestID := GetRequestID(c)

	stackTrace := ""
	if stackTracer, ok := err.(interface{ StackTrace() errors.StackTrace }); ok {
		stackTrace = fmt.Sprintf("%+v", stackTracer.StackTrace())
	}

	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		details := map[string]interface{}{
			"error": appErr.Unwrap().Error(),
			"stack": stackTrace,
		}
		response := models.ErrorResponseWithDetails(appErr.Code(), appErr.Message(), requestID, details)
		c.JSON(appErr.HTTPStatusCode(), response)
		return
	}

	details := map[string]interface{}{
		"error": err.Error(),
		"stack": stackTrace,
	}
	response := models.ErrorResponseWithDetails("UNEXPECTED_ERROR", "An unexpected error occurred", requestID, details)
	c.JSON(http.StatusInternalServerError, response)
}
