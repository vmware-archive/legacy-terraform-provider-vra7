package client

import (
	"io"
)

const (
	// Common header key
	ContentTypeHeader   = "Content-Type"
	ContentLengthHeader = "Content-Length"
	AuthorizationHeader = "Authorization"
	ConnectionHeader    = "Connection"
	AcceptHeader        = "Accept"

	// Common header values
	AppJson         = "application/json"
	AppYaml         = "application/yaml"
	CloseConnection = "close"

	// Http Methods
	GET    = "GET"
	POST   = "POST"
	PATCH  = "PATCH"
	PUT    = "PUT"
	DELETE = "DELETE"
)

func Get(encodedUrl string, headers map[string]string) (*APIResponse, error) {
	req := &APIRequest{
		Method:  GET,
		URL:     encodedUrl,
		Headers: headers,
	}
	return doRequest(req)
}

func Post(url string, body io.Reader, headers map[string]string) (*APIResponse, error) {
	req := &APIRequest{
		Method:  POST,
		URL:     url,
		Body:    body,
		Headers: headers,
	}
	return doRequest(req)
}

func Put(url string, body io.Reader, headers map[string]string) (*APIResponse, error) {
	req := &APIRequest{
		Method:  PUT,
		URL:     url,
		Body:    body,
		Headers: headers,
	}
	return doRequest(req)
}

func Patch(url string, body io.Reader, headers map[string]string) (*APIResponse, error) {
	req := &APIRequest{
		Method:  PATCH,
		URL:     url,
		Body:    body,
		Headers: headers,
	}
	return doRequest(req)
}

func Delete(url string, body io.Reader, headers map[string]string) (*APIResponse, error) {
	req := &APIRequest{
		Method:  DELETE,
		URL:     url,
		Body:    body,
		Headers: headers,
	}
	return doRequest(req)
}

func doRequest(req *APIRequest) (*APIResponse, error) {
	apiResp, err := DoRequest(req, false)
	if err != nil {
		return nil, err
	}
	return apiResp, nil
}
