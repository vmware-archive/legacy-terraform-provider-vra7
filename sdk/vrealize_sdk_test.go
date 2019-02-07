package sdk

// import (
// 	"errors"
// 	"fmt"
// 	"testing"
//
// 	"gopkg.in/jarcoal/httpmock.v1"
// )
//
// func init() {
// 	fmt.Println("init")
// 	// These are mock test credentials
// 	client = NewClient(
// 		"admin@myvra.local",
// 		"pass!@#",
// 		"vsphere.local",
// 		"http://localhost/",
// 		true,
// 	)
// }
//
// var catalogItemID1 = "e5dd4fba-45ed-4943-b1fc-7f96239286be"
// var catalogItemID2 = "e5dd4fba-45ed-4943-b1fc-7f96239286b1"
//
// var entitledCatalogItemViewsResp = `{"links":[{"@type":"link","rel":"next",
// "href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItemViews?page=2&limit=20"}],
// "content":[{"@type":"ConsumerEntitledCatalogItemView","entitledOrganizations":[{"tenantRef":"vsphere.local",
// "tenantLabel":"vsphere.local","subtenantRef":"53619006-56bb-4788-9723-9eab79752cc1","subtenantLabel":"Content"}],
// "catalogItemId":"502efc1b-d5ce-4ef9-99ee-d4e2a741747c","name":"CentOS 6.3 - IPAM EXT",
// "description":"CentOS 6.3 IaaS Blueprint w/Infoblox IPAM","isNoteworthy":false,"dateCreated":"2016-09-26T13:42:51.564Z",
// "lastUpdatedDate":"2017-01-06T05:11:51.682Z","links":[{"@type":"link","rel":"GET: Request Template",
// "href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItems/502efc1b-d5ce-4ef9-99ee-d4e2a741747c/requests/template{?businessGroupId,requestedFor}"},
// {"@type":"link","rel":"POST: Submit Request",
// "href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItems/502efc1b-d5ce-4ef9-99ee-d4e2a741747c/requests{?businessGroupId,requestedFor}"}],
// "iconId":"502efc1b-d5ce-4ef9-99ee-d4e2a741747c","catalogItemTypeRef":{"id":"com.vmware.csp.component.cafe.composition.blueprint",
// "label":"Composite Blueprint"},"serviceRef":{"id":"baad0ad2-8b96-4347-b188-f534dad53a0d","label":"Infrastructure"},
// "outputResourceTypeRef":{"id":"composition.resource.type.deployment","label":"Deployment"}}],"metadata":{"size":20,
// "totalElements":44,"totalPages":3,"number":1,"offset":0}}`
//
// var entitledCatalogItemViewsErrorResp = `{"errors":[{"code":20116,"source":null,
// "message":"Unable to find the specified catalog item in the service catalog: e5dd4fba-45ed-4943-b1fc-07f96239286b.",
// "systemMessage":"Unable to find the specified catalog item in the service catalog: e5dd4fba-45ed-4943-b1fc-07f96239286b.",
// "moreInfoUrl":null}]}`
//
// var catalogItemViewsErrorResp = `{"errors":[{"code":50505,"source":null,"message":"System exception.",
// "systemMessage":null,"moreInfoUrl":null}]}`
//
// var catalogItemTemplateResp = `{"catalogItem":{"callbacks":null,"catalogItemTypeRef":{"id":"com.vmware.csp.component.cafe.composition.blueprint",
// "label":"Composite Blueprint"},"dateCreated":"2015-12-22T03:16:19.289Z","description":"CentOS 6.3 IaaS Blueprint",
// "forms":{"itemDetails":{"type":"external","formId":"composition.catalog.item.details"},"catalogRequestInfoHidden":true,
// "requestFormScale":"BIG","requestSubmission":{"type":"extension","extensionId":"com.vmware.vcac.core.design.blueprints.requestForm",
// "extensionPointId":null},"requestDetails":{"type":"extension","extensionId":"com.vmware.vcac.core.design.blueprints.requestDetailsForm",
// "extensionPointId":null},"requestPreApproval":null,"requestPostApproval":null},"iconId":"e5dd4fba-45ed-4943-b1fc-7f96239286be",
// "id":"e5dd4fba-45ed-4943-b1fc-7f96239286be","isNoteworthy":false,"lastUpdatedDate":"2017-01-06T05:12:56.690Z",
// "name":"CentOS 6.3","organization":{"tenantRef":"vsphere.local","tenantLabel":"vsphere.local","subtenantRef":null,
// "subtenantLabel":null},"outputResourceTypeRef":{"id":"composition.resource.type.deployment","label":"Deployment"},
// "providerBinding":{"bindingId":"vsphere.local!::!CentOS63","providerRef":{"id":"2fbaabc5-3a48-488a-9f2a-a42616345445",
// "label":"Blueprint Service"}},"serviceRef":{"id":"baad0ad2-8b96-4347-b188-f534dad53a0d","label":"Infrastructure"},
// "status":"PUBLISHED","statusName":"Published","quota":0,"version":4,"requestable":true},"entitledOrganizations":[{"tenantRef":"vsphere.local",
// "tenantLabel":"vsphere.local","subtenantRef":"53619006-56bb-4788-9723-9eab79752cc1","subtenantLabel":"Content"}]}`
//
// func TestFetchCatalogItemByName(t *testing.T) {
// 	httpmock.Activate()
// 	defer httpmock.DeactivateAndReset()
//
// 	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItemViews",
// 		httpmock.NewStringResponder(200, entitledCatalogItemViewsResp))
//
// 	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItemViews?page=1&limit=44",
// 		httpmock.NewStringResponder(200, entitledCatalogItemViewsResp))
//
// 	//Fetch catalog Item by correct and exact name
// 	catalogItemID, err := client.readCatalogItemIDByName("CentOS 6.3 - IPAM EXT")
//
// 	if err != nil {
// 		t.Errorf("Error while fetching catalog Item %v", err)
// 	}
//
// 	if catalogItemID == "" {
// 		t.Errorf("Catalog ID is nil")
// 	}
//
// 	// Fetch catalog item by false name
// 	// The name is not present in catalog item list
// 	// This should return an error
// 	catalogItemID, err = client.readCatalogItemIDByName("Cent OS 6.3")
//
// 	if catalogItemID != "" {
// 		t.Errorf("Catalog Item ID is not nil")
// 	}
//
// 	if err != nil {
// 		t.Errorf("Error while fetching catalog item %v", err)
// 	}
//
// 	// Fetch catalog item by correct and incomplete name
// 	// The name provided (CentOS) to the readCatalogItemIDByName() is substring of correct name (CentOS_6.3)
// 	// This should return empty catalogItemID with suggestions with full name
// 	catalogItemID, err = client.readCatalogItemIDByName("CentOS")
//
// 	if catalogItemID != "" {
// 		t.Errorf("Catalog Item ID is not nil")
// 	}
//
// 	if err == nil {
// 		t.Errorf("Error should have been occurred while fetching the catalog item with encomplete name")
// 	}
// 	httpmock.Reset()
//
// 	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItemViews",
// 		httpmock.NewStringResponder(404, entitledCatalogItemViewsResp))
//
// 	// Fetch catalog item by correct and incomplete name
// 	// The name provided (CentOS) to the readCatalogItemIDByName() is substring of correct name (CentOS_6.3)
// 	// This should return empty catalogItemID with suggestions with full name
// 	catalogItemID, err = client.readCatalogItemIDByName("CentOS")
//
// 	if catalogItemID != "" {
// 		t.Errorf("Catalog Item ID is not nil")
// 	}
//
// 	if err == nil {
// 		t.Errorf("Error should have been occurred while fetching the catalog item with encomplete name")
// 	}
// 	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItemViews",
// 		httpmock.NewStringResponder(404, entitledCatalogItemViewsErrorResp))
// 	catalogItemID, err = client.readCatalogItemIDByName("CentOS")
//
// 	if err == nil {
// 		t.Errorf("Data fetched with wrong catalog ID")
// 	}
//
// 	if catalogItemID != "" {
// 		t.Errorf("Wrong catalog data got fetched")
// 	}
// 	httpmock.Reset()
//
// 	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItemViews",
// 		httpmock.NewStringResponder(200, entitledCatalogItemViewsResp))
//
// 	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItemViews?page=1&limit=44",
// 		httpmock.NewStringResponder(404, entitledCatalogItemViewsErrorResp))
//
// 	catalogItemID, err = client.readCatalogItemIDByName("CentOS")
//
// 	if err == nil {
// 		t.Errorf("Data fetched with wrong catalog item ID")
// 	}
//
// 	if catalogItemID != "" {
// 		t.Errorf("Wrong catalog item data got fetched")
// 	}
//
// 	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItemViews",
// 		httpmock.NewErrorResponder(errors.New(catalogItemViewsErrorResp)))
// 	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItemViews",
// 		httpmock.NewStringResponder(500, entitledCatalogItemViewsErrorResp))
// 	catalogItemID, err = client.readCatalogItemIDByName("CentOS")
// }
//
// // This unit test function contains test cases related to catalog ID in tf config files
// // Scenarios :
// // 1) Fetch catalog item details using correct catalog item ID
// // 2) Fetch catalog item details with invalid catalog item ID
// func TestFetchCatalogItemByID(t *testing.T) {
// 	httpmock.Activate()
// 	defer httpmock.DeactivateAndReset()
// 	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItems/"+catalogItemID1,
// 		httpmock.NewStringResponder(200, catalogItemTemplateResp))
//
// 	catalogItemName, err := client.readCatalogItemNameByID(catalogItemID1)
//
// 	if err != nil {
// 		t.Errorf("Error while fetching catalog item %v", err)
// 	}
//
// 	if catalogItemName == "" {
// 		t.Errorf("Catalog Item Name is is nil")
// 	}
//
// 	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItems/"+catalogItemID1,
// 		httpmock.NewStringResponder(404, entitledCatalogItemViewsErrorResp))
//
// 	catalogItemName, err = client.readCatalogItemNameByID(catalogItemID2)
//
// 	if err == nil {
// 		t.Errorf("Data fetched with wrong catalog ID")
// 	}
//
// 	if catalogItemName != "" {
// 		t.Errorf("Wrong catalog item data got fetched")
// 	}
//
// 	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItems/"+catalogItemID2,
// 		httpmock.NewStringResponder(404, entitledCatalogItemViewsErrorResp))
// 	catalogItemName, err = client.readCatalogItemNameByID(catalogItemID2)
//
// 	if err == nil {
// 		t.Errorf("Data fetched with wrong catalog ID")
// 	}
//
// 	if catalogItemName != "" {
// 		t.Errorf("Wrong catalog data got fetched")
// 	}
// }
