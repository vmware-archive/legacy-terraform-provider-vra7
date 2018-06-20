package vrealize

import (
	"encoding/json"
	"crypto/tls"
	"fmt"
	"log"
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
		log.Printf("%s\n", apiError.Error())
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

//DestroyMachine - To set resource destroy call
func (vRAClient *APIClient) DestroyMachine(destroyTemplate *ActionTemplate, resourceViewTemplate *ResourceView) (*ActionResponseTemplate, error) {
	//Get a destroy template URL from given resource template
	var destroyactionURL string
	destroyactionURL = getactionURL(resourceViewTemplate, "POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.destroy.name}")
	//Raise an error if any exception raised while fetching delete resource URL
	if len(destroyactionURL) == 0 {
		return nil, fmt.Errorf("Resource is not created or not found")
	}

	actionResponse := new(ActionResponseTemplate)
	apiError := new(APIError)

	//Set a REST call with delete resource request and delete resource template as a data
	resp, err := vRAClient.HTTPClient.New().Post(destroyactionURL).
		BodyJSON(destroyTemplate).Receive(actionResponse, apiError)

	if resp.StatusCode != 201 {
		return nil, err
	}

	if !apiError.isEmpty() {
		return nil, apiError
	}

	return actionResponse, nil
}

//PowerOffMachine - To set resource power-off call
func (vRAClient *APIClient) PowerOffMachine(powerOffTemplate *ActionTemplate, resourceViewTemplate *ResourceView) (*ActionResponseTemplate, error) {
	//Get power-off resource URL from given template
	var powerOffMachineactionURL string
	powerOffMachineactionURL = getactionURL(resourceViewTemplate, "POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.PowerOff}")
	//Raise an exception if error got while fetching URL
	if len(powerOffMachineactionURL) == 0 {
		return nil, fmt.Errorf("resource is not created or not found")
	}

	actionResponse := new(ActionResponseTemplate)
	apiError := new(APIError)

	//Set a rest call to power-off the resource with resource power-off template as a data
	response, err := vRAClient.HTTPClient.New().Post(powerOffMachineactionURL).
		BodyJSON(powerOffTemplate).Receive(actionResponse, apiError)

	response.Close = true
	if response.StatusCode == 201 {
		return actionResponse, nil
	}

	if !apiError.isEmpty() {
		return nil, apiError
	}

	if err != nil {
		return nil, err
	}

	return nil, err
}

//GetRequestStatus - To read request status of resource
// which is used to show information to user post create call.
func (vRAClient *APIClient) GetRequestStatus(ResourceID string) (*RequestStatusView, error) {
	//Form a URL to read request status
	path := fmt.Sprintf("catalog-service/api/consumer/requests/%s", ResourceID)
	RequestStatusViewTemplate := new(RequestStatusView)
	apiError := new(APIError)
	//Set a REST call and fetch a resource request status
	_, err := vRAClient.HTTPClient.New().Get(path).Receive(RequestStatusViewTemplate, apiError)
	if err != nil {
		return nil, err
	}
	if !apiError.isEmpty() {
		return nil, apiError
	}
	return RequestStatusViewTemplate, nil
}

// GetDeploymentState - Read the state of a vRA7 Deployment
func (vRAClient *APIClient) GetDeploymentState(CatalogRequestId string) (*ResourceView, error) {
	//Form an URL to fetch resource list view
	path := fmt.Sprintf("catalog-service/api/consumer/requests/%s"+
		"/resourceViews", CatalogRequestId)
	ResourceView := new(ResourceView)
	apiError := new(APIError)
	//Set a REST call to fetch resource view data
	_, err := vRAClient.HTTPClient.New().Get(path).Receive(ResourceView, apiError)
	if err != nil {
		return nil, err
	}
	if !apiError.isEmpty() {
		return nil, apiError
	}
	return ResourceView, nil
}

//RequestMachine - To set create resource REST call
func (vRAClient *APIClient) RequestMachine(template *CatalogItemTemplate) (*RequestMachineResponse, error) {
	//Form a path to set a REST call to create a machine
	path := fmt.Sprintf("/catalog-service/api/consumer/entitledCatalogItems/%s"+
		"/requests", template.CatalogItemID)

	requestMachineRes := new(RequestMachineResponse)
	apiError := new(APIError)

	jsonBody, jErr := json.Marshal(template)
	if jErr != nil {
		log.Printf("Error marshalling template as JSON")
		return nil, jErr
	}

	log.Printf("JSON Request Info: %s", jsonBody)
	//Set a REST call to create a machine
	_, err := vRAClient.HTTPClient.New().Post(path).BodyJSON(template).
		Receive(requestMachineRes, apiError)

	if err != nil {
		return nil, err
	}

	if !apiError.isEmpty() {
		return nil, apiError
	}

	return requestMachineRes, nil
}
