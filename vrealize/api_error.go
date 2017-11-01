package vrealize

import "fmt"

//APIError struct is used to store REST call errors
type APIError struct {
	Errors []struct {
		Code          int    `json:"code"`
		Message       string `json:"message"`
		SystemMessage string `json:"systemMessage"`
	} `json:"errors"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("vRealize API: %+v", e.Errors)
}

func (e APIError) isEmpty() bool {
	return len(e.Errors) == 0
}
