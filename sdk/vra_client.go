package sdk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

// NewClient creates a new APIClient object
func NewClient(user, password, tenant, baseURL string, insecure bool) APIClient {

	transport := http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: insecure,
	}
	httpClient := &http.Client{
		// Timeout:   clientTimeout,
		Transport: transport,
	}
	log.Critical("user %s ", user)
	apiClient := APIClient{
		Username: user,
		Password: password,
		Tenant:   tenant,
		BaseURL:  baseURL,
		Insecure: insecure,
		Client:   httpClient,
	}
	return apiClient
}

// DoRequest makes the request and returns the response
func (c *APIClient) DoRequest(req *APIRequest, login bool) (*APIResponse, error) {
	r, err := FromAPIRequestToHTTPRequest(req)
	if err != nil {
		return nil, err
	}
	if !login {
		err = c.AddToken(r)
		if err != nil {
			return nil, err
		}
	}
	r.Header.Add(ConnectionHeader, CloseConnection)
	log.Critical("the request object is %v ", r)
	resp, err := c.Client.Do(r)
	if err != nil {
		log.Error("An error occurred when calling %v on %v. Error: %v", req.Method, req.URL, err)
		return nil, err
	}
	log.Info("Check the status of the request %s \n The response is: %s", req.URL, string(resp.Status))
	return FromHTTPRespToAPIResp(resp)
}

// AddToken gets the token and adds to the request header
func (c *APIClient) AddToken(req *http.Request) error {
	log.Info("Get Token for the Request to: %v", req.URL)
	var (
		token string
		err   error
	)

	token, err = c.Authenticate()
	if err != nil {
		return err
	}
	req.Header.Add(AuthorizationHeader, fmt.Sprintf("Bearer %s", token))
	return nil
}

// Authenticate authenticates for the first time when the provider is invoked
func (c *APIClient) Authenticate() (string, error) {
	uri := fmt.Sprintf("%s"+Tokens, c.BaseURL)
	data := AuthenticationRequest{
		Username: c.Username,
		Password: c.Password,
		Tenant:   c.Tenant,
	}

	jsonData, _ := json.Marshal(data)

	req := &APIRequest{
		Method: POST,
		Body:   bytes.NewBufferString(string(jsonData)),
		URL:    uri,
	}
	req.AddHeader(AcceptHeader, AppJSON)
	req.AddHeader(ContentTypeHeader, AppJSON)

	return c.DoLogin(req)
}

// DoLogin returns the bearer token
func (c *APIClient) DoLogin(apiReq *APIRequest) (string, error) {
	apiResp, err := c.DoRequest(apiReq, true)
	if err != nil {
		return "", err
	}
	response := &AuthResponse{}

	err = json.Unmarshal(apiResp.Body, response)
	if err != nil {
		return "", err
	}
	return response.ID, nil
}
