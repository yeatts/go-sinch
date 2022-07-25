package sms

import (
	"fmt"
	"net/http"
	"sync"
)

type Client struct {
	mu         sync.Mutex
	authToken  string
	baseURL    string
	planID     string
	httpClient *http.Client
}

const (
	BaseURLv1   = "sms.api.sinch.com/xms/v1/"
	USBaseURLv1 = "https://us" + BaseURLv1
	EUBaseURLv1 = "https://eu" + BaseURLv1
	AUBaseURLv1 = "https://au" + BaseURLv1
	BRBaseURLv1 = "https://br" + BaseURLv1
	CABaseURLv1 = "https://cn" + BaseURLv1
)

type Executable interface {
	US() Executable
	EU() Executable
	AU() Executable
	BR() Executable
	CA() Executable
	NewSendRequest() Sendable
	WithAuthToken(authToken string) Executable
	WithPlanID(planID string) Executable
	WithCustomBaseURL(baseURL string) Executable
	WithCustomHTTPClient(httpClient *http.Client) Executable
	Execute(req *http.Request, resourceName string) (*http.Response, error)
}

// NewClient creates a new Client with an embedded http.Client that implements the Executable interface
func NewClient() Executable {
	return &Client{
		baseURL: USBaseURLv1,
		httpClient: &http.Client{
			Transport: &http.Transport{},
		},
	}
}

// US is a shortcut for New().WithCustomBaseURL(USBaseURLv1)
func (c *Client) US() Executable {
	defer c.transaction()()
	c.baseURL = USBaseURLv1
	return c
}

// EU is a shortcut for New().WithCustomBaseURL(EUBaseURLv1)
func (c *Client) EU() Executable {
	defer c.transaction()()
	c.baseURL = EUBaseURLv1
	return c
}

// AU is a shortcut for New().WithCustomBaseURL(AUBaseURLv1)
func (c *Client) AU() Executable {
	defer c.transaction()()
	c.baseURL = AUBaseURLv1
	return c
}

// BR is a shortcut for New().WithCustomBaseURL(BRBaseURLv1)
func (c *Client) BR() Executable {
	defer c.transaction()()
	c.baseURL = BRBaseURLv1
	return c
}

// CA is a shortcut for New().WithCustomBaseURL(CABaseURLv1)
func (c *Client) CA() Executable {
	defer c.transaction()()
	c.baseURL = CABaseURLv1
	return c
}

// WithAuthToken sets the auth token for the client
func (c *Client) WithAuthToken(authToken string) Executable {
	defer c.transaction()()
	c.authToken = authToken
	return c
}

// WithPlanID sets the plan ID for the client
func (c *Client) WithPlanID(planID string) Executable {
	defer c.transaction()()
	c.planID = planID
	return c
}

// WithCustomBaseURL sets the base URL for the client
func (c *Client) WithCustomBaseURL(baseURL string) Executable {
	defer c.transaction()()
	c.baseURL = baseURL
	return c
}

// WithCustomHTTPClient allows you to set a custom http.Client for the SMS client to use for http requests
func (c *Client) WithCustomHTTPClient(httpClient *http.Client) Executable {
	defer c.transaction()()
	c.httpClient = httpClient
	return c
}

// transaction locks the client and returns a function that unlocks it. Common usage pattern is to defer a call to
// the returned function.
// Example:
// 	defer client.transaction()()
func (c *Client) transaction() func() {
	c.mu.Lock()
	return c.mu.Unlock
}

// Execute executes the given request with the client's http.Client and returns the response object.
func (c *Client) Execute(req *http.Request, resourceName string) (*http.Response, error) {
	defer c.transaction()()
	if c.authToken == "" {
		return nil, NoAuthTokenError
	}
	if c.planID == "" {
		return nil, NoPlanIDError
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)
	req.Header.Set("Content-Type", "application/json")
	req.URL.Path = fmt.Sprintf("%s/%s", c.baseURL, resourceName)
	return c.httpClient.Do(req)
}
