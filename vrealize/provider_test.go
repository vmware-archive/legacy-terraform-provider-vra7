package vrealize

import (
	"github.com/vmware/terraform-provider-vra7/utils"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"gopkg.in/jarcoal/httpmock.v1"
	"github.com/vmware/terraform-provider-vra7/sdk"
	"testing"
)

var (
	testProviders map[string]terraform.ResourceProvider
    testProvider *schema.Provider
	client   sdk.APIClient
	user     = "admin@myvra.local"
	password = "pass!@#"
	tenant   = "vsphere.local"
	baseURL  = "http://localhost"
	insecure = true

	validAuthResponse =
	`{  
		"expires":"2019-02-26T03:32:35.000Z",
		"id":"MTU1MTEyMzE1NTc5ODpiYTZkYjdhNjZlNGNkYjZmZTBiMjp0ZW5hbnQ6cWV1c2VybmFtZTpmcml0ekBjb2tlLnNxYS1ob3Jpem9uLmxvY2FsZXhwaXJhdGlvbjoxNTUxMTUxOTU1MDAwOmMyNGVjNTFiNzE1OTJhZDZjNTljMTUwMDkxMjcyNzUyZDkzNzQ0ODRkMTVlZGFhNWM0MDhjYmQ3YTM2MTljZGNiNjM3MjM1NmY1MzZlYTk1YzUyMGZiZDVjMTkzMzg3YjQzZmMwNmNlMGI5YjJkZmIwNzhlZGU2NzdiNTk3MWFk",
		"tenant":"qe"
	 }`

	 errorAuthResponse =
	 `{  
		"errors":[  
		   {  
			  "code":90135,
			  "source":null,
			  "message":"Unable to authenticate user fritz@coke.sqa-horizon.local in tenant q.",
			  "systemMessage":"90135-Unable to authenticate user fritz@coke.sqa-horizon.local in tenant q.",
			  "moreInfoUrl":null
		   }
		]
	 }`
)

func init() {

	fmt.Println("init")
	testProvider = Provider().(*schema.Provider)
	testProviders = map[string]terraform.ResourceProvider{
		"vra7": testProvider,
	}
	client = sdk.NewClient(user, password, tenant, baseURL, insecure)
}

func TestValidateProvider(t *testing.T) {
	httpmock.ActivateNonDefault(client.Client)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost/identity/api/tokens",
		httpmock.NewStringResponder(200, validAuthResponse))

	err := client.Authenticate()
	utils.AssertNilError(t, err)

	httpmock.RegisterResponder("POST", "http://localhost/identity/api/tokens",
		httpmock.NewStringResponder(90135, errorAuthResponse))

	err = client.Authenticate()
	utils.AssertNotNilError(t, err)
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
