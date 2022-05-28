package logger

import "sync"

var (
	mu sync.RWMutex
	gl = NewLogger()
)

// L returns the global Logger.
func L() *Logger {
	mu.RLock()
	l := gl
	mu.RUnlock()
	return l
}
