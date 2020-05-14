package elasticsearch

import (
	"fmt"

	"github.com/pkg/errors"
)

type ErrNotFound struct {
	msg string
}

// Error so that ErrNotFound implements error interface
func (e *ErrNotFound) Error() string {
	return e.msg
}

// NewErrNotFound is constructor for ErrNotFound
func NewErrNotFound(format string, a ...interface{}) *ErrNotFound {
	return &ErrNotFound{
		msg: fmt.Sprintf(format, a...),
	}
}

// IsErrNotFound returns true if error is ErrNotFound
func IsErrNotFound(err error) bool {
	err = errors.Cause(err)

	_, ok := err.(*ErrNotFound)

	return ok
}
