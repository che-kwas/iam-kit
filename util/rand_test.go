package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RandString(t *testing.T) {
	assert := assert.New(t)

	r1 := RandString(10)
	assert.Equal(10, len(r1))
	r2 := RandString(10)
	assert.Equal(10, len(r2))
	assert.NotEqual(r1, r2)

	r3 := RandString(64)
	assert.Equal(64, len(r3))
}
