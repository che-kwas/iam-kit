package middleware

import (
	"github.com/gin-gonic/gin"
)

// AuthStrategy is used to do authentication.
type AuthStrategy interface {
	AuthFunc() gin.HandlerFunc
}

// AuthOperator is used to switch between different auth strategy.
type AuthOperator struct {
	strategy AuthStrategy
}

// SetStrategy sets the strategy to operator.
func (operator *AuthOperator) SetStrategy(strategy AuthStrategy) {
	operator.strategy = strategy
}

// AuthFunc executes authentication.
func (operator *AuthOperator) AuthFunc() gin.HandlerFunc {
	return operator.strategy.AuthFunc()
}
