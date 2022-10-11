package sms // import sinchsms "github.com/thezmc/go-sinch/sms"

import (
	"github.com/thezmc/go-sinch/pkg/api"
	"github.com/thezmc/go-sinch/pkg/sinch"
)

type Client struct {
	SinchAPI api.Client
	PlanID   string
}

const (
	BaseURLv1   = ".sms.api.sinch.com/xms/v1"
	USBaseURLv1 = "https://us" + BaseURLv1
	EUBaseURLv1 = "https://eu" + BaseURLv1
	AUBaseURLv1 = "https://au" + BaseURLv1
	BRBaseURLv1 = "https://br" + BaseURLv1
	CABaseURLv1 = "https://cn" + BaseURLv1
)

func (c *Client) US() *Client {
	c.SinchAPI.WithBaseURL(USBaseURLv1)
	return c
}

// EU is a shortcut for New().WithBaseURL(EUBaseURLv1)
func (c *Client) EU() *Client {
	c.SinchAPI.WithBaseURL(EUBaseURLv1)
	return c
}

// AU is a shortcut for New().WithBaseURL(AUBaseURLv1)
func (c *Client) AU() *Client {
	c.SinchAPI.WithBaseURL(AUBaseURLv1)
	return c
}

// BR is a shortcut for New().WithBaseURL(BRBaseURLv1)
func (c *Client) BR() *Client {
	c.SinchAPI.WithBaseURL(BRBaseURLv1)
	return c
}

// CA is a shortcut for New().WithBaseURL(CABaseURLv1)
func (c *Client) CA() *Client {
	c.SinchAPI.WithBaseURL(CABaseURLv1)
	return c
}

// WithPlanID sets the plan ID for the client
func (c *Client) WithPlanID(planID string) *Client {
	c.PlanID = planID
	return c
}

func (c *Client) WithSinchAPI(sinchAPI *api.Client) *Client {
	c.SinchAPI = *sinchAPI
	return c
}

func (c *Client) URL() string {
	return c.SinchAPI.BaseURL + "/" + c.PlanID
}

func (c *Client) Validate() error {
	if c.PlanID == "" {
		return NoPlanIDError
	}
	return nil
}

// Do executes the given request with the client's http.Client and returns the response object.
func (c *Client) Do(req sinch.APIRequest, resp sinch.APIResponse) error {
	return c.SinchAPI.Do(c, req, resp)
}
