// Code generated by "codegen ./code"; DO NOT EDIT.

package code

import "github.com/che-kwas/iam-kit/errcoder"

// init register error codes defines in `github.com/marmotedu/errors`
func init() {
	errcoder.Register(ErrSuccess, 200, "OK")
	errcoder.Register(ErrUnknown, 500, "Internal server error")
	errcoder.Register(ErrBadParams, 400, "Bad request parameters")
	errcoder.Register(ErrNotFound, 404, "Not found")
	errcoder.Register(ErrPasswordInvalid, 401, "Password invalid")
	errcoder.Register(ErrHeaderInvalid, 401, "Authorization header invalid")
	errcoder.Register(ErrSignatureInvalid, 401, "Signature invalid")
	errcoder.Register(ErrTokenInvalid, 401, "Token invalid")
	errcoder.Register(ErrTokenExpired, 401, "Token expired")
	errcoder.Register(ErrPermissionDenied, 403, "Permission denied")
	errcoder.Register(ErrDatabase, 500, "Database error")
}
