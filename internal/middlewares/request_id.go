package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	RequestIDKey        = "X-Request-ID"
	RequestIDContextKey = "request_id"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDKey)

		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set(RequestIDContextKey, requestID)

		c.Header(RequestIDKey, requestID)

		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDContextKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}
