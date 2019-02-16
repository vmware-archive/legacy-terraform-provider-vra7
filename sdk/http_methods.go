package sdk

import (
	"io"
)

// Get HTTP GET method
func (c *APIClient) Get(encodedURL string, headers map[string]string) (*APIResponse, error) {
	req := &APIRequest{
		Method:  GET,
		URL:     encodedURL,
		Headers: headers,
	}
	return c.doRequest(req)
}

// Post HTTP POST method
func (c *APIClient) Post(url string, body io.Reader, headers map[string]string) (*APIResponse, error) {
	req := &APIRequest{
		Method:  POST,
		URL:     url,
		Body:    body,
		Headers: headers,
	}
	return c.doRequest(req)
}

// Put HTTP PUT method
func (c *APIClient) Put(url string, body io.Reader, headers map[string]string) (*APIResponse, error) {
	req := &APIRequest{
		Method:  PUT,
		URL:     url,
		Body:    body,
		Headers: headers,
	}
	return c.doRequest(req)
}

// Patch HTTP PATCH method
func (c *APIClient) Patch(url string, body io.Reader, headers map[string]string) (*APIResponse, error) {
	req := &APIRequest{
		Method:  PATCH,
		URL:     url,
		Body:    body,
		Headers: headers,
	}
	return c.doRequest(req)
}

// Delete HTTP DELETE method
func (c *APIClient) Delete(url string, body io.Reader, headers map[string]string) (*APIResponse, error) {
	req := &APIRequest{
		Method:  DELETE,
		URL:     url,
		Body:    body,
		Headers: headers,
	}
	return c.doRequest(req)
}

// doRequest makes the request and returns the response
func (c *APIClient) doRequest(req *APIRequest) (*APIResponse, error) {
	apiResp, err := c.DoRequest(req, false)
	if err != nil {
		return nil, err
	}
	return apiResp, nil
}
