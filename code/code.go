// Package code defines the base errors with corresponding error code.
package code // import "github.com/che-kwas/iam-kit/code"

//go:generate codegen
//go:generate codegen -doc -output ../../../errcode_base.md

// Common: basic errors (1000xx).
const (
	// ErrSuccess - 200: OK.
	ErrSuccess int = iota + 100001

	// ErrUnknown - 500: Internal server error.
	ErrUnknown

	// ErrBadParams - 400: Bad request parameters.
	ErrBadParams

	// ErrPageNotFound - 404: Not found.
	ErrNotFound
)

// common: authorization / authentication errors (1001xx).
const (
	// ErrPasswordInvalid - 401: Password invalid.
	ErrPasswordInvalid int = iota + 100101

	// ErrHeaderInvalid - 401: Authorization header invalid.
	ErrHeaderInvalid

	// ErrSignatureInvalid - 401: Signature invalid.
	ErrSignatureInvalid

	// ErrTokenInvalid - 401: Token invalid.
	ErrTokenInvalid

	// ErrTokenExpired - 401: Token expired.
	ErrTokenExpired

	// PermissionDenied - 403: Permission denied.
	ErrPermissionDenied
)

// common: database errors (1002xx).
const (
	// ErrDatabase - 500: Database error.
	ErrDatabase int = iota + 100201
)
