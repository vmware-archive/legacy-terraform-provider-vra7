package vrealize

// import (
// 	"errors"
// 	"fmt"
// 	"github.com/hashicorp/terraform/helper/schema"
// 	"github.com/hashicorp/terraform/terraform"
// 	"gopkg.in/jarcoal/httpmock.v1"
// 	"testing"
// )
//
// var testProviders map[string]terraform.ResourceProvider
// var testProvider *schema.Provider
//
// func init() {
// 	testProvider = Provider().(*schema.Provider)
// 	testProviders = map[string]terraform.ResourceProvider{
// 		"vra7": testProvider,
// 	}
//
// 	t := new(testing.T)
// 	fmt.Println("init")
// 	client = NewClient(
// 		"admin@myvra.local",
// 		"pass!@#",
// 		"vsphere.local",
// 		"http://localhost/",
// 		true,
// 	)
// 	httpmock.Activate()
// 	defer httpmock.DeactivateAndReset()
//
// 	httpmock.RegisterResponder("POST", "http://localhost/identity/api/tokens",
// 		httpmock.NewStringResponder(200, `{"expires":"2017-07-25T15:18:49.000Z",
// 		"id":"MTUwMDk2NzEyOTEyOTplYTliNTA3YTg4MjZmZjU1YTIwZjp0ZW5hbnQ6dnNwaGVyZS5sb2NhbHVzZX
// 		JuYW1lOmphc29uQGNvcnAubG9jYWxleHBpcmF0aW9uOjE1MDA5OTU5MjkwMDA6ZjE1OTQyM2Y1NjQ2YzgyZjY
// 		4Yjg1NGFjMGNkNWVlMTNkNDhlZTljNjY3ZTg4MzA1MDViMTU4Y2U3MzBkYjQ5NmQ5MmZhZWM1MWYzYTg1ZWM4
// 		ZDhkYmFhMzY3YTlmNDExZmM2MTRmNjk5MGQ1YjRmZjBhYjgxMWM0OGQ3ZGVmNmY=","tenant":"vsphere.local"}`))
//
// 	client.Authenticate()
//
// 	httpmock.RegisterResponder("POST", "http://localhost/identity/api/tokens",
// 		httpmock.NewErrorResponder(errors.New(`{"errors":[{"code":90135,"source":null,"message":"Unable to authenticate user jason@corp.local1 in tenant vsphere.local.","systemMessage":"90135-Unable to authenticate user jason@corp.local1 in tenant vsphere.local.","moreInfoUrl":null}]}`)))
//
// 	err := client.Authenticate()
// 	if err == nil {
// 		t.Errorf("Authentication should fail")
// 	}
// }
//
// func TestValidateProvider(t *testing.T) {
// 	httpmock.RegisterResponder("POST", "http://localhost/identity/api/tokens",
// 		httpmock.NewErrorResponder(errors.New(`{"errors":[{"code":90135,"source":null,"message":"Unable to authenticate user jason@corp.local1 in tenant vsphere.local.","systemMessage":"90135-Unable to authenticate user jason@corp.local1 in tenant vsphere.local.","moreInfoUrl":null}]}`)))
//
// 	err := client.Authenticate()
// 	if err == nil {
// 		t.Errorf("Authentication should fail")
// 	}
// }
//
// func TestProvider(t *testing.T) {
// 	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
// 		t.Fatalf("err: %s", err)
// 	}
// }
