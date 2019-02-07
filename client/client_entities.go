package client

import (
	"io"
	"net/http"
	"time"
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

// Represents an error from the Photon API.
type APIError struct {
	Message        string `json:"message"`
	HttpStatusCode int    `json:"-"` // Not part of API contract
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
