package client

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	logging "github.com/op/go-logging"
	"github.com/vmware/terraform-provider-vra7/utils"
)

var (
	log       = logging.MustGetLogger(utils.LoggerID)
	apiClient APIClient
)

type APIClient struct {
	Username    string
	Password    string
	BaseURL     string
	Tenant      string
	Insecure    bool
	BearerToken string
	Client      *http.Client
}

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

// Reads an error out of the HTTP response, or does nothing if
// no error occured.
func GetApiError(respBody []byte, statusCode int) error {
	apiError := APIError{}
	unmarshalErr := json.Unmarshal(respBody, &apiError)
	if unmarshalErr != nil {
		// Do not return this error just log it.
		log.Error("Error is %v ", unmarshalErr)
		apiError.Message = string(respBody)
	}

	apiError.HttpStatusCode = statusCode
	return apiError
}

func FromApiRequestToHttpRequest(apiReq *APIRequest) (*http.Request, error) {
	req, err := http.NewRequest(apiReq.Method, apiReq.URL, apiReq.Body)
	if err != nil {
		return nil, err
	}

	if apiReq.ContentType() == "" {
		apiReq.AddHeader(ContentTypeHeader, AppJson)
	}
	apiReq.CopyHeadersTo(req)
	return req, nil
}

func FromHttpRespToApiResp(resp *http.Response) (*APIResponse, error) {
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, GetApiError(respBody, resp.StatusCode)
	}

	apiResp := &APIResponse{}
	apiResp.Body = respBody
	apiResp.Headers = resp.Header
	apiResp.Status = resp.Status
	apiResp.StatusCode = resp.StatusCode
	return apiResp, nil
}

func BuildEncodedURL(relativePath string, queryParameters map[string]string) string {
	//Todo it might be better to swith to Viper to load all the config at once
	serverURL := apiClient.BaseURL

	var queryURL *url.URL
	queryURL, err := url.Parse(serverURL)
	if err != nil {
		log.Error("Error %v ", err)
	}

	queryURL.Path += relativePath

	parameters := url.Values{}
	if queryParameters != nil {
		for key, value := range queryParameters {
			parameters.Add(key, value)
		}
	} else {
		return queryURL.String()
	}
	queryURL.RawQuery = parameters.Encode()
	return queryURL.String()
}
