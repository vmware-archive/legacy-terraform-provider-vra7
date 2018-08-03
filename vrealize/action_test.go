package vrealize

import (
	"errors"
	"fmt"
	"gopkg.in/jarcoal/httpmock.v1"
	"testing"
	"encoding/json"
)

//var client APIClient

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

var resourceViewsResp = `{"links":[],"content":[{"@type":"CatalogResourceView","resourceId":"588f5e28-5572-495e-b104-e1237eaf2b98","iconId":"e5dd4fba-45ed-4943-b1fc-7f96239286be",
"name":"machine 3 with timeout 1 min","description":"","status":null,"catalogItemId":"e5dd4fba-45ed-4943-b1fc-7f96239286be",
"catalogItemLabel":"CentOS 6.3","requestId":"666d77e3-7642-492d-aad1-82b8edd30e56","requestState":"SUCCESSFUL",
"resourceType":"composition.resource.type.deployment","owners":["Jason Cloud Admin"],
"businessGroupId":"53619006-56bb-4788-9723-9eab79752cc1","tenantId":"vsphere.local","dateCreated":"2018-04-20T12:17:32.177Z",
"lastUpdated":"2018-04-20T12:23:03.303Z","lease":{"start":"2018-04-20T12:17:32.164Z","end":null},"costs":null,
"costToDate":null,"totalCost":null,"parentResourceId":null,"hasChildren":true,"data":{},"links":[{"@type":"link",
"rel":"GET: Catalog Item","href":"https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItemViews/e5dd4fba-45ed-4943-b1fc-7f96239286be"},
{"@type":"link","rel":"GET: Request",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/requests/666d77e3-7642-492d-aad1-82b8edd30e56"},
{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changelease.name}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/588f5e28-5572-495e-b104-e1237eaf2b98/actions/561be422-ece6-4316-8acb-a8f3dbb8ed0c/requests/template"},
{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changelease.name}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/588f5e28-5572-495e-b104-e1237eaf2b98/actions/561be422-ece6-4316-8acb-a8f3dbb8ed0c/requests"},
{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changeowner.name}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/588f5e28-5572-495e-b104-e1237eaf2b98/actions/59249166-e427-4082-a3dc-eb7223bb2de1/requests/template"},
{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changeowner.name}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/588f5e28-5572-495e-b104-e1237eaf2b98/actions/59249166-e427-4082-a3dc-eb7223bb2de1/requests"},
{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.destroy.name}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/588f5e28-5572-495e-b104-e1237eaf2b98/actions/3da0ca14-e7e2-4d7b-89cb-c6db57440d72/requests/template"},
{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.destroy.name}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/588f5e28-5572-495e-b104-e1237eaf2b98/actions/3da0ca14-e7e2-4d7b-89cb-c6db57440d72/requests"},
{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.archive.name}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/588f5e28-5572-495e-b104-e1237eaf2b98/actions/9725d56e-461a-471a-be00-b1856681c6d0/requests/template"},
{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.archive.name}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/588f5e28-5572-495e-b104-e1237eaf2b98/actions/9725d56e-461a-471a-be00-b1856681c6d0/requests"},
{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scalein.name}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/588f5e28-5572-495e-b104-e1237eaf2b98/actions/85e090f9-9529-4101-9691-6bab1b0a1f77/requests/template"},
{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scalein.name}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/588f5e28-5572-495e-b104-e1237eaf2b98/actions/85e090f9-9529-4101-9691-6bab1b0a1f77/requests"},
{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scaleout.name}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/588f5e28-5572-495e-b104-e1237eaf2b98/actions/ab5795f5-32ad-4f6c-8598-1d3a7d190caa/requests/template"},
{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scaleout.name}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/588f5e28-5572-495e-b104-e1237eaf2b98/actions/ab5795f5-32ad-4f6c-8598-1d3a7d190caa/requests"},
{"@type":"link","rel":"GET: Child Resources",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resourceViews?managedOnly=false&withExtendedData=true&withOperations=true&%24filter=parentResource%20eq%20%27588f5e28-5572-495e-b104-e1237eaf2b98%27"}]},
{"@type":"CatalogResourceView","resourceId":"4f58732f-62c7-4d38-a78b-b2cf34ee45df",
"iconId":"Infrastructure.CatalogItem.Machine.Virtual.vSphere","name":"dev-444","description":"Basic IaaS CentOS Machine",
"status":"On","catalogItemId":null,"catalogItemLabel":null,"requestId":"666d77e3-7642-492d-aad1-82b8edd30e56",
"requestState":"SUCCESSFUL","resourceType":"Infrastructure.Virtual","owners":["Jason Cloud Admin"],
"businessGroupId":"53619006-56bb-4788-9723-9eab79752cc1","tenantId":"vsphere.local","dateCreated":"2018-04-20T12:22:51.322Z",
"lastUpdated":"2018-04-20T12:23:03.303Z","lease":{"start":"2018-04-20T12:17:32.164Z","end":null},"costs":null,"costToDate":null,
"totalCost":null,"parentResourceId":"588f5e28-5572-495e-b104-e1237eaf2b98","hasChildren":false,"data":{"ChangeLease":true,
"ChangeOwner":true,"Component":"CentOS_6.3","ConnectViaNativeVmrc":true,"ConnectViaVmrc":true,"CreateSnapshot":true,
"DISK_VOLUMES":[{"componentTypeId":"com.vmware.csp.component.iaas.proxy.provider","componentId":null,
"classId":"dynamicops.api.model.DiskInputModel","typeFilter":null,"data":{"DISK_CAPACITY":3,"DISK_INPUT_ID":"DISK_INPUT_ID1",
"DISK_LABEL":"Hard disk 1"}}],"Destroy":true,"EXTERNAL_REFERENCE_ID":"vm-989","Expire":true,"IS_COMPONENT_MACHINE":false,
"InstallTools":true,"MachineBlueprintName":"CentOS 6.3","MachineCPU":1,"MachineDailyCost":0,"MachineDestructionDate":null,
"MachineExpirationDate":null,"MachineGroupName":"Content","MachineGuestOperatingSystem":"CentOS 4/5/6/7 (64-bit)",
"MachineInterfaceDisplayName":"vSphere (vCenter)","MachineInterfaceType":"vSphere","MachineMemory":512,"MachineName":"dev-444",
"MachineReservationName":"Content Cluster Reservation","MachineStorage":3,"MachineType":"Virtual",
"NETWORK_LIST":[{"componentTypeId":"com.vmware.csp.component.iaas.proxy.provider","componentId":null,
"classId":"dynamicops.api.model.NetworkViewModel","typeFilter":null,"data":{"NETWORK_MAC_ADDRESS":"00:50:56:ae:e5:87",
"NETWORK_NAME":"VM Network","NETWORK_NETWORK_NAME":"corp192168110024","NETWORK_PROFILE":"corp-192.168.110.0/24"}}],
"PowerOff":true,"Reboot":true,"Reconfigure":true,"Reprovision":true,"Reset":true,"SNAPSHOT_LIST":[],"Shutdown":true,
"Suspend":true,"VirtualMachine.Admin.UUID":"502ee2ef-d81f-6965-7d3a-08e23291ace5",
"endpointExternalReferenceId":"d322b019-58d4-4d6f-9f8b-d28695a716c0","ip_address":"192.168.100.136",
"machineId":"7c4c92ae-4c00-45a1-9664-f26a8754ae66"},"links":[{"@type":"link","rel":"GET: Request",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/requests/666d77e3-7642-492d-aad1-82b8edd30e56"},{"@type":"link",
"rel":"GET: Parent Resource","href":"https://vra-01a.corp.local/catalog-service/api/consumer/resourceViews/588f5e28-5572-495e-b104-e1237eaf2b98"},
{"@type":"link","rel":"GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.virtual.CreateSnapshot}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/06b9e681-0f76-4f95-90b3-6e657f5fbf23/requests/template"},
{"@type":"link","rel":"POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.virtual.CreateSnapshot}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/06b9e681-0f76-4f95-90b3-6e657f5fbf23/requests"},
{"@type":"link","rel":"GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.virtual.Destroy}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/654b4c71-e84f-40c7-9439-fd409fea7323/requests/template"},
{"@type":"link","rel":"POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.virtual.Destroy}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/654b4c71-e84f-40c7-9439-fd409fea7323/requests"},
{"@type":"link","rel":"GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.InstallTools}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/8dcec322-db95-451f-ad56-ac37e406672a/requests/template"},
{"@type":"link","rel":"POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.InstallTools}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/8dcec322-db95-451f-ad56-ac37e406672a/requests"},
{"@type":"link","rel":"GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Reset}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/f928ebb1-6bfd-46b6-b912-46ed83facd4b/requests/template"},
{"@type":"link","rel":"POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Reset}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/f928ebb1-6bfd-46b6-b912-46ed83facd4b/requests"},
{"@type":"link","rel":"GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.PowerOff}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/b37c071e-06ce-4842-b194-0f64a829908f/requests/template"},
{"@type":"link","rel":"POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.PowerOff}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/b37c071e-06ce-4842-b194-0f64a829908f/requests"},
{"@type":"link","rel":"GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Reboot}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/45395e7d-75f1-4829-957b-64c5538c667d/requests/template"},
{"@type":"link","rel":"POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Reboot}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/45395e7d-75f1-4829-957b-64c5538c667d/requests"},
{"@type":"link","rel":"GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Reconfigure}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/eb98e7d9-9de2-4600-9888-0c0f0d6d696d/requests/template"},
{"@type":"link","rel":"POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Reconfigure}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/eb98e7d9-9de2-4600-9888-0c0f0d6d696d/requests"},
{"@type":"link","rel":"GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Reprovision}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/1cd874e9-90d6-4a87-b75f-b2aa584fae0e/requests/template"},
{"@type":"link","rel":"POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Reprovision}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/1cd874e9-90d6-4a87-b75f-b2aa584fae0e/requests"},
{"@type":"link","rel":"GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Shutdown}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/17b3729a-da19-495b-bb59-efbfb028695d/requests/template"},
{"@type":"link","rel":"POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Shutdown}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/17b3729a-da19-495b-bb59-efbfb028695d/requests"},
{"@type":"link","rel":"GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Suspend}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/56d16cd0-4e4f-400f-93ff-acd550d40aee/requests/template"},
{"@type":"link","rel":"POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Suspend}",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/56d16cd0-4e4f-400f-93ff-acd550d40aee/requests"}]},
{"@type":"CatalogResourceView","resourceId":"97556f48-bb95-4687-a9d8-273d0263716b","iconId":"existing_network",
"name":"corp192168110024","description":"Default External Network Profile (VM Network, vSphere)","status":null,
"catalogItemId":null,"catalogItemLabel":null,"requestId":"666d77e3-7642-492d-aad1-82b8edd30e56","requestState":"SUCCESSFUL",
"resourceType":"Infrastructure.Network.Network.Existing","owners":["Jason Cloud Admin"],
"businessGroupId":"53619006-56bb-4788-9723-9eab79752cc1","tenantId":"vsphere.local","dateCreated":"2018-04-20T12:17:48.181Z",
"lastUpdated":"2018-04-20T12:23:03.303Z","lease":{"start":"2018-04-20T12:17:32.164Z","end":null},"costs":null,"costToDate":null,
"totalCost":null,"parentResourceId":"588f5e28-5572-495e-b104-e1237eaf2b98","hasChildren":false,
"data":{"Description":"Default External Network Profile (VM Network, vSphere)","Name":"corp192168110024","_archiveDays":5,
"_deploymentName":"machine 3 with timeout 1 min","_hasChildren":false,"_leaseDays":null,"_number_of_instances":1,
"dns":{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,
"classId":"Infrastructure.Network.Network.DnsWins","typeFilter":null,"data":{"alternate_wins":null,
"dns_search_suffix":"corp.local","dns_suffix":"corp.local","preferred_wins":null,"primary_dns":"192.168.110.10",
"secondary_dns":null}},"gateway":"192.168.110.1","ip_ranges":[{"componentTypeId":"com.vmware.csp.iaas.blueprint.service",
"componentId":null,"classId":"Infrastructure.Network.Network.IpRanges","typeFilter":null,"data":{"description":"",
"end_ip":"192.168.110.250","externalId":"","id":"f91f513e-9ed0-4b43-bd25-fcdfe9ea0870","name":"IP Range",
"start_ip":"192.168.110.200"}}],"network_profile":"corp-192.168.110.0/24","providerBindingId":"CentOS63",
"providerId":"2fbaabc5-3a48-488a-9f2a-a42616345445","subnet_mask":"255.255.255.0",
"subtenantId":"53619006-56bb-4788-9723-9eab79752cc1"},"links":[{"@type":"link","rel":"GET: Request",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/requests/666d77e3-7642-492d-aad1-82b8edd30e56"},
{"@type":"link","rel":"GET: Parent Resource",
"href":"https://vra-01a.corp.local/catalog-service/api/consumer/resourceViews/588f5e28-5572-495e-b104-e1237eaf2b98"}]}],
"metadata":{"size":20,"totalElements":3,"totalPages":1,"number":1,"offset":0}}`

var catalogItemId = "666d77e3-7642-492d-aad1-82b8edd30e56"

func TestPowerOffAction(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/requests/"+catalogItemId+"/resourceViews",
		httpmock.NewStringResponder(200, resourceViewsResp))

	httpmock.RegisterResponder("GET", "https://vra-01a.corp.local/catalog-service/api/consumer/resources/4f58732f-62c7-4d38-a78b-b2cf34ee45df/actions/b37c071e-06ce-4842-b194-0f64a829908f/requests/template",
		httpmock.NewStringResponder(200, `{"type":"com.vmware.vcac.catalog.domain.request.CatalogResourceRequest",
			"resourceId":"4f58732f-62c7-4d38-a78b-b2cf34ee45df","actionId":"b37c071e-06ce-4842-b194-0f64a829908f","description":null,
			"data":{}}`))

	templateResources, errTemplate := client.GetDeploymentState(catalogItemId)
	if errTemplate != nil {
		t.Errorf("Failed to get the template resources %v", errTemplate)
	}

	_, _, err := client.GetPowerOffActionTemplate(templateResources)

	if err != nil {
		t.Errorf("Fail to get destroy action template %v", err)
	}

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/3da0ca14-e7e2-4d7b-89cb-c6db57440d72/requests/template",
		httpmock.NewErrorResponder(errors.New(`{"errors":[{"code":50505,"source":null,"message":"System exception.",
			"systemMessage":null,"moreInfoUrl":null}]}`)))

	_, _, err = client.GetPowerOffActionTemplate(templateResources)

	if err == nil {
		t.Errorf("Fail to get destroy action template exception.")
	}
}

func TestGetActionTemplate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	resourceView := new(ResourceView)
	str := `{"links":[],"content":[{"@type":"CatalogResourceView","resourceId":"b313acd6-0738-439c-b601-e3ebf9ebb49b","iconId":"502efc1b-d5ce-4ef9-99ee-d4e2a741747c","name":"CentOS 6.3 - IPAM EXT-95563173","description":"","status":null,"catalogItemId":"502efc1b-d5ce-4ef9-99ee-d4e2a741747c","catalogItemLabel":"CentOS 6.3 - IPAM EXT","requestId":"dcb12203-93f4-4873-a7d5-1757f3696141","requestState":"SUCCESSFUL","resourceType":"composition.resource.type.deployment","owners":["Jason Cloud Admin"],"businessGroupId":"53619006-56bb-4788-9723-9eab79752cc1","tenantId":"vsphere.local","dateCreated":"2017-07-17T13:26:42.102Z","lastUpdated":"2017-07-17T13:33:25.521Z","lease":{"start":"2017-07-17T13:26:42.079Z","end":null},"costs":null,"costToDate":null,"totalCost":null,"parentResourceId":null,"hasChildren":true,"data":{},"links":[{"@type":"link","rel":"GET: Catalog Item","href":"http://localhost/catalog-service/api/consumer/entitledCatalogItemViews/502efc1b-d5ce-4ef9-99ee-d4e2a741747c"},{"@type":"link","rel":"GET: Request","href":"http://localhost/catalog-service/api/consumer/requests/dcb12203-93f4-4873-a7d5-1757f3696141"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changelease.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/561be422-ece6-4316-8acb-a8f3dbb8ed0c/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changelease.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/561be422-ece6-4316-8acb-a8f3dbb8ed0c/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changeowner.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/59249166-e427-4082-a3dc-eb7223bb2de1/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changeowner.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/59249166-e427-4082-a3dc-eb7223bb2de1/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.destroy.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/3da0ca14-e7e2-4d7b-89cb-c6db57440d72/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.destroy.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/3da0ca14-e7e2-4d7b-89cb-c6db57440d72/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.archive.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/9725d56e-461a-471a-be00-b1856681c6d0/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.archive.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/9725d56e-461a-471a-be00-b1856681c6d0/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scalein.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/85e090f9-9529-4101-9691-6bab1b0a1f77/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scalein.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/85e090f9-9529-4101-9691-6bab1b0a1f77/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scaleout.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/ab5795f5-32ad-4f6c-8598-1d3a7d190caa/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scaleout.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/ab5795f5-32ad-4f6c-8598-1d3a7d190caa/requests"},{"@type":"link","rel":"GET: Child Resources","href":"http://localhost/catalog-service/api/consumer/resourceViews?managedOnly=false&withExtendedData=true&withOperations=true&%24filter=parentResource%20eq%20%27b313acd6-0738-439c-b601-e3ebf9ebb49b%27"}]},{"@type":"CatalogResourceView","resourceId":"51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5","iconId":"Infrastructure.CatalogItem.Machine.Virtual.vSphere","name":"Content0061","description":"Basic IaaS CentOS Machine","status":"Missing","catalogItemId":null,"catalogItemLabel":null,"requestId":"dcb12203-93f4-4873-a7d5-1757f3696141","requestState":"SUCCESSFUL","resourceType":"Infrastructure.Virtual","owners":["Jason Cloud Admin"],"businessGroupId":"53619006-56bb-4788-9723-9eab79752cc1","tenantId":"vsphere.local","dateCreated":"2017-07-17T13:33:16.686Z","lastUpdated":"2017-07-17T13:33:25.521Z","lease":{"start":"2017-07-17T13:26:42.079Z","end":null},"costs":null,"costToDate":null,"totalCost":null,"parentResourceId":"b313acd6-0738-439c-b601-e3ebf9ebb49b","hasChildren":false,"data":{"Component":"CentOS_6.3","DISK_VOLUMES":[{"componentTypeId":"com.vmware.csp.component.iaas.proxy.provider","componentId":null,"classId":"dynamicops.api.model.DiskInputModel","typeFilter":null,"data":{"DISK_CAPACITY":3,"DISK_INPUT_ID":"DISK_INPUT_ID1","DISK_LABEL":"Hard disk 1"}}],"Destroy":true,"EXTERNAL_REFERENCE_ID":"vm-773","IS_COMPONENT_MACHINE":false,"MachineBlueprintName":"CentOS 6.3 - IPAM EXT","MachineCPU":1,"MachineDailyCost":0,"MachineDestructionDate":null,"MachineExpirationDate":null,"MachineGroupName":"Content","MachineGuestOperatingSystem":"CentOS 4/5/6/7 (64-bit)","MachineInterfaceDisplayName":"vSphere (vCenter)","MachineInterfaceType":"vSphere","MachineMemory":512,"MachineName":"Content0061","MachineReservationName":"IPAM Sandbox","MachineStorage":3,"MachineType":"Virtual","NETWORK_LIST":[{"componentTypeId":"com.vmware.csp.component.iaas.proxy.provider","componentId":null,"classId":"dynamicops.api.model.NetworkViewModel","typeFilter":null,"data":{"NETWORK_ADDRESS":"192.168.110.150","NETWORK_MAC_ADDRESS":"00:50:56:ae:31:bd","NETWORK_NAME":"VM Network","NETWORK_NETWORK_NAME":"ipamext1921681100","NETWORK_PROFILE":"ipam-ext-192.168.110.0"}}],"SNAPSHOT_LIST":[],"Unregister":true,"VirtualMachine.Admin.UUID":"502e9fb3-6f0d-0b1e-f90f-a769fd406620","endpointExternalReferenceId":"d322b019-58d4-4d6f-9f8b-d28695a716c0","ip_address":"192.168.110.150","machineId":"4fc33663-992d-49f8-af17-df7ce4831aa0"},"links":[{"@type":"link","rel":"GET: Request","href":"http://localhost/catalog-service/api/consumer/requests/dcb12203-93f4-4873-a7d5-1757f3696141"},{"@type":"link","rel":"GET: Parent Resource","href":"http://localhost/catalog-service/api/consumer/resourceViews/b313acd6-0738-439c-b601-e3ebf9ebb49b"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.virtual.Destroy}","href":"http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/654b4c71-e84f-40c7-9439-fd409fea7323/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.virtual.Destroy}","href":"http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/654b4c71-e84f-40c7-9439-fd409fea7323/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Unregister}","href":"http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/f3ae9408-885a-4a3a-9200-43366f2aa163/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Unregister}","href":"http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/f3ae9408-885a-4a3a-9200-43366f2aa163/requests"}]},{"@type":"CatalogResourceView","resourceId":"169b596f-e4c0-4b25-ba44-18cb19c0fd65","iconId":"existing_network","name":"ipamext1921681100","description":"Infoblox External Network","status":null,"catalogItemId":null,"catalogItemLabel":null,"requestId":"dcb12203-93f4-4873-a7d5-1757f3696141","requestState":"SUCCESSFUL","resourceType":"Infrastructure.Network.Network.Existing","owners":["Jason Cloud Admin"],"businessGroupId":"53619006-56bb-4788-9723-9eab79752cc1","tenantId":"vsphere.local","dateCreated":"2017-07-17T13:27:17.526Z","lastUpdated":"2017-07-17T13:33:25.521Z","lease":{"start":"2017-07-17T13:26:42.079Z","end":null},"costs":null,"costToDate":null,"totalCost":null,"parentResourceId":"b313acd6-0738-439c-b601-e3ebf9ebb49b","hasChildren":false,"data":{"Description":"Infoblox External Network","IPAMEndpointId":"1c2b6237-540a-43c3-8c06-b37a1d274b44","IPAMEndpointName":"Infoblox - nios01a","Name":"ipamext1921681100","_archiveDays":5,"_hasChildren":false,"_leaseDays":null,"_number_of_instances":1,"dns":{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Network.Network.DnsWins","typeFilter":null,"data":{"alternate_wins":null,"dns_search_suffix":null,"dns_suffix":null,"preferred_wins":null,"primary_dns":null,"secondary_dns":null}},"gateway":null,"ip_ranges":[{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Network.Network.IpRanges","typeFilter":null,"data":{"description":"","end_ip":"","externalId":"network/default-vra/192.168.110.0/24","id":"b078d23a-1c3d-4458-ab57-e352c80e6d55","name":"192.168.110.0/24","start_ip":""}}],"network_profile":"ipam-ext-192.168.110.0","providerBindingId":"CentOS63Infoblox","providerId":"2fbaabc5-3a48-488a-9f2a-a42616345445","subnet_mask":"255.255.255.0","subtenantId":"53619006-56bb-4788-9723-9eab79752cc1"},"links":[{"@type":"link","rel":"GET: Request","href":"http://localhost/catalog-service/api/consumer/requests/dcb12203-93f4-4873-a7d5-1757f3696141"},{"@type":"link","rel":"GET: Parent Resource","href":"http://localhost/catalog-service/api/consumer/resourceViews/b313acd6-0738-439c-b601-e3ebf9ebb49b"}]}],"metadata":{"size":20,"totalElements":3,"totalPages":1,"number":1,"offset":0}}`
	json.Unmarshal([]byte(str), &resourceView)
	client.GetActionTemplate(resourceView, "")

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/3da0ca14-e7e2-4d7b-89cb-c6db57440d72/requests/template",
		httpmock.NewErrorResponder(errors.New(`{"errors":[{"code":50505,"source":null,"message":"System exception.","systemMessage":null,"moreInfoUrl":null}]}`)))
	actionURL := "GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.destroy.name}"
	_,_, err := client.GetActionTemplate(resourceView, actionURL)
	if err == nil {
		t.Errorf("error expected while fetching action template")
	}
}

func TestDeploymentLeaseActionLinks(t *testing.T) {
	// check deployment action links with false resource data
	resourceView := new(ResourceView)
	actionLinks1 := client.getDeploymentLeaseActionLinks(resourceView)
	// action links mapn should be empty with false resource map / empty resource map
	if len(actionLinks1) != 0 {
		t.Errorf("empty action link map expected")
	}
	str := `{"links":[],"content":[{"@type":"CatalogResourceView","resourceId":"b313acd6-0738-439c-b601-e3ebf9ebb49b","iconId":"502efc1b-d5ce-4ef9-99ee-d4e2a741747c","name":"CentOS 6.3 - IPAM EXT-95563173","description":"","status":null,"catalogItemId":"502efc1b-d5ce-4ef9-99ee-d4e2a741747c","catalogItemLabel":"CentOS 6.3 - IPAM EXT","requestId":"dcb12203-93f4-4873-a7d5-1757f3696141","requestState":"SUCCESSFUL","resourceType":"composition.resource.type.deployment","owners":["Jason Cloud Admin"],"businessGroupId":"53619006-56bb-4788-9723-9eab79752cc1","tenantId":"vsphere.local","dateCreated":"2017-07-17T13:26:42.102Z","lastUpdated":"2017-07-17T13:33:25.521Z","lease":{"start":"2017-07-17T13:26:42.079Z","end":null},"costs":null,"costToDate":null,"totalCost":null,"parentResourceId":null,"hasChildren":true,"data":{},"links":[{"@type":"link","rel":"GET: Catalog Item","href":"http://localhost/catalog-service/api/consumer/entitledCatalogItemViews/502efc1b-d5ce-4ef9-99ee-d4e2a741747c"},{"@type":"link","rel":"GET: Request","href":"http://localhost/catalog-service/api/consumer/requests/dcb12203-93f4-4873-a7d5-1757f3696141"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changelease.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/561be422-ece6-4316-8acb-a8f3dbb8ed0c/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changelease.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/561be422-ece6-4316-8acb-a8f3dbb8ed0c/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changeowner.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/59249166-e427-4082-a3dc-eb7223bb2de1/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changeowner.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/59249166-e427-4082-a3dc-eb7223bb2de1/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.destroy.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/3da0ca14-e7e2-4d7b-89cb-c6db57440d72/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.destroy.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/3da0ca14-e7e2-4d7b-89cb-c6db57440d72/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.archive.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/9725d56e-461a-471a-be00-b1856681c6d0/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.archive.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/9725d56e-461a-471a-be00-b1856681c6d0/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scalein.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/85e090f9-9529-4101-9691-6bab1b0a1f77/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scalein.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/85e090f9-9529-4101-9691-6bab1b0a1f77/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scaleout.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/ab5795f5-32ad-4f6c-8598-1d3a7d190caa/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scaleout.name}","href":"http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/ab5795f5-32ad-4f6c-8598-1d3a7d190caa/requests"},{"@type":"link","rel":"GET: Child Resources","href":"http://localhost/catalog-service/api/consumer/resourceViews?managedOnly=false&withExtendedData=true&withOperations=true&%24filter=parentResource%20eq%20%27b313acd6-0738-439c-b601-e3ebf9ebb49b%27"}]},{"@type":"CatalogResourceView","resourceId":"51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5","iconId":"Infrastructure.CatalogItem.Machine.Virtual.vSphere","name":"Content0061","description":"Basic IaaS CentOS Machine","status":"Missing","catalogItemId":null,"catalogItemLabel":null,"requestId":"dcb12203-93f4-4873-a7d5-1757f3696141","requestState":"SUCCESSFUL","resourceType":"Infrastructure.Virtual","owners":["Jason Cloud Admin"],"businessGroupId":"53619006-56bb-4788-9723-9eab79752cc1","tenantId":"vsphere.local","dateCreated":"2017-07-17T13:33:16.686Z","lastUpdated":"2017-07-17T13:33:25.521Z","lease":{"start":"2017-07-17T13:26:42.079Z","end":null},"costs":null,"costToDate":null,"totalCost":null,"parentResourceId":"b313acd6-0738-439c-b601-e3ebf9ebb49b","hasChildren":false,"data":{"Component":"CentOS_6.3","DISK_VOLUMES":[{"componentTypeId":"com.vmware.csp.component.iaas.proxy.provider","componentId":null,"classId":"dynamicops.api.model.DiskInputModel","typeFilter":null,"data":{"DISK_CAPACITY":3,"DISK_INPUT_ID":"DISK_INPUT_ID1","DISK_LABEL":"Hard disk 1"}}],"Destroy":true,"EXTERNAL_REFERENCE_ID":"vm-773","IS_COMPONENT_MACHINE":false,"MachineBlueprintName":"CentOS 6.3 - IPAM EXT","MachineCPU":1,"MachineDailyCost":0,"MachineDestructionDate":null,"MachineExpirationDate":null,"MachineGroupName":"Content","MachineGuestOperatingSystem":"CentOS 4/5/6/7 (64-bit)","MachineInterfaceDisplayName":"vSphere (vCenter)","MachineInterfaceType":"vSphere","MachineMemory":512,"MachineName":"Content0061","MachineReservationName":"IPAM Sandbox","MachineStorage":3,"MachineType":"Virtual","NETWORK_LIST":[{"componentTypeId":"com.vmware.csp.component.iaas.proxy.provider","componentId":null,"classId":"dynamicops.api.model.NetworkViewModel","typeFilter":null,"data":{"NETWORK_ADDRESS":"192.168.110.150","NETWORK_MAC_ADDRESS":"00:50:56:ae:31:bd","NETWORK_NAME":"VM Network","NETWORK_NETWORK_NAME":"ipamext1921681100","NETWORK_PROFILE":"ipam-ext-192.168.110.0"}}],"SNAPSHOT_LIST":[],"Unregister":true,"VirtualMachine.Admin.UUID":"502e9fb3-6f0d-0b1e-f90f-a769fd406620","endpointExternalReferenceId":"d322b019-58d4-4d6f-9f8b-d28695a716c0","ip_address":"192.168.110.150","machineId":"4fc33663-992d-49f8-af17-df7ce4831aa0"},"links":[{"@type":"link","rel":"GET: Request","href":"http://localhost/catalog-service/api/consumer/requests/dcb12203-93f4-4873-a7d5-1757f3696141"},{"@type":"link","rel":"GET: Parent Resource","href":"http://localhost/catalog-service/api/consumer/resourceViews/b313acd6-0738-439c-b601-e3ebf9ebb49b"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.virtual.Destroy}","href":"http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/654b4c71-e84f-40c7-9439-fd409fea7323/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.virtual.Destroy}","href":"http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/654b4c71-e84f-40c7-9439-fd409fea7323/requests"},{"@type":"link","rel":"GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Unregister}","href":"http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/f3ae9408-885a-4a3a-9200-43366f2aa163/requests/template"},{"@type":"link","rel":"POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Unregister}","href":"http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/f3ae9408-885a-4a3a-9200-43366f2aa163/requests"}]},{"@type":"CatalogResourceView","resourceId":"169b596f-e4c0-4b25-ba44-18cb19c0fd65","iconId":"existing_network","name":"ipamext1921681100","description":"Infoblox External Network","status":null,"catalogItemId":null,"catalogItemLabel":null,"requestId":"dcb12203-93f4-4873-a7d5-1757f3696141","requestState":"SUCCESSFUL","resourceType":"Infrastructure.Network.Network.Existing","owners":["Jason Cloud Admin"],"businessGroupId":"53619006-56bb-4788-9723-9eab79752cc1","tenantId":"vsphere.local","dateCreated":"2017-07-17T13:27:17.526Z","lastUpdated":"2017-07-17T13:33:25.521Z","lease":{"start":"2017-07-17T13:26:42.079Z","end":null},"costs":null,"costToDate":null,"totalCost":null,"parentResourceId":"b313acd6-0738-439c-b601-e3ebf9ebb49b","hasChildren":false,"data":{"Description":"Infoblox External Network","IPAMEndpointId":"1c2b6237-540a-43c3-8c06-b37a1d274b44","IPAMEndpointName":"Infoblox - nios01a","Name":"ipamext1921681100","_archiveDays":5,"_hasChildren":false,"_leaseDays":null,"_number_of_instances":1,"dns":{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Network.Network.DnsWins","typeFilter":null,"data":{"alternate_wins":null,"dns_search_suffix":null,"dns_suffix":null,"preferred_wins":null,"primary_dns":null,"secondary_dns":null}},"gateway":null,"ip_ranges":[{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Network.Network.IpRanges","typeFilter":null,"data":{"description":"","end_ip":"","externalId":"network/default-vra/192.168.110.0/24","id":"b078d23a-1c3d-4458-ab57-e352c80e6d55","name":"192.168.110.0/24","start_ip":""}}],"network_profile":"ipam-ext-192.168.110.0","providerBindingId":"CentOS63Infoblox","providerId":"2fbaabc5-3a48-488a-9f2a-a42616345445","subnet_mask":"255.255.255.0","subtenantId":"53619006-56bb-4788-9723-9eab79752cc1"},"links":[{"@type":"link","rel":"GET: Request","href":"http://localhost/catalog-service/api/consumer/requests/dcb12203-93f4-4873-a7d5-1757f3696141"},{"@type":"link","rel":"GET: Parent Resource","href":"http://localhost/catalog-service/api/consumer/resourceViews/b313acd6-0738-439c-b601-e3ebf9ebb49b"}]}],"metadata":{"size":20,"totalElements":3,"totalPages":1,"number":1,"offset":0}}`
	json.Unmarshal([]byte(str), &resourceView)

	// check deployment action links with false resource data
	actionLinks1 = client.getDeploymentLeaseActionLinks(resourceView)

	// action links mapn should not be empty with correct resource map
	if len(actionLinks1) == 0 {
		t.Errorf("empty action link map was not expected")
	}
}