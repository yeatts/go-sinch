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

// New creates a new Client with an embedded http.Client
func New() *Client {
	return &Client{
		baseURL: USBaseURLv1,
		httpClient: &http.Client{
			Transport: &http.Transport{},
		},
	}
}

// US is a shortcut for New().WithCustomBaseURL(USBaseURLv1)
func (client *Client) US() *Client {
	defer client.transaction()()
	client.baseURL = USBaseURLv1
	return client
}

// EU is a shortcut for New().WithCustomBaseURL(EUBaseURLv1)
func (client *Client) EU() *Client {
	defer client.transaction()()
	client.baseURL = EUBaseURLv1
	return client
}

// AU is a shortcut for New().WithCustomBaseURL(AUBaseURLv1)
func (client *Client) AU() *Client {
	defer client.transaction()()
	client.baseURL = AUBaseURLv1
	return client
}

// BR is a shortcut for New().WithCustomBaseURL(BRBaseURLv1)
func (client *Client) BR() *Client {
	defer client.transaction()()
	client.baseURL = BRBaseURLv1
	return client
}

// CA is a shortcut for New().WithCustomBaseURL(CABaseURLv1)
func (client *Client) CA() *Client {
	defer client.transaction()()
	client.baseURL = CABaseURLv1
	return client
}

// WithAuthToken sets the auth token for the client
func (client *Client) WithAuthToken(authToken string) *Client {
	defer client.transaction()()
	client.authToken = authToken
	return client
}

// WithPlanID sets the plan ID for the client
func (client *Client) WithPlanID(planID string) *Client {
	defer client.transaction()()
	client.planID = planID
	return client
}

// WithCustomBaseURL sets the base URL for the client
func (client *Client) WithCustomBaseURL(baseURL string) *Client {
	defer client.transaction()()
	client.baseURL = baseURL
	return client
}

func (client *Client) WithCustomHTTPClient(httpClient *http.Client) *Client {
	defer client.transaction()()
	client.httpClient = httpClient
	return client
}

// transaction locks the client and returns a function that unlocks it. Common usage pattern is to defer a call to
// the returned function.
// Example:
// 	defer client.transaction()()
func (client *Client) transaction() func() {
	client.mu.Lock()
	return client.mu.Unlock
}

// Execute executes the given request with the client's http.Client and returns the response object.
func (client *Client) Execute(req *http.Request, resourceName string) (*http.Response, error) {
	defer client.transaction()()
	req.Header.Set("Authorization", "Bearer "+client.authToken)
	req.Header.Set("Content-Type", "application/json")
	req.URL.Path = fmt.Sprintf("%s/%s", client.baseURL, resourceName)
	return client.httpClient.Do(req)
}
