package api

import (
	"bytes"
	"io"
	"net/http"

	"github.com/thezmc/go-sinch/pkg/sinch"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func (api *Client) WithBaseURL(baseURL string) *Client {
	api.BaseURL = baseURL
	return api
}

func (api *Client) WithHTTPClient(httpClient *http.Client) *Client {
	api.HTTPClient = httpClient
	return api
}

func (c Client) Validate() error {
	if c.BaseURL == "" {
		return NoBaseURLError
	}
	if c.HTTPClient == nil {
		return NilHTTPClientError
	}
	return nil
}

func (c Client) Do(client sinch.APIClient, req sinch.APIRequest, recv sinch.APIResponse) error {
	if err := Validate(c, client, req); err != nil {
		return err
	}

	queryString, err := req.QueryString()
	if err != nil {
		return err
	}

	body, err := req.Body()
	if err != nil {
		return err
	}

	url := client.URL() + req.Path() + queryString
	httpReq, err := http.NewRequest(req.Method(), url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	_, err = client.Authenticate(httpReq)
	if err != nil {
		return err
	}

	if httpReq.ContentLength > 0 {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	httpResp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != req.ExpectedStatusCode() {
		return UnexpectedStatusCodeErr(req.ExpectedStatusCode(), httpResp.StatusCode)
	}

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}

	return recv.FromJSON(respBody)
}
