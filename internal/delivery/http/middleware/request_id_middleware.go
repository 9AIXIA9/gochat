package middleware

import (
	"gochat/internal/shared/kernel"
	"gochat/pkg/ctxutil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

// NewRequestIDMiddleware ensures every request carries a request_id for correlation.
func NewRequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(requestIDHeader)
		if rid == "" {
			rid = uuid.Must(uuid.NewV7()).String()
		}

		requestID := kernel.OperationID(rid)
		c.Request = c.Request.WithContext(ctxutil.WithRequestID(c.Request.Context(), requestID))
		c.Writer.Header().Set(requestIDHeader, rid)

		c.Next()
	}
}
