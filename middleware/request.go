package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateRequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.New()
		str := strings.ReplaceAll(id.String(), "-", "")
		requestId := strings.ToLower(str)
		c.Writer.Header().Set("X-Request-Id", requestId)
		c.Set("RequestId", requestId)
		c.Next()
	}
}
