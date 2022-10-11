package numbers

import (
	"github.com/thezmc/go-sinch/pkg/api"
	"github.com/thezmc/go-sinch/pkg/sinch"
)

type Client struct {
	SinchAPI  api.Client
	ProjectID string
}

const (
	BaseURLv1 = "https://numbers.api.sinch.com/v1/projects"
)

func (c *Client) WithSinchAPI(sinchAPI *api.Client) *Client {
	c.SinchAPI = *sinchAPI
	return c
}

func (c *Client) WithProjectID(projectID string) *Client {
	c.ProjectID = projectID
	return c
}

func (c *Client) Validate() error {
	if c.ProjectID == "" {
		return ProjectIDRequiredError
	}
	return nil
}

func (c *Client) URL() string {
	return c.SinchAPI.BaseURL + "/" + c.ProjectID
}

func (c *Client) Do(action sinch.Action[sinch.APIRequest, sinch.APIResponse]) error {
	return c.SinchAPI.Do(c, action.Request(), action.Response())
}
