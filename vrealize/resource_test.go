package vrealize

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/terraform-provider-vra7/sdk"
	"github.com/vmware/terraform-provider-vra7/utils"
)

var (
	mockUser     = "admin@myvra.local"
	mockPassword = "pass!@#"
	mockTenant   = "vsphere.local"
	mockBaseURL  = "http://localhost"
	mockInsecure = true
)

func init() {
	fmt.Println("init")
	client = sdk.NewClient(mockUser, mockPassword, mockTenant, mockBaseURL, mockInsecure)
}

func TestConfigValidityFunction(t *testing.T) {

	mockRequestTemplate := GetMockRequestTemplate()

	// a resource_configuration map is created with valid components
	// all combinations of components name and properties are created with dots
	mockConfigResourceMap := make(map[string]interface{})
	mockConfigResourceMap["mock.test.machine1.cpu"] = 2
	mockConfigResourceMap["mock.test.machine1.mock.storage"] = 8

	resourceSchema := resourceSchema()

	resourceDataMap := map[string]interface{}{
		utils.CatalogItemID:         "abcdefghijklmn",
		utils.ResourceConfiguration: mockConfigResourceMap,
	}

	mockResourceData := schema.TestResourceDataRaw(t, resourceSchema, resourceDataMap)

	p := readProviderConfiguration(mockResourceData)

	readProviderConfiguration(mockResourceData)
	err := p.checkResourceConfigValidity(mockRequestTemplate)
	utils.AssertNilError(t, err)

	mockConfigResourceMap["machine2.mock.cpu"] = 2
	mockConfigResourceMap["machine2.storage"] = 2

	resourceDataMap = map[string]interface{}{
		utils.CatalogItemID:         "abcdefghijklmn",
		utils.ResourceConfiguration: mockConfigResourceMap,
	}

	mockResourceData = schema.TestResourceDataRaw(t, resourceSchema, resourceDataMap)
	readProviderConfiguration(mockResourceData)

	err = p.checkResourceConfigValidity(mockRequestTemplate)
	utils.AssertNilError(t, err)

	mockConfigResourceMap["mock.machine3.vSphere.mock.cpu"] = 2
	resourceDataMap = map[string]interface{}{
		utils.CatalogItemID:         "abcdefghijklmn",
		utils.ResourceConfiguration: mockConfigResourceMap,
	}

	mockResourceData = schema.TestResourceDataRaw(t, resourceSchema, resourceDataMap)
	p = readProviderConfiguration(mockResourceData)

	var mockInvalidKeys []string
	mockInvalidKeys = append(mockInvalidKeys, "mock.machine3.vSphere.mock.cpu")

	validityErr := fmt.Sprintf(ConfigInvalidError, strings.Join(mockInvalidKeys, ", "))
	err = p.checkResourceConfigValidity(mockRequestTemplate)
	// this should throw an error as none of the string combinations (mock, mock.machine3, mock.machine3.vsphere, etc)
	// matches the component names(mock.test.machine1 and machine2) in the request template
	utils.AssertNotNilError(t, err)
	utils.AssertEqualsString(t, validityErr, err.Error())
}

// creates a mock request template from a request template template json file
func GetMockRequestTemplate() *sdk.CatalogItemRequestTemplate {

	mockRequestTemplateStruct := sdk.CatalogItemRequestTemplate{}
	json.Unmarshal([]byte(mockRequestTemplate), &mockRequestTemplateStruct)

	return &mockRequestTemplateStruct

}
