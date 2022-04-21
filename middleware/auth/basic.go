package auth

import (
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"

	"github.com/che-kwas/iam-kit/errcode"
	"github.com/che-kwas/iam-kit/middleware"
	"github.com/che-kwas/iam-kit/util"
)

// BasicStrategy defines Basic authentication strategy.
type BasicStrategy struct {
	compare func(username string, password string) bool
}

var _ middleware.AuthStrategy = &BasicStrategy{}

// NewBasicStrategy creates basic strategy with compare function.
func NewBasicStrategy(compare func(username string, password string) bool) BasicStrategy {
	return BasicStrategy{
		compare: compare,
	}
}

// AuthFunc implements the AuthStrategy interface.
func (b BasicStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			util.WriteResponse(
				c,
				errors.WithCode(errcode.ErrHeaderInvalid, "Authorization header is invalid."),
				nil,
			)
			c.Abort()

			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !b.compare(pair[0], pair[1]) {
			util.WriteResponse(
				c,
				errors.WithCode(errcode.ErrHeaderInvalid, "Authorization header is invalid."),
				nil,
			)
			c.Abort()

			return
		}

		c.Set(middleware.UsernameKey, pair[0])

		c.Next()
	}
}
