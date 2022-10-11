package sinch

import (
	"fmt"

	"go.uber.org/multierr"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type Errors []error

func (e Errors) Error() string {
	if len(e) == 0 {
		return ""
	}
	return multierr.Combine(e...).Error()
}

const (
	NoAuthTokenError          = Error("an auth token is required")
	NoBaseURLError            = Error("a base URL is required")
	NilHTTPClientError        = Error("an HTTP client is required")
	UnexpectedStatusCodeError = Error("unexpected status code")
	NilValidatableError       = Error("validatable cannot be nil")
	NilClientError            = Error("client cannot be nil")
	InvalidRequestTypeError   = Error("invalid request type")
)

func UnexpectedStatusCodeErr(exp, actual int) error {
	return multierr.Combine(UnexpectedStatusCodeError, fmt.Errorf("expected %d, got %d", exp, actual))
}
