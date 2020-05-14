package cherwell

import (
	"encoding/json"
	"strings"
)

// Constants for Cherwell error codes
const (
	RecordNotFoundError     = "RECORDNOTFOUND"
	BusObNotValidError      = "BusObNotValid"
	BadRequestError         = "BADREQUEST"
	GeneralFailureError     = "GENERALFAILURE"
	ValueNotValidError      = "ValueNotValid"
	ExpressionNotFoundError = "EXPRESSIONNOTFOUND"
	UndefinedError          = "UNDEFINED_LINK_ERROR"
	DuplicateLinkError      = "DUPLICATE_LINK_ERROR"
	DuplicateEntryError     = "DuplicateEntry"
	MarshalError            = "MarshalError"
	UnmarshalError          = "UnmarshalError"
	CreateRequestError      = "CreateRequestError"
	DoRequestError          = "DoRequestError"
	AuthorizationError      = "AuthorizationError"
	ReadResponseError       = "ReadResponseError"
)

// ErrorData holds common error info in response
type ErrorData struct {
	ErrorCode    string `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	HasError     bool   `json:"hasError"`
}

// GetErrorObject gets error object based on error response
func (r *ErrorData) GetErrorObject() error {
	switch r.ErrorCode {
	case BusObNotValidError:
		return &BusObNotValid{Message: r.ErrorMessage}
	case RecordNotFoundError:
		return &RecordNotFound{Message: r.ErrorMessage}
	case GeneralFailureError:
		return &GeneralFailure{Message: r.ErrorMessage}
	}
	return &CherwError{Code: r.ErrorCode, Message: r.ErrorMessage}
}

// CherwError describes a cherwell api error
type CherwError struct {
	Code    string
	Message string
}

func (r ErrorData) Error() string {
	return r.ErrorMessage
}

func (ge *CherwError) Error() string {
	return ge.Message
}

// Errors is a set of errors
type Errors struct {
	Errors []error
}

// NewErrorSet creates a new Errors instance
func NewErrorSet(es ...error) *Errors {
	return &Errors{Errors: es}
}

func (es *Errors) Error() string {
	errs := make([]string, len(es.Errors))
	for i, err := range es.Errors {
		errs[i] = err.Error()
	}

	return strings.Join(errs, "\n")
}

// Add appends an error or set of errors to existing set
func (es *Errors) Add(e ...error) {
	es.Errors = append(es.Errors, e...)
}

// IsEmpty checks if set of errors is not empty
func (es *Errors) IsEmpty() bool {
	return len(es.Errors) == 0
}

// BusObNotValid recognizes cherwell invalid business object error
type BusObNotValid struct {
	Message string
}

func (ge BusObNotValid) Error() string {
	return ge.Message
}

// RecordNotFound recognizes cherwell record not found error
type RecordNotFound struct {
	Message string
}

func (ge RecordNotFound) Error() string {
	return ge.Message
}

// GeneralFailure recognizes cherwell general failure error
type GeneralFailure struct {
	Message string
}

func (ge GeneralFailure) Error() string {
	return ge.Message
}

// InvalidFilterOperator recognizes cherwell invalid filter operator error
type InvalidFilterOperator struct {
	Message string
}

func (ifo InvalidFilterOperator) Error() string {
	return ifo.Message
}

// errorFromResponse builds error based on error response
func errorFromResponse(body string) error {
	errBody := ErrorData{}
	if err := json.Unmarshal([]byte(body), &errBody); err != nil {
		return &GeneralFailure{
			Message: body,
		}
	}

	return errBody.GetErrorObject()
}
