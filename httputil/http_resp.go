// Package httputil provide useful functions to handle http response.
package httputil // import "github.com/che-kwas/iam-kit/httputil"

import (
	"net/http"

	"github.com/che-kwas/iam-kit/logger"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
)

// ErrResponse defines the return messages when an error occurred.
type ErrResponse struct {
	// Code defines the business error code.
	Code int `json:"code"`
	// Message defines the user-safe error message.
	Message string `json:"message"`
}

// WriteResponse writes the response data or an error into HTTP reponse.
func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		logger.L().X(c).Errorf("response error: %v", err)
		coder := errors.ParseCoder(err)

		c.JSON(coder.HTTPStatus(), ErrResponse{
			Code:    coder.Code(),
			Message: coder.String(),
		})

		return
	}

	logger.L().X(c).Infof("response data: %v", data)
	c.JSON(http.StatusOK, data)
}
