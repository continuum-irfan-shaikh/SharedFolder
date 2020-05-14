package types

import "fmt"

// MultiError - error that contains multiple errors and implements error interface
type MultiError []error

// NewMultiError - creates NewMultiError instance
func NewMultiError() MultiError {
	err := make(MultiError, 0)
	return err
}

// ToError - converts to error, returns nil for empty internal errors
func (m *MultiError) ToError() error {
	if len(*m) == 0 {
		return nil
	}
	return error(m)
}

// Error - returns formatted error message
func (m *MultiError) Error() string {
	var msg string
	for _, err := range *m {
		switch err.(type) {
		case *MultiError:
			msg += err.Error()
		default:
			msg = fmt.Sprintf("%s\n%+v\n", msg, err)
		}
	}
	return msg
}
