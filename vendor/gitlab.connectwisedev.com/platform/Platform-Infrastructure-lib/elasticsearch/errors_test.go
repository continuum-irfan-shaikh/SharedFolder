package elasticsearch

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewErrNotFound(t *testing.T) {
	err := NewErrNotFound("error")
	assert.NotNil(t, err)

	msg := err.Error()
	assert.NotNil(t, msg)

	ok := IsErrNotFound(errors.New("error"))
	assert.False(t, ok)
}
