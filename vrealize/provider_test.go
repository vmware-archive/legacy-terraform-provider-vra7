package vrealize

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/terraform-provider-vra7/sdk"
	"github.com/vmware/terraform-provider-vra7/utils"
	"gopkg.in/jarcoal/httpmock.v1"
	"testing"
)

var (
	client   sdk.APIClient
	user     = "admin@myvra.local"
	password = "pass!@#"
	tenant   = "vsphere.local"
	baseURL  = "http://localhost"
	insecure = true
)

func init() {

	fmt.Println("init")
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
