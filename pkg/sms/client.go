package sms // import sinchsms "github.com/thezmc/go-sinch/sms"

import (
	"net/http"
	"sync"

	"github.com/thezmc/go-sinch/pkg/interfaces"
)

type client struct {
	mutex      sync.Mutex
	authToken  string
	baseURL    string
	planID     string
	httpClient *http.Client
}

const (
	BaseURLv1   = ".sms.api.sinch.com/xms/v1"
	USBaseURLv1 = "https://us" + BaseURLv1
	EUBaseURLv1 = "https://eu" + BaseURLv1
	AUBaseURLv1 = "https://au" + BaseURLv1
	BRBaseURLv1 = "https://br" + BaseURLv1
	CABaseURLv1 = "https://cn" + BaseURLv1
)

// NewClient creates a new Client with an embedded http.Client that implements the Executable interface
func NewClient() interfaces.SMSClient {
	return &client{
		baseURL:    USBaseURLv1,
		httpClient: new(http.Client),
	}
}

// US is a shortcut for New().WithBaseURL(USBaseURLv1). The base URL is set to the US base URL by default.
func (c *client) US() interfaces.SMSClient {
	defer c.transact()()
	c.baseURL = USBaseURLv1
	return c
}

// EU is a shortcut for New().WithBaseURL(EUBaseURLv1)
func (c *client) EU() interfaces.SMSClient {
	defer c.transact()()
	c.baseURL = EUBaseURLv1
	return c
}

// AU is a shortcut for New().WithBaseURL(AUBaseURLv1)
func (c *client) AU() interfaces.SMSClient {
	defer c.transact()()
	c.baseURL = AUBaseURLv1
	return c
}

// BR is a shortcut for New().WithBaseURL(BRBaseURLv1)
func (c *client) BR() interfaces.SMSClient {
	defer c.transact()()
	c.baseURL = BRBaseURLv1
	return c
}

// CA is a shortcut for New().WithBaseURL(CABaseURLv1)
func (c *client) CA() interfaces.SMSClient {
	defer c.transact()()
	c.baseURL = CABaseURLv1
	return c
}

// WithAuthToken sets the auth token for the client
func (c *client) WithAuthToken(authToken string) interfaces.SMSClient {
	defer c.transact()()
	c.authToken = authToken
	return c
}

// WithPlanID sets the plan ID for the client
func (c *client) WithPlanID(planID string) interfaces.SMSClient {
	defer c.transact()()
	c.planID = planID
	return c
}

// WithCustomBaseURL sets the base URL for the client
func (c *client) WithBaseURL(baseURL string) interfaces.SMSClient {
	defer c.transact()()
	c.baseURL = baseURL
	return c
}

// WithCustomHTTPClient allows you to set a custom http.Client for the SMS client to use for http requests
func (c *client) WithCustomHTTPClient(httpClient *http.Client) interfaces.SMSClient {
	defer c.transact()()
	c.httpClient = httpClient
	return c
}

// transact locks the client and returns a function that unlocks it. Common usage pattern is to defer a call to
// the returned function.
// Example:
//
//	defer client.transact()()
func (c *client) transact() func() {
	c.mutex.Lock()
	return c.mutex.Unlock
}

// Execute executes the given request with the client's http.Client and returns the response object.
func (c *client) Execute(req *http.Request, resourceName string, queryString string) (*http.Response, error) {
	defer c.transact()()
	if c.authToken == "" {
		return nil, NoAuthTokenError
	}
	if c.planID == "" {
		return nil, NoPlanIDError
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)
	req.Header.Set("Content-Type", "application/json")
	return c.httpClient.Do(req)
}
