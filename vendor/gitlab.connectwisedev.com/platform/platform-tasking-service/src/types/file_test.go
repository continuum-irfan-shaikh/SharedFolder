package types

import (
	"fmt"
	"testing"
)

func TestMultiError_Error(t *testing.T) {
	err := NewMultiError()
	err2 := NewMultiError()
	err.ToError()
	err = append(err, fmt.Errorf("asd"))
	err = append(err, &err2)
	err.ToError()
	err.Error()
}
