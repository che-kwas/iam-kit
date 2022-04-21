package util

import (
	"net/http"

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
		coder := errors.ParseCoder(err)
		c.JSON(coder.HTTPStatus(), ErrResponse{
			Code:    coder.Code(),
			Message: coder.String(),
		})

		return
	}

	c.JSON(http.StatusOK, data)
}
