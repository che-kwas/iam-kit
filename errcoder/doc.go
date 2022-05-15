// Package errcoder implements `github.com/marmotedu/errors` Coder interface.
//
// errcoder only allowed the following http status:
// StatusOK                           = 200 // RFC 7231, 6.3.1
// StatusBadRequest                   = 400 // RFC 7231, 6.5.1
// StatusUnauthorized                 = 401 // RFC 7235, 3.1
// StatusForbidden                    = 403 // RFC 7231, 6.5.3
// StatusNotFound                     = 404 // RFC 7231, 6.5.4
// StatusInternalServerError          = 500 // RFC 7231, 6.6.1
package errcoder // import "github.com/che-kwas/iam-kit/errcoder"
