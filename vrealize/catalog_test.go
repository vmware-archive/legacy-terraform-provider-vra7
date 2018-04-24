package vrealize

import (
	"errors"
	"fmt"
	"gopkg.in/jarcoal/httpmock.v1"
	"testing"
)

func init() {
	fmt.Println("init")
	client = NewClient(
		"admin@myvra.local",
		"pass!@#",
		"vsphere.local",
		"http://localhost/",
		true,
	)
}

func TestFetchCatalogByName(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItemViews",
		httpmock.NewStringResponder(200, `{"links":[{"@type":"link","rel":"next","href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItemViews?page=2&limit=20"}],"content":[{"@type":"ConsumerEntitledCatalogItemView","entitledOrganizations":[{"tenantRef":"vsphere.local","tenantLabel":"vsphere.local","subtenantRef":"53619006-56bb-4788-9723-9eab79752cc1","subtenantLabel":"Content"}],"catalogItemId":"502efc1b-d5ce-4ef9-99ee-d4e2a741747c","name":"CentOS 6.3 - IPAM EXT","description":"CentOS 6.3 IaaS Blueprint w/Infoblox IPAM","isNoteworthy":false,"dateCreated":"2016-09-26T13:42:51.564Z","lastUpdatedDate":"2017-01-06T05:11:51.682Z","links":[{"@type":"link","rel":"GET: Request Template","href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItems/502efc1b-d5ce-4ef9-99ee-d4e2a741747c/requests/template{?businessGroupId,requestedFor}"},{"@type":"link","rel":"POST: Submit Request","href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItems/502efc1b-d5ce-4ef9-99ee-d4e2a741747c/requests{?businessGroupId,requestedFor}"}],"iconId":"502efc1b-d5ce-4ef9-99ee-d4e2a741747c","catalogItemTypeRef":{"id":"com.vmware.csp.component.cafe.composition.blueprint","label":"Composite Blueprint"},"serviceRef":{"id":"baad0ad2-8b96-4347-b188-f534dad53a0d","label":"Infrastructure"},"outputResourceTypeRef":{"id":"composition.resource.type.deployment","label":"Deployment"}}],"metadata":{"size":20,"totalElements":44,"totalPages":3,"number":1,"offset":0}}`))

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItemViews?page=1&limit=44",
		httpmock.NewStringResponder(200, `{"links":[{"@type":"link","rel":"next","href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItemViews?page=2&limit=20"}],"content":[{"@type":"ConsumerEntitledCatalogItemView","entitledOrganizations":[{"tenantRef":"vsphere.local","tenantLabel":"vsphere.local","subtenantRef":"53619006-56bb-4788-9723-9eab79752cc1","subtenantLabel":"Content"}],"catalogItemId":"e5dd4fba-45ed-4943-b1fc-7f96239286be","name":"CentOS 6.3","description":"CentOS 6.3 IaaS Blueprint","isNoteworthy":false,"dateCreated":"2015-12-22T03:16:19.289Z","lastUpdatedDate":"2017-01-06T05:12:56.690Z","links":[{"@type":"link","rel":"GET: Request Template","href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItems/e5dd4fba-45ed-4943-b1fc-7f96239286be/requests/template{?businessGroupId,requestedFor}"},{"@type":"link","rel":"POST: Submit Request","href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItems/e5dd4fba-45ed-4943-b1fc-7f96239286be/requests{?businessGroupId,requestedFor}"}],"iconId":"e5dd4fba-45ed-4943-b1fc-7f96239286be","catalogItemTypeRef":{"id":"com.vmware.csp.component.cafe.composition.blueprint","label":"Composite Blueprint"},"serviceRef":{"id":"baad0ad2-8b96-4347-b188-f534dad53a0d","label":"Infrastructure"},"outputResourceTypeRef":{"id":"composition.resource.type.deployment","label":"Deployment"}},{"@type":"ConsumerEntitledCatalogItemView","entitledOrganizations":[{"tenantRef":"vsphere.local","tenantLabel":"vsphere.local","subtenantRef":"53619006-56bb-4788-9723-9eab79752cc1","subtenantLabel":"Content"}],"catalogItemId":"bd4cadc5-b76d-4e2c-b89c-e0ceedc987e0","name":"Wordpress for Windows","description":"Install Wordpress on Windows in IIS","isNoteworthy":false,"dateCreated":"2016-09-28T14:26:54.446Z","lastUpdatedDate":"2016-10-02T05:42:49.373Z","links":[{"@type":"link","rel":"GET: Request Template","href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItems/bd4cadc5-b76d-4e2c-b89c-e0ceedc987e0/requests/template{?businessGroupId,requestedFor}"},{"@type":"link","rel":"POST: Submit Request","href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItems/bd4cadc5-b76d-4e2c-b89c-e0ceedc987e0/requests{?businessGroupId,requestedFor}"}],"iconId":"bd4cadc5-b76d-4e2c-b89c-e0ceedc987e0","catalogItemTypeRef":{"id":"com.vmware.csp.component.cafe.composition.blueprint","label":"Composite Blueprint"},"serviceRef":{"id":"3ec368fb-8237-4430-a83f-39199e9aea6d","label":"Platform"},"outputResourceTypeRef":{"id":"composition.resource.type.deployment","label":"Deployment"}},{"@type":"ConsumerEntitledCatalogItemView","entitledOrganizations":[{"tenantRef":"vsphere.local","tenantLabel":"vsphere.local","subtenantRef":"53619006-56bb-4788-9723-9eab79752cc1","subtenantLabel":"Content"}],"catalogItemId":"7017b51e-2495-4f1e-9846-c94fc99c861f","name":"Windows 2012 R2 with IIS","description":"Windows 2012 R2 with IIS","isNoteworthy":false,"dateCreated":"2016-10-02T05:54:51.918Z","lastUpdatedDate":"2016-12-22T21:58:53.906Z","links":[{"@type":"link","rel":"GET: Request Template","href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItems/7017b51e-2495-4f1e-9846-c94fc99c861f/requests/template{?businessGroupId,requestedFor}"},{"@type":"link","rel":"POST: Submit Request","href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItems/7017b51e-2495-4f1e-9846-c94fc99c861f/requests{?businessGroupId,requestedFor}"}],"iconId":"7017b51e-2495-4f1e-9846-c94fc99c861f","catalogItemTypeRef":{"id":"com.vmware.csp.component.cafe.composition.blueprint","label":"Composite Blueprint"},"serviceRef":{"id":"baad0ad2-8b96-4347-b188-f534dad53a0d","label":"Infrastructure"},"outputResourceTypeRef":{"id":"composition.resource.type.deployment","label":"Deployment"}}],"metadata":{"size":20,"totalElements":44,"totalPages":3,"number":1,"offset":0}}`))

	//Fetch catalog by correct and exact name
	catalogID, err := client.readCatalogIDByName("CentOS 6.3 - IPAM EXT")

	if err != nil {
		t.Errorf("Error while fetching catalog %v", err)
	}

	if catalogID == nil || catalogID == "" {
		t.Errorf("Catalog ID is nil")
	}

	//Fetch catalog by correct and false name
	catalogID, err = client.readCatalogIDByName("Cent OS 6.3")

	if catalogID != nil && catalogID != "" {
		t.Errorf("Catalog ID is not nil")
	}

	if err != nil {
		t.Errorf("Error while fetching catalog %v", err)
	}

	//Fetch catalog by correct and incomplete name
	catalogID, err = client.readCatalogIDByName("CentOS")

	if catalogID != nil && catalogID != "" {
		t.Errorf("Catalog ID is not nil")
	}

	if err == nil {
		t.Errorf("Error should have been occurred while fetching the catalog with encomplete name")
	}
	httpmock.Reset()
	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItemViews",
		httpmock.NewStringResponder(404, `{"links":[{"@type":"link","rel":"next","href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItemViews?page=2&limit=20"}],"content":[{"@type":"ConsumerEntitledCatalogItemView","entitledOrganizations":[{"tenantRef":"vsphere.local","tenantLabel":"vsphere.local","subtenantRef":"53619006-56bb-4788-9723-9eab79752cc1","subtenantLabel":"Content"}],"catalogItemId":"502efc1b-d5ce-4ef9-99ee-d4e2a741747c","name":"CentOS 6.3 - IPAM EXT","description":"CentOS 6.3 IaaS Blueprint w/Infoblox IPAM","isNoteworthy":false,"dateCreated":"2016-09-26T13:42:51.564Z","lastUpdatedDate":"2017-01-06T05:11:51.682Z","links":[{"@type":"link","rel":"GET: Request Template","href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItems/502efc1b-d5ce-4ef9-99ee-d4e2a741747c/requests/template{?businessGroupId,requestedFor}"},{"@type":"link","rel":"POST: Submit Request","href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItems/502efc1b-d5ce-4ef9-99ee-d4e2a741747c/requests{?businessGroupId,requestedFor}"}],"iconId":"502efc1b-d5ce-4ef9-99ee-d4e2a741747c","catalogItemTypeRef":{"id":"com.vmware.csp.component.cafe.composition.blueprint","label":"Composite Blueprint"},"serviceRef":{"id":"baad0ad2-8b96-4347-b188-f534dad53a0d","label":"Infrastructure"},"outputResourceTypeRef":{"id":"composition.resource.type.deployment","label":"Deployment"}}],"metadata":{"size":20,"totalElements":44,"totalPages":3,"number":1,"offset":0}}`))

	//Fetch catalog by correct and incomplete name
	catalogID, err = client.readCatalogIDByName("CentOS")

	if catalogID != nil {
		t.Errorf("Catalog ID is not nil")
	}

	if err == nil {
		t.Errorf("Error should have been occurred while fetching the catalog with encomplete name")
	}

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItemViews",
		httpmock.NewStringResponder(404, `{"errors":[{"code":20116,"source":null,"message":"Unable to find the specified catalog item in the service catalog: e5dd4fba-45ed-4943-b1fc-07f96239286b.","systemMessage":"Unable to find the specified catalog item in the service catalog: e5dd4fba-45ed-4943-b1fc-07f96239286b.","moreInfoUrl":null}]}`))
	catalogID, err = client.readCatalogIDByName("CentOS")

	if err == nil {
		t.Errorf("Data fetched with wrong catalog ID")
	}

	if catalogID != nil {
		t.Errorf("Wrong catalog data got fetched")
	}
	httpmock.Reset()
	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItemViews",
		httpmock.NewStringResponder(200, `{"links":[{"@type":"link","rel":"next","href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItemViews?page=2&limit=20"}],"content":[{"@type":"ConsumerEntitledCatalogItemView","entitledOrganizations":[{"tenantRef":"vsphere.local","tenantLabel":"vsphere.local","subtenantRef":"53619006-56bb-4788-9723-9eab79752cc1","subtenantLabel":"Content"}],"catalogItemId":"502efc1b-d5ce-4ef9-99ee-d4e2a741747c","name":"CentOS 6.3 - IPAM EXT","description":"CentOS 6.3 IaaS Blueprint w/Infoblox IPAM","isNoteworthy":false,"dateCreated":"2016-09-26T13:42:51.564Z","lastUpdatedDate":"2017-01-06T05:11:51.682Z","links":[{"@type":"link","rel":"GET: Request Template","href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItems/502efc1b-d5ce-4ef9-99ee-d4e2a741747c/requests/template{?businessGroupId,requestedFor}"},{"@type":"link","rel":"POST: Submit Request","href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItems/502efc1b-d5ce-4ef9-99ee-d4e2a741747c/requests{?businessGroupId,requestedFor}"}],"iconId":"502efc1b-d5ce-4ef9-99ee-d4e2a741747c","catalogItemTypeRef":{"id":"com.vmware.csp.component.cafe.composition.blueprint","label":"Composite Blueprint"},"serviceRef":{"id":"baad0ad2-8b96-4347-b188-f534dad53a0d","label":"Infrastructure"},"outputResourceTypeRef":{"id":"composition.resource.type.deployment","label":"Deployment"}}],"metadata":{"size":20,"totalElements":44,"totalPages":3,"number":1,"offset":0}}`))

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItemViews?page=1&limit=44",
		httpmock.NewStringResponder(404, `{"errors":[{"code":20116,"source":null,"message":"Unable to find the specified catalog item in the service catalog: e5dd4fba-45ed-4943-b1fc-07f96239286b.","systemMessage":"Unable to find the specified catalog item in the service catalog: e5dd4fba-45ed-4943-b1fc-07f96239286b.","moreInfoUrl":null}]}`))

	catalogID, err = client.readCatalogIDByName("CentOS")

	if err == nil {
		t.Errorf("Data fetched with wrong catalog ID")
	}

	if catalogID != nil {
		t.Errorf("Wrong catalog data got fetched")
	}

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItemViews",
		httpmock.NewErrorResponder(errors.New(`{"errors":[{"code":50505,"source":null,"message":"System exception.","systemMessage":null,"moreInfoUrl":null}]}`)))
	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItemViews",
		httpmock.NewStringResponder(500, `{"errors":[{"code":20116,"source":null,"message":"Unable to find the specified catalog item in the service catalog: e5dd4fba-45ed-4943-b1fc-07f96239286b.","systemMessage":"Unable to find the specified catalog item in the service catalog: e5dd4fba-45ed-4943-b1fc-07f96239286b.","moreInfoUrl":null}]}`))
	catalogID, err = client.readCatalogIDByName("CentOS")
}

func TestFetchCatalogByID(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItems/e5dd4fba-45ed-4943-b1fc-7f96239286be",
		httpmock.NewStringResponder(200, `{"catalogItem":{"callbacks":null,"catalogItemTypeRef":{"id":"com.vmware.csp.component.cafe.composition.blueprint","label":"Composite Blueprint"},"dateCreated":"2015-12-22T03:16:19.289Z","description":"CentOS 6.3 IaaS Blueprint","forms":{"itemDetails":{"type":"external","formId":"composition.catalog.item.details"},"catalogRequestInfoHidden":true,"requestFormScale":"BIG","requestSubmission":{"type":"extension","extensionId":"com.vmware.vcac.core.design.blueprints.requestForm","extensionPointId":null},"requestDetails":{"type":"extension","extensionId":"com.vmware.vcac.core.design.blueprints.requestDetailsForm","extensionPointId":null},"requestPreApproval":null,"requestPostApproval":null},"iconId":"e5dd4fba-45ed-4943-b1fc-7f96239286be","id":"e5dd4fba-45ed-4943-b1fc-7f96239286be","isNoteworthy":false,"lastUpdatedDate":"2017-01-06T05:12:56.690Z","name":"CentOS 6.3","organization":{"tenantRef":"vsphere.local","tenantLabel":"vsphere.local","subtenantRef":null,"subtenantLabel":null},"outputResourceTypeRef":{"id":"composition.resource.type.deployment","label":"Deployment"},"providerBinding":{"bindingId":"vsphere.local!::!CentOS63","providerRef":{"id":"2fbaabc5-3a48-488a-9f2a-a42616345445","label":"Blueprint Service"}},"serviceRef":{"id":"baad0ad2-8b96-4347-b188-f534dad53a0d","label":"Infrastructure"},"status":"PUBLISHED","statusName":"Published","quota":0,"version":4,"requestable":true},"entitledOrganizations":[{"tenantRef":"vsphere.local","tenantLabel":"vsphere.local","subtenantRef":"53619006-56bb-4788-9723-9eab79752cc1","subtenantLabel":"Content"}]}`))

	catalogName, err := client.readCatalogNameByID("e5dd4fba-45ed-4943-b1fc-7f96239286be")

	if err != nil {
		t.Errorf("Error while fetching catalog %v", err)
	}

	if catalogName == nil {
		t.Errorf("Catalog Name is is nil")
	}

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItems/e5dd4fba-45ed-4943-b1fc-7f96239286be",
		httpmock.NewStringResponder(404, `{"errors":[{"code":20116,"source":null,"message":"Unable to find the specified catalog item in the service catalog: e5dd4fba-45ed-4943-b1fc-07f96239286b.","systemMessage":"Unable to find the specified catalog item in the service catalog: e5dd4fba-45ed-4943-b1fc-07f96239286b.","moreInfoUrl":null}]}`))

	catalogName, err = client.readCatalogNameByID("e5dd4fba-45ed-4943-b1fc-7f96239286b1")

	if err == nil {
		t.Errorf("Data fetched with wrong catalog ID")
	}

	if catalogName != nil {
		t.Errorf("Wrong catalog data got fetched")
	}

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/entitledCatalogItems/e5dd4fba-45ed-4943-b1fc-7f96239286b1",
		httpmock.NewStringResponder(404, `{"errors":[{"code":20116,"source":null,"message":"Unable to find the specified catalog item in the service catalog: e5dd4fba-45ed-4943-b1fc-07f96239286b.","systemMessage":"Unable to find the specified catalog item in the service catalog: e5dd4fba-45ed-4943-b1fc-07f96239286b.","moreInfoUrl":null}]}`))
	catalogName, err = client.readCatalogNameByID("e5dd4fba-45ed-4943-b1fc-7f96239286b1")

	if err == nil {
		t.Errorf("Data fetched with wrong catalog ID")
	}

	if catalogName != nil {
		t.Errorf("Wrong catalog data got fetched")
	}
}
