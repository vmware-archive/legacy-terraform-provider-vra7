package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
)

// NewClient creates a new APIClient object
func NewClient(d *schema.ResourceData) APIClient {

	transport := http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: d.Get("insecure").(bool),
	}
	httpClient := &http.Client{
		// Timeout:   clientTimeout,
		Transport: transport,
	}
	apiClient = APIClient{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Tenant:   d.Get("tenant").(string),
		BaseURL:  d.Get("host").(string),
		Insecure: d.Get("insecure").(bool),
		Client:   httpClient,
	}
	return apiClient
}

// DoRequest makes the request and returns the response
func DoRequest(req *APIRequest, login bool) (*APIResponse, error) {
	r, err := FromAPIRequestToHTTPRequest(req)
	if err != nil {
		return nil, err
	}
	if !login {
		err = AddToken(r)
		if err != nil {
			return nil, err
		}
	}
	r.Header.Add(ConnectionHeader, CloseConnection)
	resp, err := apiClient.Client.Do(r)
	if err != nil {
		log.Error("An error occurred when calling %v on %v. Error: %v", req.Method, req.URL, err)
		return nil, err
	}
	log.Info("Check the status of the request %s \n The response is: %s", req.URL, string(resp.Status))
	return FromHTTPRespToAPIResp(resp)
}

// AddToken gets the token and adds to the request header
func AddToken(req *http.Request) error {
	log.Info("Get Token for the Request to: %v", req.URL)
	var (
		token string
		err   error
	)

	token, err = apiClient.Authenticate()
	if err != nil {
		return err
	}
	req.Header.Add(AuthorizationHeader, fmt.Sprintf("Bearer %s", token))
	return nil
}

// Authenticate authenticates for the first time when the provider is invoked
func (apiClient *APIClient) Authenticate() (string, error) {
	uri := fmt.Sprintf("%s/identity/api/tokens", apiClient.BaseURL)
	data := AuthenticationRequest{
		Username: apiClient.Username,
		Password: apiClient.Password,
		Tenant:   apiClient.Tenant,
	}

	jsonData, _ := json.Marshal(data)

	req := &APIRequest{
		Method: POST,
		Body:   bytes.NewBufferString(string(jsonData)),
		URL:    uri,
	}
	req.AddHeader(AcceptHeader, AppJSON)
	req.AddHeader(ContentTypeHeader, AppJSON)

	return DoLogin(req)
}

// DoLogin returns the bearer token
func DoLogin(apiReq *APIRequest) (string, error) {
	apiResp, err := DoRequest(apiReq, true)
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
