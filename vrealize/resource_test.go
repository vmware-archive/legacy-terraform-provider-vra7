package vrealize

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/terraform-provider-vra7/utils"
	"gopkg.in/jarcoal/httpmock.v1"
)

var client APIClient

func TestNewClient(t *testing.T) {
	username := "admin@myvra.local"
	password := "pass!@#"
	tenant := "vshpere.local"
	baseURL := "http://localhost/"

	client := NewClient(
		username,
		password,
		tenant,
		baseURL,
		true,
	)

	if client.Username != username {
		t.Errorf("Expected username %v, got %v ", username, client.Username)
	}

	if client.Password != password {
		t.Errorf("Expected password %v, got %v ", password, client.Password)
	}

	if client.Tenant != tenant {
		t.Errorf("Expected tenant %v, got %v ", tenant, client.Tenant)
	}

	if client.BaseURL != baseURL {
		t.Errorf("Expected BaseUrl %v, got %v ", baseURL, client.BaseURL)
	}
}

func TestClient_Authenticate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost/identity/api/tokens",
		httpmock.NewStringResponder(200, `{
		  "expires": "2017-07-25T15:18:49.000Z",
		  "id": "MTUwMDk2NzEyOTEyOTplYTliNTA3YTg4MjZmZjU1YTIwZjp0ZW5hbnQ6dnNwaGVyZS5sb2NhbHVzZXJuYW1lOmphc29uQGNvcnAubG9jYWxleHBpcmF0aW9uOjE1MDA5OTU5MjkwMDA6ZjE1OTQyM2Y1NjQ2YzgyZjY4Yjg1NGFjMGNkNWVlMTNkNDhlZTljNjY3ZTg4MzA1MDViMTU4Y2U3MzBkYjQ5NmQ5MmZhZWM1MWYzYTg1ZWM4ZDhkYmFhMzY3YTlmNDExZmM2MTRmNjk5MGQ1YjRmZjBhYjgxMWM0OGQ3ZGVmNmY=",
		  "tenant": "vsphere.local"
		}`))

	err := client.Authenticate()

	if len(client.BearerToken) == 0 {
		t.Error("Fail to set BearerToken.")
	}

	httpmock.RegisterResponder("POST", "http://localhost/identity/api/tokens",
		httpmock.NewErrorResponder(errors.New(`{
		  "errors": [
			{
			  "code": 90135,
			  "source": null,
			  "message": "Unable to authenticate user jason@corp.local1 in tenant vsphere.local.",
			  "systemMessage": "90135-Unable to authenticate user jason@corp.local1 in tenant vsphere.local.",
			  "moreInfoUrl": null
			}
		  ]
		}`)))

	err = client.Authenticate()

	if err == nil {
		t.Errorf("Authentication should fail")
	}
}

func TestAPIClient_GetCatalogItemRequestTemplate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/"+
		"api/consumer/entitledCatalogItems/e5dd4fba-45ed-4943-b1fc-7f96239286be/requests/template",
		httpmock.NewStringResponder(200, `{"type":"com.vmware.vcac.catalog.domain.request.CatalogItemProvisioningRequest","catalogItemId":"e5dd4fba-45ed-4943-b1fc-7f96239286be","requestedFor":"jason@corp.local","businessGroupId":"53619006-56bb-4788-9723-9eab79752cc1","description":null,"reasons":null,"data":{"CentOS_6.3":{"componentTypeId":"com.vmware.csp.component.cafe.composition","componentId":null,"classId":"Blueprint.Component.Declaration","typeFilter":"CentOS63*CentOS_6.3","data":{"_allocation":{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.Allocation","typeFilter":null,"data":{"machines":[{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.Allocation.Machine","typeFilter":null,"data":{"machine_id":"","nics":[{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.Nic","typeFilter":null,"data":{"address":"","assignment_type":"Static","external_address":"","id":null,"load_balancing":null,"network":null,"network_profile":null}}]}}]}},"_cluster":1,"_hasChildren":false,"cpu":1,"datacenter_location":null,"description":"Basic IaaS CentOS Machine","disks":[{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.MachineDisk","typeFilter":null,"data":{"capacity":3,"custom_properties":null,"id":1450725224066,"initial_location":"","is_clone":true,"label":"Hard disk 1","storage_reservation_policy":"","userCreated":false,"volumeId":0}}],"guest_customization_specification":"CentOS","max_network_adapters":-1,"max_per_user":0,"max_volumes":60,"memory":512,"nics":null,"os_arch":"x86_64","os_distribution":null,"os_type":"Linux","os_version":null,"property_groups":null,"reservation_policy":null,"security_groups":[],"security_tags":[],"storage":3}},"_archiveDays":5,"_leaseDays":null,"_number_of_instances":1,"corp192168110024":{"componentTypeId":"com.vmware.csp.component.cafe.composition","componentId":null,"classId":"Blueprint.Component.Declaration","typeFilter":"CentOS63*corp192168110024","data":{"_hasChildren":false}}}}`))

	template, err := client.GetCatalogItemRequestTemplate("e5dd4fba-45ed-4943-b1fc-7f96239286be")

	if err != nil {
		t.Errorf("Fail to get catalog Item template %v.", err)
	}

	if len(template.CatalogItemID) == 0 {
		t.Errorf("Catalog Item id is empty.")
	}
	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/"+
		"api/consumer/entitledCatalogItems/e5dd4fba-45ed-4943-b1fc-7f96239286be/requests/template",
		httpmock.NewErrorResponder(errors.New(`"errors":[{"code":20116,"source":null,"message":"Unable to find the specified catalog item in the service catalog: ae368563-867e-40c1-a09e-0aeec55c9e81.","systemMessage":"Unable to find the specified catalog item in the service catalog: ae368563-867e-40c1-a09e-0aeec55c9e81.","moreInfoUrl":null}]}`)))

	template, err = client.GetCatalogItemRequestTemplate("e5dd4fba-45ed-4943-b1fc-7f96239286be")

	if err == nil {
		t.Errorf("Fail to generate exception")
	}

	if template != nil {
		t.Errorf("Catalog Item id is not empty.")
	}
}

func TestAPIClient_RequestCatalogItem(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/"+
		"api/consumer/entitledCatalogItems/e5dd4fba-45ed-4943-b1fc-7f96239286be/requests/template",
		httpmock.NewStringResponder(200, `{"type":"com.vmware.vcac.catalog.domain.request.CatalogItemProvisioningRequest","catalogItemId":"e5dd4fba-45ed-4943-b1fc-7f96239286be","requestedFor":"jason@corp.local","businessGroupId":"53619006-56bb-4788-9723-9eab79752cc1","description":null,"reasons":null,"data":{"CentOS_6.3":{"componentTypeId":"com.vmware.csp.component.cafe.composition","componentId":null,"classId":"Blueprint.Component.Declaration","typeFilter":"CentOS63*CentOS_6.3","data":{"_allocation":{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.Allocation","typeFilter":null,"data":{"machines":[{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.Allocation.Machine","typeFilter":null,"data":{"machine_id":"","nics":[{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.Nic","typeFilter":null,"data":{"address":"","assignment_type":"Static","external_address":"","id":null,"load_balancing":null,"network":null,"network_profile":null}}]}}]}},"_cluster":1,"_hasChildren":false,"cpu":1,"datacenter_location":null,"description":"Basic IaaS CentOS Machine","disks":[{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.MachineDisk","typeFilter":null,"data":{"capacity":3,"custom_properties":null,"id":1450725224066,"initial_location":"","is_clone":true,"label":"Hard disk 1","storage_reservation_policy":"","userCreated":false,"volumeId":0}}],"guest_customization_specification":"CentOS","max_network_adapters":-1,"max_per_user":0,"max_volumes":60,"memory":512,"nics":null,"os_arch":"x86_64","os_distribution":null,"os_type":"Linux","os_version":null,"property_groups":null,"reservation_policy":null,"security_groups":[],"security_tags":[],"storage":3}},"_archiveDays":5,"_leaseDays":null,"_number_of_instances":1,"corp192168110024":{"componentTypeId":"com.vmware.csp.component.cafe.composition","componentId":null,"classId":"Blueprint.Component.Declaration","typeFilter":"CentOS63*corp192168110024","data":{"_hasChildren":false}}}}`))

	httpmock.RegisterResponder("POST", "http://localhost/catalog-service/"+
		"api/consumer/entitledCatalogItems/e5dd4fba-45ed-4943-b1fc-7f96239286be/requests",
		httpmock.NewStringResponder(201, `{"@type":"CatalogItemRequest","id":"b2907df7-6c36-4e30-9c62-a21f293b067a","iconId":"composition.blueprint.png","version":0,"requestNumber":null,"state":"SUBMITTED","description":null,"reasons":null,"requestedFor":"jason@corp.local","requestedBy":"jason@corp.local","organization":{"tenantRef":"vsphere.local","tenantLabel":null,"subtenantRef":"29a02ed9-7e63-4c77-8a15-c930afb0e3d8","subtenantLabel":null},"requestorEntitlementId":"e0d6ce92-6e23-4f75-a787-4564699b2895","preApprovalId":null,"postApprovalId":null,"dateCreated":"2017-08-10T13:38:25.395Z","lastUpdated":"2017-08-10T13:38:25.395Z","dateSubmitted":"2017-08-10T13:38:25.395Z","dateApproved":null,"dateCompleted":null,"quote":{"leasePeriod":null,"leaseRate":null,"totalLeaseCost":null},"requestCompletion":null,"requestData":{"entries":[{"key":"MySQL_1","value":{"type":"complex","componentTypeId":"com.vmware.csp.component.cafe.composition","componentId":null,"classId":"Blueprint.Component.Declaration","typeFilter":"checkcloudclient*MySQL_1","values":{"entries":[{"key":"_hasChildren","value":{"type":"boolean","value":false}},{"key":"dbpassword","value":{"type":"secureString","value":"catalog~+gzbqycW+GiAqOREkOs7+mW9D4Og83AKc4FE46i2Z6Y="}}]}}},{"key":"Apache_Load_Balancer_1","value":{"type":"complex","componentTypeId":"com.vmware.csp.component.cafe.composition","componentId":null,"classId":"Blueprint.Component.Declaration","typeFilter":"checkcloudclient*Apache_Load_Balancer_1","values":{"entries":[{"key":"http_node_ips","value":{"type":"multiple","elementTypeId":"STRING","items":[{"type":"string","value":"None"}]}},{"key":"_hasChildren","value":{"type":"boolean","value":false}},{"key":"http_proxy_port","value":{"type":"string","value":"8081"}},{"key":"tomcat_context","value":null},{"key":"JAVA_HOME","value":{"type":"string","value":"/opt/vmware-jre"}},{"key":"appsrv_routes","value":{"type":"multiple","elementTypeId":"STRING","items":[{"type":"string","value":"None"}]}},{"key":"use_ajp","value":{"type":"string","value":"NO"}},{"key":"http_node_port","value":{"type":"multiple","elementTypeId":"STRING","items":[{"type":"string","value":"8080"}]}},{"key":"http_port","value":{"type":"string","value":"80"}},{"key":"autogen_sticky_cookie","value":{"type":"string","value":"NO"}}]}}},{"key":"corp192168110024","value":{"type":"complex","componentTypeId":"com.vmware.csp.component.cafe.composition","componentId":null,"classId":"Blueprint.Component.Declaration","typeFilter":"checkcloudclient*corp192168110024","values":{"entries":[{"key":"_hasChildren","value":{"type":"boolean","value":false}}]}}},{"key":"providerId","value":{"type":"string","value":"2fbaabc5-3a48-488a-9f2a-a42616345445"}},{"key":"subtenantId","value":{"type":"string","value":"29a02ed9-7e63-4c77-8a15-c930afb0e3d8"}},{"key":"vSphere__vCenter__Machine_2","value":{"type":"complex","componentTypeId":"com.vmware.csp.component.cafe.composition","componentId":null,"classId":"Blueprint.Component.Declaration","typeFilter":"checkcloudclient*vSphere__vCenter__Machine_2","values":{"entries":[{"key":"snapshot_name","value":null},{"key":"source_machine","value":null},{"key":"memory","value":{"type":"integer","value":512}},{"key":"disks","value":{"type":"multiple","elementTypeId":"COMPLEX","items":[{"type":"complex","componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.MachineDisk","typeFilter":null,"values":{"entries":[{"key":"is_clone","value":{"type":"boolean","value":false}},{"key":"initial_location","value":{"type":"string","value":""}},{"key":"volumeId","value":{"type":"string","value":"0"}},{"key":"id","value":{"type":"integer","value":1502347498478}},{"key":"label","value":{"type":"string","value":""}},{"key":"userCreated","value":{"type":"boolean","value":true}},{"key":"storage_reservation_policy","value":{"type":"string","value":""}},{"key":"capacity","value":{"type":"integer","value":1}}]}}]}},{"key":"description","value":null},{"key":"storage","value":{"type":"integer","value":1}},{"key":"source_machine_name","value":null},{"key":"guest_customization_specification","value":null},{"key":"_hasChildren","value":{"type":"boolean","value":true}},{"key":"os_distribution","value":null},{"key":"reservation_policy","value":null},{"key":"max_network_adapters","value":{"type":"integer","value":-1}},{"key":"machine_prefix","value":null},{"key":"max_per_user","value":{"type":"integer","value":0}},{"key":"nics","value":null},{"key":"source_machine_vmsnapshot","value":null},{"key":"_allocation","value":{"type":"complex","componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.Allocation","typeFilter":null,"values":{"entries":[{"key":"machines","value":null}]}}},{"key":"display_location","value":{"type":"boolean","value":false}},{"key":"os_version","value":null},{"key":"os_arch","value":{"type":"string","value":"x86_64"}},{"key":"cpu","value":{"type":"integer","value":1}},{"key":"datacenter_location","value":null},{"key":"property_groups","value":null},{"key":"_cluster","value":{"type":"integer","value":1}},{"key":"security_groups","value":{"type":"multiple","elementTypeId":"ENTITY_REFERENCE","items":[]}},{"key":"max_volumes","value":{"type":"integer","value":60}},{"key":"os_type","value":{"type":"string","value":"Linux"}},{"key":"source_machine_external_snapshot","value":null},{"key":"security_tags","value":{"type":"multiple","elementTypeId":"ENTITY_REFERENCE","items":[]}}]}}},{"key":"_leaseDays","value":null},{"key":"providerBindingId","value":{"type":"string","value":"checkcloudclient"}},{"key":"_number_of_instances","value":{"type":"integer","value":1}},{"key":"vSphere__vCenter__Machine_1","value":{"type":"complex","componentTypeId":"com.vmware.csp.component.cafe.composition","componentId":null,"classId":"Blueprint.Component.Declaration","typeFilter":"checkcloudclient*vSphere__vCenter__Machine_1","values":{"entries":[{"key":"snapshot_name","value":null},{"key":"source_machine","value":null},{"key":"memory","value":{"type":"integer","value":512}},{"key":"disks","value":{"type":"multiple","elementTypeId":"COMPLEX","items":[{"type":"complex","componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.MachineDisk","typeFilter":null,"values":{"entries":[{"key":"is_clone","value":{"type":"boolean","value":false}},{"key":"initial_location","value":{"type":"string","value":"hd-1"}},{"key":"volumeId","value":{"type":"string","value":"0"}},{"key":"id","value":{"type":"integer","value":1502345335122}},{"key":"label","value":{"type":"string","value":""}},{"key":"userCreated","value":{"type":"boolean","value":true}},{"key":"storage_reservation_policy","value":{"type":"string","value":""}},{"key":"capacity","value":{"type":"integer","value":3}}]}}]}},{"key":"description","value":null},{"key":"storage","value":{"type":"integer","value":3}},{"key":"source_machine_name","value":null},{"key":"guest_customization_specification","value":null},{"key":"_hasChildren","value":{"type":"boolean","value":true}},{"key":"os_distribution","value":null},{"key":"reservation_policy","value":null},{"key":"max_network_adapters","value":{"type":"integer","value":-1}},{"key":"machine_prefix","value":null},{"key":"max_per_user","value":{"type":"integer","value":0}},{"key":"nics","value":null},{"key":"source_machine_vmsnapshot","value":null},{"key":"_allocation","value":{"type":"complex","componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.Allocation","typeFilter":null,"values":{"entries":[{"key":"machines","value":null}]}}},{"key":"display_location","value":{"type":"boolean","value":false}},{"key":"os_version","value":null},{"key":"os_arch","value":{"type":"string","value":"x86_64"}},{"key":"cpu","value":{"type":"integer","value":1}},{"key":"datacenter_location","value":null},{"key":"property_groups","value":null},{"key":"_cluster","value":{"type":"integer","value":1}},{"key":"security_groups","value":{"type":"multiple","elementTypeId":"ENTITY_REFERENCE","items":[]}},{"key":"max_volumes","value":{"type":"integer","value":60}},{"key":"os_type","value":{"type":"string","value":"Linux"}},{"key":"source_machine_external_snapshot","value":null},{"key":"security_tags","value":{"type":"multiple","elementTypeId":"ENTITY_REFERENCE","items":[]}}]}}}]},"retriesRemaining":3,"requestedItemName":"myCompositeBlueprint","requestedItemDescription":"","components":null,"stateName":null,"catalogItemRef":{"id":"a3647254-3c50-4fe6-a630-69ae28bf3c81","label":"myCompositeBlueprint"},"catalogItemProviderBinding":{"bindingId":"vsphere.local!::!checkcloudclient","providerRef":{"id":"2fbaabc5-3a48-488a-9f2a-a42616345445","label":"Blueprint Service"}},"waitingStatus":"NOT_WAITING","executionStatus":"STARTED","approvalStatus":"PENDING","phase":"PENDING_PRE_APPROVAL"}`))

	template, err := client.GetCatalogItemRequestTemplate("e5dd4fba-45ed-4943-b1fc-7f96239286be")
	if err != nil {
		t.Errorf("Failed to get catalog item template %v.", err)
	}
	if len(template.CatalogItemID) == 0 {
		t.Errorf("Catalog Item Id is empty.")
	}

	catalogRequest, errorRequestCatalogItem := client.RequestCatalogItem(template)

	if errorRequestCatalogItem != nil {
		t.Errorf("Failed to request the catalog item %v.", errorRequestCatalogItem)
	}

	if len(catalogRequest.ID) == 0 {
		t.Errorf("Failed to request catalog item.")
	}

	httpmock.RegisterResponder("POST", "http://localhost/catalog-service/"+
		"api/consumer/entitledCatalogItems/e5dd4fba-45ed-4943-b1fc-7f96239286be/requests",
		httpmock.NewErrorResponder(errors.New(`{"errors":[{"code":20116,"source":null,"message":"Unable to find the specified catalog item in the service catalog: a3647254-3c50-4fe6-a636-9ae28bf3c811.","systemMessage":"Unable to find the specified catalog item in the service catalog: a3647254-3c50-4fe6-a636-9ae28bf3c811.","moreInfoUrl":null}]}`)))

	catalogRequest, errorRequestCatalogItem = client.RequestCatalogItem(template)

	if errorRequestCatalogItem == nil {
		t.Errorf("Failed to generate exception.")
	}

	if catalogRequest != nil {
		t.Errorf("Catalog item request initiated successfully.")
	}
}
func TestAPIClient_GetResourceViews(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/"+
		"api/consumer/requests/937099db-5174-4862-99a3-9c2666bfca28/resourceViews",
		httpmock.NewStringResponder(200, `{"links":[],"content":[{"@type":"CatalogResourceView","resourceId":"b313acd6-0738-439c-b601-e3ebf9ebb49b","iconId":"502efc1b-d5ce-4ef9-99ee-d4e2a741747c","name":"CentOS 6.3 - IPAM EXT-95563173","description":"","status":null,"catalogItemId":"502efc1b-d5ce-4ef9-99ee-d4e2a741747c","catalogItemLabel":"CentOS 6.3 - IPAM EXT","requestId":"dcb12203-93f4-4873-a7d5-1757f3696141","requestState":"SUCCESSFUL","resourceType":"composition.resource.type.deployment","owners":["Jason Cloud Admin"],"businessGroupId":"53619006-56bb-4788-9723-9eab79752cc1","tenantId":"vsphere.local","dateCreated":"2017-07-17T13:26:42.102Z","lastUpdated":"2017-07-17T13:33:25.521Z","lease":{"start":"2017-07-17T13:26:42.079Z","end":null},"costs":null,"costToDate":null,"totalCost":null,"parentResourceId":null,"hasChildren":true,"data":{},"links":[{"@type":"link","rel":"GET: Catalog Item","href":"http://localhost/catalog-service/api/consumer/entitledCatalogItemViews/502efc1b-d5ce-4ef9-99ee-d4e2a741747c"},{"@type":"link","rel":"GET: Request","href":"http://localhost/catalog-service/api/consumer/requests/dcb12203-93f4-4873-a7d5-1757f3696141"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changelease.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/561be422-ece6-4316-8acb-a8f3dbb8ed0c/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changelease.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/561be422-ece6-4316-8acb-a8f3dbb8ed0c/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changeowner.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/59249166-e427-4082-a3dc-eb7223bb2de1/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changeowner.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/59249166-e427-4082-a3dc-eb7223bb2de1/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.destroy.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/3da0ca14-e7e2-4d7b-89cb-c6db57440d72/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.destroy.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/3da0ca14-e7e2-4d7b-89cb-c6db57440d72/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.archive.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/9725d56e-461a-471a-be00-b1856681c6d0/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.archive.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/9725d56e-461a-471a-be00-b1856681c6d0/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scalein.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/85e090f9-9529-4101-9691-6bab1b0a1f77/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scalein.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/85e090f9-9529-4101-9691-6bab1b0a1f77/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scaleout.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/ab5795f5-32ad-4f6c-8598-1d3a7d190caa/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scaleout.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/ab5795f5-32ad-4f6c-8598-1d3a7d190caa/requests"},{"@type":"link","rel":"GET: Child Resources","href":"http://localhost/catalog-service/api/consumer/resourceViews?managedOnly=false&withExtendedData=true&withOperations=true&%24filter=parentResource%20eq%20%27b313acd6-0738-439c-b601-e3ebf9ebb49b%27"}]},{"@type":"CatalogResourceView","resourceId":"51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5","iconId":"Infrastructure.CatalogItem.Machine.Virtual.vSphere","name":"Content0061","description":"Basic IaaS CentOS Machine","status":"Missing","catalogItemId":null,"catalogItemLabel":null,"requestId":"dcb12203-93f4-4873-a7d5-1757f3696141","requestState":"SUCCESSFUL","resourceType":"Infrastructure.Virtual","owners":["Jason Cloud Admin"],"businessGroupId":"53619006-56bb-4788-9723-9eab79752cc1","tenantId":"vsphere.local","dateCreated":"2017-07-17T13:33:16.686Z","lastUpdated":"2017-07-17T13:33:25.521Z","lease":{"start":"2017-07-17T13:26:42.079Z","end":null},"costs":null,"costToDate":null,"totalCost":null,"parentResourceId":"b313acd6-0738-439c-b601-e3ebf9ebb49b","hasChildren":false,"data":{"Component":"CentOS_6.3","DISK_VOLUMES":[{"componentTypeId":"com.vmware.csp.component.iaas.proxy.provider","componentId":null,"classId":"dynamicops.api.model.DiskInputModel","typeFilter":null,"data":{"DISK_CAPACITY":3,"DISK_INPUT_ID":"DISK_INPUT_ID1","DISK_LABEL":"Hard disk 1"}}],"Destroy":true,"EXTERNAL_REFERENCE_ID":"vm-773","IS_COMPONENT_MACHINE":false,"MachineBlueprintName":"CentOS 6.3 - IPAM EXT","MachineCPU":1,"MachineDailyCost":0,"MachineDestructionDate":null,"MachineExpirationDate":null,"MachineGroupName":"Content","MachineGuestOperatingSystem":"CentOS 4/5/6/7 (64-bit)","MachineInterfaceDisplayName":"vSphere (vCenter)","MachineInterfaceType":"vSphere","MachineMemory":512,"MachineName":"Content0061","MachineReservationName":"IPAM Sandbox","MachineStorage":3,"MachineType":"Virtual","NETWORK_LIST":[{"componentTypeId":"com.vmware.csp.component.iaas.proxy.provider","componentId":null,"classId":"dynamicops.api.model.NetworkViewModel","typeFilter":null,"data":{"NETWORK_ADDRESS":"192.168.110.150","NETWORK_MAC_ADDRESS":"00:50:56:ae:31:bd","NETWORK_NAME":"VM Network","NETWORK_NETWORK_NAME":"ipamext1921681100","NETWORK_PROFILE":"ipam-ext-192.168.110.0"}}],"SNAPSHOT_LIST":[],"Unregister":true,"VirtualMachine.Admin.UUID":"502e9fb3-6f0d-0b1e-f90f-a769fd406620","endpointExternalReferenceId":"d322b019-58d4-4d6f-9f8b-d28695a716c0","ip_address":"192.168.110.150","machineId":"4fc33663-992d-49f8-af17-df7ce4831aa0"},"links":[{"@type":"link","rel":"GET: Request","href":"http://localhost/catalog-service/api/consumer/requests/dcb12203-93f4-4873-a7d5-1757f3696141"},{"@type":"link","rel":"GET: Parent Resource","href":"http://localhost/catalog-service/api/consumer/resourceViews/b313acd6-0738-439c-b601-e3ebf9ebb49b"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.virtual.Destroy}","href":"http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/654b4c71-e84f-40c7-9439-fd409fea7323/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.virtual.Destroy}","href":"http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/654b4c71-e84f-40c7-9439-fd409fea7323/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Unregister}","href":"http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/f3ae9408-885a-4a3a-9200-43366f2aa163/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Unregister}","href":"http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/f3ae9408-885a-4a3a-9200-43366f2aa163/requests"}]},{"@type":"CatalogResourceView","resourceId":"169b596f-e4c0-4b25-ba44-18cb19c0fd65","iconId":"existing_network","name":"ipamext1921681100","description":"Infoblox External Network","status":null,"catalogItemId":null,"catalogItemLabel":null,"requestId":"dcb12203-93f4-4873-a7d5-1757f3696141","requestState":"SUCCESSFUL","resourceType":"Infrastructure.Network.Network.Existing","owners":["Jason Cloud Admin"],"businessGroupId":"53619006-56bb-4788-9723-9eab79752cc1","tenantId":"vsphere.local","dateCreated":"2017-07-17T13:27:17.526Z","lastUpdated":"2017-07-17T13:33:25.521Z","lease":{"start":"2017-07-17T13:26:42.079Z","end":null},"costs":null,"costToDate":null,"totalCost":null,"parentResourceId":"b313acd6-0738-439c-b601-e3ebf9ebb49b","hasChildren":false,"data":{"Description":"Infoblox External Network","IPAMEndpointId":"1c2b6237-540a-43c3-8c06-b37a1d274b44","IPAMEndpointName":"Infoblox - nios01a","Name":"ipamext1921681100","_archiveDays":5,"_hasChildren":false,"_leaseDays":null,"_number_of_instances":1,"dns":{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Network.Network.DnsWins","typeFilter":null,"data":{"alternate_wins":null,"dns_search_suffix":null,"dns_suffix":null,"preferred_wins":null,"primary_dns":null,"secondary_dns":null}},"gateway":null,"ip_ranges":[{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Network.Network.IpRanges","typeFilter":null,"data":{"description":"","end_ip":"","externalId":"network/default-vra/192.168.110.0/24","id":"b078d23a-1c3d-4458-ab57-e352c80e6d55","name":"192.168.110.0/24","start_ip":""}}],"network_profile":"ipam-ext-192.168.110.0","providerBindingId":"CentOS63Infoblox","providerId":"2fbaabc5-3a48-488a-9f2a-a42616345445","subnet_mask":"255.255.255.0","subtenantId":"53619006-56bb-4788-9723-9eab79752cc1"},"links":[{"@type":"link","rel":"GET: Request","href":"http://localhost/catalog-service/api/consumer/requests/dcb12203-93f4-4873-a7d5-1757f3696141"},{"@type":"link","rel":"GET: Parent Resource","href":"http://localhost/catalog-service/api/consumer/resourceViews/b313acd6-0738-439c-b601-e3ebf9ebb49b"}]}],"metadata":{"size":20,"totalElements":3,"totalPages":1,"number":1,"offset":0}}`))

	template, err := client.GetDeploymentState("937099db-5174-4862-99a3-9c2666bfca28")
	if err != nil {
		t.Errorf("Fail to get resource views %v.", err)
	}
	if len(template.Content) == 0 {
		t.Errorf("No resources provisioned.")
	}

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/"+
		"api/consumer/requests/937099db-5174-4862-99a3-9c2666bfca28/resourceViews",
		httpmock.NewErrorResponder(errors.New(`{"errors":[{"code":20111,"source":null,"message":"Unable to find the specified request in the service catalog: dcb12203-93f4-4873-a7d5-757f36961411.","systemMessage":"Unable to find the specified request in the service catalog: dcb12203-93f4-4873-a7d5-757f36961411.","moreInfoUrl":null}]}`)))

	template, err = client.GetDeploymentState("937099db-5174-4862-99a3-9c2666bfca28")
	if err == nil {
		t.Errorf("Succeed to get resource views %v.", err)
	}
	if template != nil {
		t.Errorf("Resources provisioned.")
	}

}

func TestAPIClient_GetDestroyActionTemplate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/"+
		"api/consumer/requests/937099db-5174-4862-99a3-9c2666bfca28/resourceViews",
		httpmock.NewStringResponder(200, `{"links":[],"content":[{"@type":"CatalogResourceView","resourceId":"b313acd6-0738-439c-b601-e3ebf9ebb49b","iconId":"502efc1b-d5ce-4ef9-99ee-d4e2a741747c","name":"CentOS 6.3 - IPAM EXT-95563173","description":"","status":null,"catalogItemId":"502efc1b-d5ce-4ef9-99ee-d4e2a741747c","catalogItemLabel":"CentOS 6.3 - IPAM EXT","requestId":"dcb12203-93f4-4873-a7d5-1757f3696141","requestState":"SUCCESSFUL","resourceType":"composition.resource.type.deployment","owners":["Jason Cloud Admin"],"businessGroupId":"53619006-56bb-4788-9723-9eab79752cc1","tenantId":"vsphere.local","dateCreated":"2017-07-17T13:26:42.102Z","lastUpdated":"2017-07-17T13:33:25.521Z","lease":{"start":"2017-07-17T13:26:42.079Z","end":null},"costs":null,"costToDate":null,"totalCost":null,"parentResourceId":null,"hasChildren":true,"data":{},"links":[{"@type":"link","rel":"GET: Catalog Item","href":"http://localhost/catalog-service/api/consumer/entitledCatalogItemViews/502efc1b-d5ce-4ef9-99ee-d4e2a741747c"},{"@type":"link","rel":"GET: Request","href":"http://localhost/catalog-service/api/consumer/requests/dcb12203-93f4-4873-a7d5-1757f3696141"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changelease.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/561be422-ece6-4316-8acb-a8f3dbb8ed0c/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changelease.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/561be422-ece6-4316-8acb-a8f3dbb8ed0c/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changeowner.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/59249166-e427-4082-a3dc-eb7223bb2de1/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changeowner.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/59249166-e427-4082-a3dc-eb7223bb2de1/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.destroy.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/3da0ca14-e7e2-4d7b-89cb-c6db57440d72/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.destroy.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/3da0ca14-e7e2-4d7b-89cb-c6db57440d72/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.archive.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/9725d56e-461a-471a-be00-b1856681c6d0/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.archive.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/9725d56e-461a-471a-be00-b1856681c6d0/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scalein.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/85e090f9-9529-4101-9691-6bab1b0a1f77/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scalein.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/85e090f9-9529-4101-9691-6bab1b0a1f77/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scaleout.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/ab5795f5-32ad-4f6c-8598-1d3a7d190caa/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scaleout.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/ab5795f5-32ad-4f6c-8598-1d3a7d190caa/requests"},{"@type":"link","rel":"GET: Child Resources","href":"http://localhost/catalog-service/api/consumer/resourceViews?managedOnly=false&withExtendedData=true&withOperations=true&%24filter=parentResource%20eq%20%27b313acd6-0738-439c-b601-e3ebf9ebb49b%27"}]},{"@type":"CatalogResourceView","resourceId":"51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5","iconId":"Infrastructure.CatalogItem.Machine.Virtual.vSphere","name":"Content0061","description":"Basic IaaS CentOS Machine","status":"Missing","catalogItemId":null,"catalogItemLabel":null,"requestId":"dcb12203-93f4-4873-a7d5-1757f3696141","requestState":"SUCCESSFUL","resourceType":"Infrastructure.Virtual","owners":["Jason Cloud Admin"],"businessGroupId":"53619006-56bb-4788-9723-9eab79752cc1","tenantId":"vsphere.local","dateCreated":"2017-07-17T13:33:16.686Z","lastUpdated":"2017-07-17T13:33:25.521Z","lease":{"start":"2017-07-17T13:26:42.079Z","end":null},"costs":null,"costToDate":null,"totalCost":null,"parentResourceId":"b313acd6-0738-439c-b601-e3ebf9ebb49b","hasChildren":false,"data":{"Component":"CentOS_6.3","DISK_VOLUMES":[{"componentTypeId":"com.vmware.csp.component.iaas.proxy.provider","componentId":null,"classId":"dynamicops.api.model.DiskInputModel","typeFilter":null,"data":{"DISK_CAPACITY":3,"DISK_INPUT_ID":"DISK_INPUT_ID1","DISK_LABEL":"Hard disk 1"}}],"Destroy":true,"EXTERNAL_REFERENCE_ID":"vm-773","IS_COMPONENT_MACHINE":false,"MachineBlueprintName":"CentOS 6.3 - IPAM EXT","MachineCPU":1,"MachineDailyCost":0,"MachineDestructionDate":null,"MachineExpirationDate":null,"MachineGroupName":"Content","MachineGuestOperatingSystem":"CentOS 4/5/6/7 (64-bit)","MachineInterfaceDisplayName":"vSphere (vCenter)","MachineInterfaceType":"vSphere","MachineMemory":512,"MachineName":"Content0061","MachineReservationName":"IPAM Sandbox","MachineStorage":3,"MachineType":"Virtual","NETWORK_LIST":[{"componentTypeId":"com.vmware.csp.component.iaas.proxy.provider","componentId":null,"classId":"dynamicops.api.model.NetworkViewModel","typeFilter":null,"data":{"NETWORK_ADDRESS":"192.168.110.150","NETWORK_MAC_ADDRESS":"00:50:56:ae:31:bd","NETWORK_NAME":"VM Network","NETWORK_NETWORK_NAME":"ipamext1921681100","NETWORK_PROFILE":"ipam-ext-192.168.110.0"}}],"SNAPSHOT_LIST":[],"Unregister":true,"VirtualMachine.Admin.UUID":"502e9fb3-6f0d-0b1e-f90f-a769fd406620","endpointExternalReferenceId":"d322b019-58d4-4d6f-9f8b-d28695a716c0","ip_address":"192.168.110.150","machineId":"4fc33663-992d-49f8-af17-df7ce4831aa0"},"links":[{"@type":"link","rel":"GET: Request","href":"http://localhost/catalog-service/api/consumer/requests/dcb12203-93f4-4873-a7d5-1757f3696141"},{"@type":"link","rel":"GET: Parent Resource","href":"http://localhost/catalog-service/api/consumer/resourceViews/b313acd6-0738-439c-b601-e3ebf9ebb49b"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.virtual.Destroy}","href":"http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/654b4c71-e84f-40c7-9439-fd409fea7323/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.virtual.Destroy}","href":"http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/654b4c71-e84f-40c7-9439-fd409fea7323/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Unregister}","href":"http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/f3ae9408-885a-4a3a-9200-43366f2aa163/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Unregister}","href":"http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/f3ae9408-885a-4a3a-9200-43366f2aa163/requests"}]},{"@type":"CatalogResourceView","resourceId":"169b596f-e4c0-4b25-ba44-18cb19c0fd65","iconId":"existing_network","name":"ipamext1921681100","description":"Infoblox External Network","status":null,"catalogItemId":null,"catalogItemLabel":null,"requestId":"dcb12203-93f4-4873-a7d5-1757f3696141","requestState":"SUCCESSFUL","resourceType":"Infrastructure.Network.Network.Existing","owners":["Jason Cloud Admin"],"businessGroupId":"53619006-56bb-4788-9723-9eab79752cc1","tenantId":"vsphere.local","dateCreated":"2017-07-17T13:27:17.526Z","lastUpdated":"2017-07-17T13:33:25.521Z","lease":{"start":"2017-07-17T13:26:42.079Z","end":null},"costs":null,"costToDate":null,"totalCost":null,"parentResourceId":"b313acd6-0738-439c-b601-e3ebf9ebb49b","hasChildren":false,"data":{"Description":"Infoblox External Network","IPAMEndpointId":"1c2b6237-540a-43c3-8c06-b37a1d274b44","IPAMEndpointName":"Infoblox - nios01a","Name":"ipamext1921681100","_archiveDays":5,"_hasChildren":false,"_leaseDays":null,"_number_of_instances":1,"dns":{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Network.Network.DnsWins","typeFilter":null,"data":{"alternate_wins":null,"dns_search_suffix":null,"dns_suffix":null,"preferred_wins":null,"primary_dns":null,"secondary_dns":null}},"gateway":null,"ip_ranges":[{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Network.Network.IpRanges","typeFilter":null,"data":{"description":"","end_ip":"","externalId":"network/default-vra/192.168.110.0/24","id":"b078d23a-1c3d-4458-ab57-e352c80e6d55","name":"192.168.110.0/24","start_ip":""}}],"network_profile":"ipam-ext-192.168.110.0","providerBindingId":"CentOS63Infoblox","providerId":"2fbaabc5-3a48-488a-9f2a-a42616345445","subnet_mask":"255.255.255.0","subtenantId":"53619006-56bb-4788-9723-9eab79752cc1"},"links":[{"@type":"link","rel":"GET: Request","href":"http://localhost/catalog-service/api/consumer/requests/dcb12203-93f4-4873-a7d5-1757f3696141"},{"@type":"link","rel":"GET: Parent Resource","href":"http://localhost/catalog-service/api/consumer/resourceViews/b313acd6-0738-439c-b601-e3ebf9ebb49b"}]}],"metadata":{"size":20,"totalElements":3,"totalPages":1,"number":1,"offset":0}}`))

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/3da0ca14-e7e2-4d7b-89cb-c6db57440d72/requests/template",
		httpmock.NewStringResponder(200, `{"type":"com.vmware.vcac.catalog.domain.request.CatalogResourceRequest","resourceId":"b313acd6-0738-439c-b601-e3ebf9ebb49b","actionId":"3da0ca14-e7e2-4d7b-89cb-c6db57440d72","description":null,"data":{"ForceDestroy":false}}`))

	templateResources, errTemplate := client.GetDeploymentState("937099db-5174-4862-99a3-9c2666bfca28")
	if errTemplate != nil {
		t.Errorf("Failed to get the template resources %v", errTemplate)
	}

	_, _, err := client.GetDestroyActionTemplate(templateResources)

	if err != nil {
		t.Errorf("Fail to get destroy action template %v", err)
	}

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/3da0ca14-e7e2-4d7b-89cb-c6db57440d72/requests/template",
		httpmock.NewErrorResponder(errors.New(`{"errors":[{"code":50505,"source":null,"message":"System exception.","systemMessage":null,"moreInfoUrl":null}]}`)))

	_, _, err = client.GetDestroyActionTemplate(templateResources)

	if err == nil {
		t.Errorf("Fail to get destroy action template exception.")
	}
}

func TestAPIClient_destroyMachine(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/"+
		"api/consumer/resources/ea68a707-bd35-4765-a1fc-16aca56ff844/actions/356d45a1-d30e-4def-86c9-10cb9aa2b0de/requests/template",
		httpmock.NewStringResponder(200, `{"type":"com.vmware.vcac.catalog.domain.request.CatalogResourceRequest","resourceId":"ea68a707-bd35-4765-a1fc-16aca56ff844","actionId":"356d45a1-d30e-4def-86c9-10cb9aa2b0de","description":null,"reasons":null,"data":{"ForceDestroy":false,"provider-DeleteVolumes":false}}`))

	resourceActionTemplate := new(ResourceActionTemplate)
	apiError := new(APIError)
	//getActionTemplatePath := fmt.Sprintf(utils.GET_ACTION_TEMPLATE_API, "ea68a707-bd35-4765-a1fc-16aca56ff844", "356d45a1-d30e-4def-86c9-10cb9aa2b0de")
	getActionTemplatePath := "http://localhost/catalog-service/api/consumer/resources/ea68a707-bd35-4765-a1fc-16aca56ff844/actions/356d45a1-d30e-4def-86c9-10cb9aa2b0de/requests/template"
	response, err := client.HTTPClient.New().Get(getActionTemplatePath).
		Receive(resourceActionTemplate, apiError)
	log.Info("inside test resource action template : %v ", resourceActionTemplate)
	response.Close = true

	if err != nil {
		log.Errorf("errer %v ", err.Error())
		t.Errorf("Expected no error but found %v ", err.Error())
	}

	if apiError != nil && !apiError.isEmpty() {
		log.Errorf("errer %v ", apiError.Error())
		t.Errorf("Expected no error but found %v ", apiError.Error())
	}

	postActionTemplatePath := "http://localhost/catalog-service/api/consumer/resources/ea68a707-bd35-4765-a1fc-16aca56ff844/actions/356d45a1-d30e-4def-86c9-10cb9aa2b0de/requests"
	httpmock.RegisterResponder("POST", postActionTemplatePath,
		httpmock.NewStringResponder(201, ``))

	err = client.DestroyMachine(resourceActionTemplate, postActionTemplatePath)
	if err != nil {
		t.Errorf("Expected no error but found %v ", err.Error())
	}
}

func TestChangeValueFunction(t *testing.T) {
	request_template_original := CatalogItemRequestTemplate{}
	request_template_backup := CatalogItemRequestTemplate{}
	strJson := `{"type":"com.vmware.vcac.catalog.domain.request.CatalogItemProvisioningRequest","catalogItemId":"e5dd4fba-45ed-4943-b1fc-7f96239286be","requestedFor":"jason@corp.local","businessGroupId":"53619006-56bb-4788-9723-9eab79752cc1","description":null,"reasons":null,"data":{"CentOS_6.3":{"componentTypeId":"com.vmware.csp.component.cafe.composition","componentId":null,"classId":"Blueprint.Component.Declaration","typeFilter":"CentOS63*CentOS_6.3","data":{"_allocation":{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.Allocation","typeFilter":null,"data":{"machines":[{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.Allocation.Machine","typeFilter":null,"data":{"machine_id":"","nics":[{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.Nic","typeFilter":null,"data":{"address":"","assignment_type":"Static","external_address":"","id":null,"load_balancing":null,"network":null,"network_profile":null}}]}}]}},"_cluster":1,"_hasChildren":false,"cpu":1,"datacenter_location":null,"description":"Basic IaaS CentOS Machine","disks":[{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.MachineDisk","typeFilter":null,"data":{"capacity":3,"custom_properties":null,"id":1450725224066,"initial_location":"","is_clone":true,"label":"Hard disk 1","storage_reservation_policy":"","userCreated":false,"volumeId":0}}],"guest_customization_specification":"CentOS","max_network_adapters":-1,"max_per_user":0,"max_volumes":60,"memory":512,"nics":null,"os_arch":"x86_64","os_distribution":null,"os_type":"Linux","os_version":null,"property_groups":null,"reservation_policy":null,"security_groups":[],"security_tags":[],"storage":3}},"_archiveDays":5,"_leaseDays":null,"_number_of_instances":1,"corp192168110024":{"componentTypeId":"com.vmware.csp.component.cafe.composition","componentId":null,"classId":"Blueprint.Component.Declaration","typeFilter":"CentOS63*corp192168110024","data":{"_hasChildren":false}}}}`
	json.Unmarshal([]byte(strJson), &request_template_original)
	json.Unmarshal([]byte(strJson), &request_template_backup)
	var flag bool

	request_template_original.Data, flag = replaceValueInRequestTemplate(request_template_original.Data, "false_field", 1000)
	if flag != false {
		t.Errorf("False value updated")
	}

	eq := reflect.DeepEqual(request_template_backup.Data, request_template_original.Data)
	if !eq {
		t.Errorf("False value updated")
	}

	request_template_original.Data, flag = replaceValueInRequestTemplate(request_template_original.Data, "storage", 1000)
	if flag == false {
		t.Errorf("Failed to update interface value")
	}

	eq2 := reflect.DeepEqual(request_template_backup.Data, request_template_original.Data)
	if eq2 {
		t.Errorf("Failed to update interface value")
	}

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
		utils.CATALOG_ID:             "abcdefghijklmn",
		utils.RESOURCE_CONFIGURATION: mockConfigResourceMap,
	}

	mockResourceData := schema.TestResourceDataRaw(t, resourceSchema, resourceDataMap)

	readProviderConfiguration(mockResourceData)
	err := checkResourceConfigValidity(mockRequestTemplate)
	if err != nil {
		t.Errorf("The terraform config is valid, failed to validate. Expecting no error, but found %v ", err.Error())
	}

	mockConfigResourceMap["machine2.mock.cpu"] = 2
	mockConfigResourceMap["machine2.storage"] = 2

	resourceDataMap = map[string]interface{}{
		utils.CATALOG_ID:             "abcdefghijklmn",
		utils.RESOURCE_CONFIGURATION: mockConfigResourceMap,
	}

	mockResourceData = schema.TestResourceDataRaw(t, resourceSchema, resourceDataMap)
	readProviderConfiguration(mockResourceData)

	err = checkResourceConfigValidity(mockRequestTemplate)
	if err != nil {
		t.Errorf("The terraform config is valid, failed to validate. Expecting no error, but found %v ", err.Error())
	}

	mockConfigResourceMap["mock.machine3.vSphere.mock.cpu"] = 2
	resourceDataMap = map[string]interface{}{
		utils.CATALOG_ID:             "abcdefghijklmn",
		utils.RESOURCE_CONFIGURATION: mockConfigResourceMap,
	}

	mockResourceData = schema.TestResourceDataRaw(t, resourceSchema, resourceDataMap)
	readProviderConfiguration(mockResourceData)

	var mockInvalidKeys []string
	mockInvalidKeys = append(mockInvalidKeys, "mock.machine3.vSphere.mock.cpu")

	validityErr := fmt.Sprintf(utils.CONFIG_INVALID_ERROR, strings.Join(mockInvalidKeys, ", "))
	err = checkResourceConfigValidity(mockRequestTemplate)
	// this should throw an error as none of the string combinations (mock, mock.machine3, mock.machine3.vsphere, etc)
	// matches the component names(mock.test.machine1 and machine2) in the request template
	if err == nil {
		t.Errorf("The terraform config is invalid. failed to validate. Expected the error %v. but found no error", validityErr)
	}

	if err.Error() != validityErr {
		t.Errorf("Expected: %v, but Found: %v", validityErr, err.Error())
	}
}

// creates a mock request template from a request template template json file
func GetMockRequestTemplate() *CatalogItemRequestTemplate {

	ps := utils.GetPathSeparator()
	filePath := os.Getenv("GOPATH") + ps + "src" + ps + "github.com" + ps +
		"vmware" + ps + "terraform-provider-vra7" + ps + "resources" + ps + "MockRequestTemplate"

	absPath, _ := filepath.Abs(filePath)

	jsonFile, err := os.Open(absPath)
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	mockRequestTemplate := CatalogItemRequestTemplate{}
	json.Unmarshal(byteValue, &mockRequestTemplate)

	return &mockRequestTemplate

}
