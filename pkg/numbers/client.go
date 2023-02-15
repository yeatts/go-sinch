package numbers

import (
	"net/http"

	"github.com/thezmc/go-sinch/pkg/api"
	"github.com/thezmc/go-sinch/pkg/sinch"
)

type Client struct {
	SinchAPI  api.Client
	ProjectID string
	KeyID     string
	KeySecret string
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

func (c *Client) WithKeyID(keyID string) *Client {
	c.KeyID = keyID
	return c
}

func (c *Client) WithKeySecret(KeySecret string) *Client {
	c.KeySecret = KeySecret
	return c
}

func (c *Client) Validate() error {
	if c.ProjectID == "" {
		return ProjectIDRequiredError
	}
	if c.KeyID == "" {
		return KeyIDRequiredError
	}
	if c.KeySecret == "" {
		return KeySecretRequiredError
	}
	return nil
}

func (c *Client) Authenticate(httpReq *http.Request) (*http.Request, error) {
	httpReq.SetBasicAuth(c.KeyID, c.KeySecret)
	return httpReq, nil
}

func (c *Client) URL() string {
	return c.SinchAPI.BaseURL + "/" + c.ProjectID
}

func (c *Client) Do(req sinch.APIRequest, resp sinch.APIResponse) error {
	return c.SinchAPI.Do(c, req, resp)
}
