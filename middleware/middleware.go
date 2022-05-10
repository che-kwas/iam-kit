package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// XRequestIDKey defines X-Request-ID key string.
const XRequestIDKey = "X-Request-ID"

const corsMaxAge = 12

// Middlewares store default middlewares.
var Middlewares = map[string]gin.HandlerFunc{
	"recovery":  gin.Recovery(),
	"secure":    Secure,
	"options":   Options,
	"nocache":   NoCache,
	"cors":      Cors(),
	"requestid": RequestID(),
}

// Secure is a middleware function that appends security and resource access headers.
func Secure(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-XSS-Protection", "1; mode=block")

	if c.Request.TLS != nil {
		c.Header("Strict-Transport-Security", "max-age=31536000")
	}
}

// Options is a middleware function that appends headers for OPTIONS requests.
func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
	}
}

// NoCache is a middleware function that appends headers
// to prevent the client from caching the HTTP response.
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Next()
}

// Cors returns a middleware that adds cors headers.
func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           corsMaxAge * time.Hour,
	})
}

// RequestID returns a middleware that injects a 'X-Request-ID' into the context and request/response header of each request.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		rid := c.GetHeader(XRequestIDKey)

		if rid == "" {
			rid = uuid.Must(uuid.NewV4()).String()
			c.Request.Header.Set(XRequestIDKey, rid)
			c.Set(XRequestIDKey, rid)
		}

		// Set XRequestIDKey header
		c.Writer.Header().Set(XRequestIDKey, rid)
		c.Next()
	}
}
