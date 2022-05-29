package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/che-kwas/iam-kit/logger"
)

// Logger returns a middleware that adds requestID and username to the loggin context.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

		param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path

		logger.L().X(c).Info(formatLogParams(param))
	}
}

func formatLogParams(param gin.LogFormatterParams) string {

	return fmt.Sprintf("%3d - [%s] \"%v %s %s\" %s",
		param.StatusCode,
		param.ClientIP,
		param.Latency,
		param.Method,
		param.Path,
		param.ErrorMessage,
	)
}
