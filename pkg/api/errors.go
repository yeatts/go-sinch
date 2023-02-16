package api

import (
	"fmt"

	"github.com/thezmc/go-sinch/pkg/sinch"
	"go.uber.org/multierr"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type Errors []error

func (e Errors) Error() string {
	return multierr.Combine(e...).Error()
}

const (
	NoBaseURLError            = Error("a base URL is required")
	NilHTTPClientError        = Error("an HTTP client is required")
	UnexpectedStatusCodeError = Error("unexpected status code")
	NilValidatableError       = Error("validatable cannot be nil")
	NilClientError            = Error("client cannot be nil")
	InvalidRequestTypeError   = Error("invalid request type")
)

func UnexpectedStatusCodeErr(exp, actual int) error {
	return Errors{
		UnexpectedStatusCodeError,
		fmt.Errorf("expected %d, got %d", exp, actual),
	}
}

type RequestValidator func(req sinch.APIRequest) error
type ClientValidator func(client sinch.APIClient) RequestValidator
