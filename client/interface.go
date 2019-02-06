package client

type HttpClient interface {
	DoRequest(req *APIRequest) (*APIResponse, error)
}
