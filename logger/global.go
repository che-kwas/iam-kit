package logger

import "sync"

var (
	gl   *Logger
	once sync.Once
)

// L returns the global logger.
func L() *Logger {
	if gl != nil {
		return gl
	}

	once.Do(func() {
		gl = NewLogger()
	})
	return gl
}
