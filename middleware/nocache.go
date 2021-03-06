package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// NoCache returns a middleware that appends headers
// to prevent the client from caching the HTTP response.
func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		c.Next()
	}
}
