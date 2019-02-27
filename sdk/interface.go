package sdk

// HTTPClient interface
type HTTPClient interface {
	DoRequest(req *APIRequest) (*APIResponse, error)
}
