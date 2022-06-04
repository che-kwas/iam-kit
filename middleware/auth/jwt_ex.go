package auth

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/marmotedu/errors"

	"github.com/che-kwas/iam-kit/code"
	"github.com/che-kwas/iam-kit/httputil"
	"github.com/che-kwas/iam-kit/middleware"
)

// Defined errors.
var (
	ErrMissingKID    = errors.New("missing kid in token header")
	ErrMissingSecret = errors.New("missing secret in cache")
)

// Secret contains the basic information of the secret key.
type Secret struct {
	Username string
	ID       string
	Key      string
	Expires  int64
}

// GetSecretFunc is a type of func for getting the Secret.
type GetSecretFunc func(kid string) (Secret, error)

// JWTExStrategy defines jwt bearer authentication strategy with user-specific secret.
type JWTExStrategy struct {
	get GetSecretFunc
}

var _ middleware.AuthStrategy = &JWTExStrategy{}

// NewJWTExStrategy creates a jwt strategy.
func NewJWTExStrategy(get GetSecretFunc) JWTExStrategy {
	return JWTExStrategy{get}
}

// AuthFunc defines jwt strategy as the gin authentication middleware.
func (j JWTExStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		if len(header) == 0 {
			httputil.WriteResponse(c, errors.WithCode(code.ErrHeaderInvalid, "Authorization header cannot be empty."), nil)
			c.Abort()

			return
		}

		var rawJWT string
		fmt.Sscanf(header, "Bearer %s", &rawJWT)

		var secret Secret
		claims := &jwt.MapClaims{}
		// Verify the token
		parsedT, err := jwt.ParseWithClaims(rawJWT, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is HMAC signature
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, ErrMissingKID
			}

			var err error
			secret, err = j.get(kid)
			if err != nil {
				return nil, ErrMissingSecret
			}

			return []byte(secret.Key), nil
		})
		if err != nil || !parsedT.Valid {
			httputil.WriteResponse(c, errors.WithCode(code.ErrTokenInvalid, err.Error()), nil)
			c.Abort()

			return
		}

		if KeyExpired(secret.Expires) {
			tm := time.Unix(secret.Expires, 0).Format("2006-01-02 15:04:05")
			httputil.WriteResponse(c, errors.WithCode(code.ErrTokenExpired, "expired at: %s", tm), nil)
			c.Abort()

			return
		}

		c.Set(middleware.UsernameKey, secret.Username)
		c.Next()
	}
}

// KeyExpired checks if a key has expired.
func KeyExpired(expires int64) bool {
	if expires >= 1 {
		return time.Now().After(time.Unix(expires, 0))
	}

	return false
}
