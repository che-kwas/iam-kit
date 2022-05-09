package db

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInstanceId(t *testing.T) {
	assert := assert.New(t)

	id1 := GetInstanceID(1, "secret-")
	id2 := GetInstanceID(2, "secret-")
	assert.True(strings.HasPrefix(id1, "secret-"))
	assert.True(strings.HasPrefix(id2, "secret-"))
	assert.True(len(id1) > 12)
	assert.True(len(id2) > 12)
	assert.NotEqual(id1, id2)
}
