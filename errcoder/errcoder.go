package errcoder

import (
	"fmt"
	"net/http"

	"github.com/marmotedu/errors"
	"github.com/novalagung/gubrak"
)

var ValidHTTPStatus = []int{200, 400, 401, 403, 404, 500}

// ErrCoder implements `github.com/marmotedu/errors` Coder interface.
type ErrCoder struct {
	code       int
	httpStatus int
	message    string
}

var _ errors.Coder = &ErrCoder{}

// Code returns the integer error code.
func (coder ErrCoder) Code() int {
	return coder.code
}

// String implements stringer.
func (coder ErrCoder) String() string {
	return coder.message
}

// Reference returns the reference document.
func (coder ErrCoder) Reference() string {
	return ""
}

// HTTPStatus returns the associated HTTP status code.
func (coder ErrCoder) HTTPStatus() int {
	if coder.httpStatus == 0 {
		return http.StatusInternalServerError
	}

	return coder.httpStatus
}

func Register(code int, httpStatus int, message string) {
	found, _ := gubrak.Includes(ValidHTTPStatus, httpStatus)
	if !found {
		panic(fmt.Sprintf("http code not in `%v`", ValidHTTPStatus))
	}

	coder := &ErrCoder{
		code:       code,
		httpStatus: httpStatus,
		message:    message,
	}

	errors.MustRegister(coder)
}
