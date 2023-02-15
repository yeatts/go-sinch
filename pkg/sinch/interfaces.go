package sinch

import "net/http"

type Validatable interface {
	Validate() error
}

type APIAction[RQ APIRequest, RS APIResponse] interface {
	Request() RQ
	Response() RS
}

type APIRequest interface {
	Validatable
	ExpectedStatusCode() int
	Method() string
	QueryString() (string, error)
	Body() ([]byte, error)
	Path() string
}

type APIResponse interface {
	FromJSON([]byte) error
}

type APIClient interface {
	Validatable
	Authenticate(*http.Request) (*http.Request, error)
	URL() string
	Do(APIRequest, APIResponse) error
}

type API interface {
	Do(client APIClient, req APIRequest, recv APIResponse) error
}

type Action[RQ APIRequest, RS APIResponse] APIAction[RQ, RS]
