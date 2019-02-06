package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	BearerTokenPrefix = "Bearer "
)

type APIResponse struct {
	Headers    http.Header
	Body       []byte
	Status     string
	StatusCode int
}

type APIRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    io.Reader
}

//AuthResponse - This struct contains response of user authentication call.
type AuthResponse struct {
	Expires time.Time `json:"expires"`
	ID      string    `json:"id"`
	Tenant  string    `json:"tenant"`
}

type AuthenticationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Tenant   string `json:"tenant"`
}

func (ar *APIRequest) AddHeader(key, val string) {
	if ar.Headers == nil {
		ar.Headers = make(map[string]string)
	}
	ar.Headers[key] = val
}

func (ar *APIRequest) ContentType() string {
	if ar.Headers == nil {
		return ""
	}

	contentType, ok := ar.Headers[ContentTypeHeader]
	if !ok {
		return ""
	}
	return contentType
}

func (ar *APIRequest) CopyHeadersTo(req *http.Request) {
	for key, val := range ar.Headers {
		req.Header.Add(key, val)
	}
}

// Represents an error from the Photon API.
type APIError struct {
	Message        string `json:"message"`
	HttpStatusCode int    `json:"-"` // Not part of API contract
}

// Implement Go error interface for ApiError
func (e APIError) Error() string {
	return fmt.Sprintf("Error: { HTTP status: '%v', message: '%v'}",
		e.HttpStatusCode, e.Message)
}

func DoRequest(req *APIRequest, login bool) (*APIResponse, error) {
	r, err := FromApiRequestToHttpRequest(req)

	if !login {
		err = addToken(r)
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
	return FromHttpRespToApiResp(resp)
}

func addToken(req *http.Request) error {
	log.Info("Get Token for the Request to: %v", req.URL)
	var (
		token string
		err   error
	)

	token, err = apiClient.Authenticate()
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return nil
}

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
	req.AddHeader(AcceptHeader, AppJson)
	req.AddHeader(ContentTypeHeader, AppJson)

	log.Info("the request object inside Authenticate is %v ", req)
	return DoLogin(req)
}

func DoLogin(apiReq *APIRequest) (string, error) {
	apiResp, err := DoRequest(apiReq, true)
	if err != nil {
		return "", err
	}

	//The response body of login request using access key and refresh token are different
	//Handle both of these two scenarios
	//This is for the response from access key login
	response := &AuthResponse{}

	err = json.Unmarshal(apiResp.Body, response)
	if err != nil {
		return "", err
	}
	return response.ID, nil
}
