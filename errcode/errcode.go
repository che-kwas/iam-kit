package errcode

import (
	"fmt"
	"net/http"

	"github.com/marmotedu/errors"
	"github.com/novalagung/gubrak"
)

var ValidHTTPStatus = []int{200, 400, 401, 403, 404, 500}

// ErrCode implements `github.com/marmotedu/errors` Coder interface.
type ErrCode struct {
	code       int
	httpStatus int
	message    string
}

var _ errors.Coder = &ErrCode{}

// Code returns the integer error code.
func (coder ErrCode) Code() int {
	return coder.code
}

// String implements stringer.
func (coder ErrCode) String() string {
	return coder.message
}

// Reference returns the reference document.
func (coder ErrCode) Reference() string {
	return ""
}

// HTTPStatus returns the associated HTTP status code.
func (coder ErrCode) HTTPStatus() int {
	if coder.httpStatus == 0 {
		return http.StatusInternalServerError
	}

	return coder.httpStatus
}

func register(code int, httpStatus int, message string) {
	found, _ := gubrak.Includes(ValidHTTPStatus, httpStatus)
	if !found {
		panic(fmt.Sprintf("http code not in `%v`", ValidHTTPStatus))
	}

	coder := &ErrCode{
		code:       code,
		httpStatus: httpStatus,
		message:    message,
	}

	errors.MustRegister(coder)
}
