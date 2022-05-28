package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/che-kwas/iam-kit/logger"
)

// UsernameKey defines the key in gin context which represents the owner of the secret.
const UsernameKey = "username"

// Context returns a middleware that injects requestID and username to gin.Context.
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(logger.KeyRequestID, c.GetString(XRequestIDKey))
		c.Set(logger.KeyUsername, c.GetString(UsernameKey))
		c.Next()
	}
}
