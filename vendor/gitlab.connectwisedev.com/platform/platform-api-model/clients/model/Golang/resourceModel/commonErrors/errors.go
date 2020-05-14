package commonErrors

import (
	"fmt"
	"github.com/pkg/errors"
	"path/filepath"
	"runtime"
	"strings"
)

type baseError struct {
	err        error
	fileName   string
	lineNumber int
}

func (err baseError) Error() string {
	return fmt.Sprintf("%s:%d: %s", err.fileName, err.lineNumber, err.err.Error())
}

//ResourceAlreadyExistError is StatusConflict (if the resource already exists) (409) error
type ResourceAlreadyExistError struct {
	baseError
}

//BadRequestError is BadRequest (400) error
type BadRequestError struct {
	baseError
}

//NotFoundError is NotFound (404) error
type NotFoundError struct {
	baseError
}

//AccessDeniedError  is accessDenied (403) error
type AccessDeniedError struct {
	baseError
}

//InternalServerError is InternalServerError(500) error
type InternalServerError struct {
	baseError
}

//NewResourceAlreadyExistError...
func NewResourceAlreadyExistError(errMsg string) error {
	return ResourceAlreadyExistError{
		baseError: newErr(errMsg),
	}
}

//NewBadRequestError...
func NewBadRequestError(errMsg string) error {
	return BadRequestError{
		baseError: newErr(errMsg),
	}
}

//NewNotFoundError...
func NewNotFoundError(errMsg string) error {
	return NotFoundError{
		baseError: newErr(errMsg),
	}
}

//NewAccessDeniedError...
func NewAccessDeniedError(errMsg string) error {
	return AccessDeniedError{
		baseError: newErr(errMsg),
	}
}

//NewInternalServerError...
func NewInternalServerError(errMsg string) error {
	return InternalServerError{
		baseError: newErr(errMsg),
	}
}

func newErr(errMsg string) baseError {
	e := baseError{
		err: errors.New(errMsg),
	}
	e.fileName, e.lineNumber = getFileNameAndLine()
	return e
}

func getFileNameAndLine() (fileName string, lineNumber int) {
	_, file, line, ok := runtime.Caller(3)
	if ok {
		filePathSlice := strings.Split(file, string(filepath.Separator))
		if len(filePathSlice) == 0 {
			fileName = "undefined"
			return
		}

		fileName = filePathSlice[len(filePathSlice)-1]
		lineNumber = line
	}
	return
}
