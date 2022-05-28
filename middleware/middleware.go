// Package middleware defines multiple gin middlewares.
package middleware // import "github.com/che-kwas/iam-kit/middleware"

import (
	"github.com/gin-gonic/gin"
)

// Middlewares store default middlewares.
var Middlewares = map[string]gin.HandlerFunc{
	"recovery":  gin.Recovery(),
	"secure":    Secure(),
	"options":   Options(),
	"nocache":   NoCache(),
	"cors":      Cors(),
	"requestid": RequestID(),
	"context":   Context(),
	"logger":    Logger(),
}
