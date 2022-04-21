package shutdown

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	callCnt    int32 = 0
	callbackOk       = func(ctx context.Context) error {
		atomic.AddInt32(&callCnt, 1)
		return nil
	}
	callbackErr = func(ctx context.Context) error {
		return errors.New("shutdown error")
	}
)

func TestAddShutdownCallback(t *testing.T) {
	assert := assert.New(t)

	gs := New(5 * time.Second)
	assert.Equal(0, len(gs.callbacks))
	gs.AddShutdownCallback(ShutdownFunc(callbackOk))
	gs.AddShutdownCallback(ShutdownFunc(callbackErr))
	assert.Equal(2, len(gs.callbacks))
}

func TestCallbacksGetCalled(t *testing.T) {
	assert := assert.New(t)
	gs := New(5 * time.Second)
	gs.AddShutdownCallback(ShutdownFunc(callbackOk))
	gs.AddShutdownCallback(ShutdownFunc(callbackOk))
	callCnt = 0

	gs.waitCallbacks()
	assert.Equal(int32(2), callCnt)
}
