package vrealize

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/dghubble/sling"
)

//APIClient - This struct is used to store information provided in .tf file under provider block
//Later on, this stores bearToken after successful authentication and uses that token for next
//REST get or post calls.
type APIClient struct {
	Username    string
	Password    string
	BaseURL     string
	Tenant      string
	Insecure    bool
	BearerToken string
	HTTPClient  *sling.Sling
}

//AuthRequest - This struct contains the user information provided by user
//and for authentication details of this struct are used.
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Tenant   string `json:"tenant"`
}

//AuthResponse - This struct contains response of user authentication call.
type AuthResponse struct {
	Expires time.Time `json:"expires"`
	ID      string    `json:"id"`
	Tenant  string    `json:"tenant"`
}

//ActionResponseTemplate - This struct contains information of blueprint of resource actions.
type ActionResponseTemplate struct {
}

//NewClient - set provider authentication details in struct
//which will be used for all REST call authentication
func NewClient(username string, password string, tenant string, baseURL string, insecure bool) APIClient {
	// This overrides the DefaultTransport which is probably ok
	// since we're generally only using a single client.
	transport := http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: insecure,
	}
	return APIClient{
		Username: username,
		Password: password,
		Tenant:   tenant,
		BaseURL:  baseURL,
		Insecure: insecure,
		HTTPClient: sling.New().Base(baseURL).
			Set("Accept", "application/json").
			Set("Content-Type", "application/json"),
	}
}

//Authenticate - set call for user authentication
func (c *APIClient) Authenticate() error {
	//Set user credentials details as a parameter to authenticate user
	params := &AuthRequest{
		Username: c.Username,
		Password: c.Password,
		Tenant:   c.Tenant,
	}

	authRes := new(AuthResponse)
	apiError := new(APIError)
	//Set a REST call to generate token using above user credentials
	_, err := c.HTTPClient.New().Post("/identity/api/tokens").BodyJSON(params).
		Receive(authRes, apiError)

	if err != nil {
		return err
	}

	if !apiError.isEmpty() {
		log.Errorf("%s\n", apiError.Error())
		return fmt.Errorf("%s", apiError.Error())
	}

	//Get a bearer token
	c.BearerToken = authRes.ID
	//Set bearer token
	c.HTTPClient = c.HTTPClient.New().Set("Authorization",
		fmt.Sprintf("Bearer %s", authRes.ID))

	//Return true on success
	return nil
}
