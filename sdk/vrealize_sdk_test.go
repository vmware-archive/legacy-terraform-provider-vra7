package sdk

import (
	"fmt"
	"testing"

	"github.com/vmware/terraform-provider-vra7/utils"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

var (
	client   APIClient
	user     = "admin@myvra.local"
	password = "pass!@#"
	tenant   = "vsphere.local"
	baseURL  = "http://localhost"
	insecure = true
)

func init() {
	fmt.Println("init")
	client = NewClient(user, password, tenant, baseURL, insecure)
}

func TestGetCatalogItemRequestTemplate(t *testing.T) {
	httpmock.ActivateNonDefault(client.Client)
	defer httpmock.DeactivateAndReset()

	catalogItemID := "feaedf73-560c-4612-a573-41667e017691"

	path := fmt.Sprintf(RequestTemplateAPI, catalogItemID)
	url := client.BuildEncodedURL(path, nil)

	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(200, requestTemplateResponse))

	catalogItemReqTemplate, err := client.GetCatalogItemRequestTemplate(catalogItemID)
	utils.AssertNilError(t, err)
	utils.AssertEqualsString(t, catalogItemID, catalogItemReqTemplate.CatalogItemID)

	catalogItemReqTemplate, err = client.GetCatalogItemRequestTemplate("635e5v-8e37efd60-hdgdh")
	utils.AssertNotNilError(t, err)

	httpmock.Reset()
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(20116, requestTemplateErrorResponse))
	invalidCatalogItemID := "feaedf73-560c-4612-a573-0041667e0176"
	catalogItemReqTemplate, err = client.GetCatalogItemRequestTemplate(invalidCatalogItemID)
	utils.AssertNotNilError(t, err)

}

func TestReadCatalogItemNameByID(t *testing.T) {

	httpmock.ActivateNonDefault(client.Client)
	defer httpmock.DeactivateAndReset()

	catalogItemID := "e5dd4fba-45ed-4943-b1fc-7f96239286be"
	path := fmt.Sprintf(EntitledCatalogItems+"/"+"%s", catalogItemID)
	url := client.BuildEncodedURL(path, nil)

	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(200, catalogItemResp))

	catalogItemName, err := client.ReadCatalogItemNameByID(catalogItemID)
	utils.AssertNilError(t, err)
	utils.AssertEqualsString(t, "CentOS 6.3", catalogItemName)

	catalogItemName, err = client.ReadCatalogItemNameByID("84rg=73dv-dd8dhy-hg")
	utils.AssertNotNilError(t, err)
	utils.AssertEqualsString(t, "", catalogItemName)
}

func TestReadCatalogItemByName(t *testing.T) {
	httpmock.ActivateNonDefault(client.Client)
	defer httpmock.DeactivateAndReset()

	path := fmt.Sprintf(EntitledCatalogItemViewsAPI)
	url := client.BuildEncodedURL(path, nil)

	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(200, entitledCatalogItemViewsResponse))

	catalogItemID, err := client.ReadCatalogItemByName("CentOs")
	utils.AssertEqualsString(t, "feaedf73-560c-4612-a573-41667e017691", catalogItemID)
	utils.AssertNilError(t, err)

	catalogItemID, err = client.ReadCatalogItemByName("Invalid Catalog Item name")
	utils.AssertEqualsString(t, "", catalogItemID)
}

func TestGetBusinessGroupID(t *testing.T) {

	httpmock.ActivateNonDefault(client.Client)
	defer httpmock.DeactivateAndReset()

	path := Tenants + "/" + tenant + "/subtenants"
	url := client.BuildEncodedURL(path, nil)

	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, subTenantsResponse))

	id, err := client.GetBusinessGroupID("Development", tenant)

	if id == "b2470b94-cbca-43db-be37-803cca7b0f1a" {
		fmt.Println("Passed")
	}

	if err != nil {
		t.Errorf("Error fetching is %v ", err)
	}
}
