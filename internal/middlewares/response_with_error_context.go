package middlewares

import (
	"easygin/pkg/apperror"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type ErrorResponse struct {
	Code    apperror.AppErrorCode `json:"code"`
	Message string                `json:"message"`
	Details interface{}           `json:"details,omitempty"`
}

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
	stackTrace := ""
	if stackTracer, ok := err.(interface{ StackTrace() errors.StackTrace }); ok {
		stackTrace = fmt.Sprintf("%+v", stackTracer.StackTrace())
	}

	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		response := ErrorResponse{
			Code:    appErr.Code(),
			Message: appErr.Message(),
		}

		response.Details = map[string]interface{}{
			"error": appErr.Unwrap().Error(),
			"stack": stackTrace,
		}

		c.JSON(appErr.HTTPStatusCode(), response)
		return
	}

	response := ErrorResponse{
		Code:    "UNEXPECTED_ERROR",
		Message: "An unexpected error occurred",
	}

	response.Details = map[string]interface{}{
		"error": err.Error(),
		"stack": stackTrace,
	}

	c.JSON(http.StatusInternalServerError, response)
}
