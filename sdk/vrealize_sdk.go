package sdk

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vmware/terraform-provider-vra7/utils"
)

// API constants
const (
	IdentityAPI                 = "/identity/api"
	Tokens                      = IdentityAPI + "/tokens"
	Tenants                     = IdentityAPI + "/tenants"
	CatalogService              = "/catalog-service"
	CatalogServiceAPI           = CatalogService + "/api"
	Consumer                    = CatalogServiceAPI + "/consumer"
	ConsumerRequests            = Consumer + "/requests"
	ConsumerResources           = Consumer + "/resources"
	EntitledCatalogItems        = Consumer + "/entitledCatalogItems"
	EntitledCatalogItemViewsAPI = Consumer + "/entitledCatalogItemViews"
	GetResourceAPI              = ConsumerRequests + "/" + "%s" + "/resources"
	PostActionTemplateAPI       = ConsumerResources + "/" + "%s" + "/actions/" + "%s" + "/requests"
	GetActionTemplateAPI        = PostActionTemplateAPI + "/template"
	GetRequestResourceViewAPI   = ConsumerRequests + "/" + "%s" + "/resourceViews"
	RequestTemplateAPI          = EntitledCatalogItems + "/" + "%s" + "/requests/template"

	// read resource machine constants

	MachineCPU             = "cpu"
	MachineStorage         = "storage"
	MachineMemory          = "memory"
	IPAddress              = "ip_address"
	MachineName            = "name"
	MachineGuestOs         = "guest_operating_system"
	MachineBpName          = "blueprint_name"
	MachineType            = "type"
	MachineReservationName = "reservation_name"
	MachineInterfaceType   = "interface_type"
	MachineID              = "id"
	MachineGroupName       = "group_name"
	MachineDestructionDate = "destruction_date"
	MachineReconfigure     = "reconfigure"
	MachinePowerOff        = "power_off"

	InProgress             = "IN_PROGRESS"
	Successful             = "SUCCESSFUL"
	Failed                 = "FAILED"
	Submitted              = "SUBMITTED"
	InfrastructureVirtual  = "Infrastructure.Virtual"
	DeploymentResourceType = "composition.resource.type.deployment"
	Component              = "Component"
	Reconfigure            = "Reconfigure"
	Destroy                = "Destroy"
)

//GetCatalogItemRequestTemplate - Call to retrieve a request template for a catalog item.
func (c *APIClient) GetCatalogItemRequestTemplate(catalogItemID string) (*CatalogItemRequestTemplate, error) {

	//Form a path to read catalog request template via REST call
	path := fmt.Sprintf(RequestTemplateAPI, catalogItemID)
	url := c.BuildEncodedURL(path, nil)
	resp, respErr := c.Get(url, nil)
	if respErr != nil {
		return nil, respErr
	}

	var requestTemplate CatalogItemRequestTemplate
	unmarshallErr := utils.UnmarshalJSON(resp.Body, &requestTemplate)
	if unmarshallErr != nil {
		return nil, unmarshallErr
	}
	return &requestTemplate, nil
}

// ReadCatalogItemNameByID - This function returns the catalog item name using catalog item ID
func (c *APIClient) ReadCatalogItemNameByID(catalogItemID string) (string, error) {

	path := fmt.Sprintf(EntitledCatalogItems+"/"+"%s", catalogItemID)
	url := c.BuildEncodedURL(path, nil)
	resp, respErr := c.Get(url, nil)
	if respErr != nil {
		return "", respErr
	}

	var response CatalogItem
	unmarshallErr := utils.UnmarshalJSON(resp.Body, &response)
	if unmarshallErr != nil {
		return "", unmarshallErr
	}
	return response.CatalogItem.Name, nil
}

// ReadCatalogItemByName to read id of catalog from vRA using catalog_name
func (c *APIClient) ReadCatalogItemByName(catalogName string) (string, error) {
	var catalogItemID string

	log.Info("readCatalogItemIdByName->catalog_name %v\n", catalogName)

	//Set a call to read number of catalogs from vRA
	path := fmt.Sprintf(EntitledCatalogItemViewsAPI)

	url := c.BuildEncodedURL(path, nil)
	resp, respErr := c.Get(url, nil)
	if respErr != nil || resp.StatusCode != 200 {
		return "", respErr
	}

	var template EntitledCatalogItemViews
	unmarshallErr := utils.UnmarshalJSON(resp.Body, &template)
	if unmarshallErr != nil {
		return "", unmarshallErr
	}

	var catalogItemNameArray []string
	interfaceArray := template.Content.([]interface{})
	catalogItemNameLen := len(catalogName)

	//Iterate over all catalog results to find out matching catalog name
	// provided in terraform configuration file
	for i := range interfaceArray {
		catalogItem := interfaceArray[i].(map[string]interface{})
		if catalogItemNameLen <= len(catalogItem["name"].(string)) {
			//If exact name matches then return respective catalog_id
			//else if provided catalog matches as a substring in name then store it in array
			if catalogName == catalogItem["name"].(string) {
				return catalogItem["catalogItemId"].(string), nil
			} else if catalogName == catalogItem["name"].(string)[0:catalogItemNameLen] {
				catalogItemNameArray = append(catalogItemNameArray, catalogItem["name"].(string))
			}
		}
	}

	// If multiple catalog items are present with provided catalog_name
	// then raise an error and show all names of catalog items with similar name
	if len(catalogItemNameArray) > 0 {
		for index := range catalogItemNameArray {
			catalogItemNameArray[index] = strconv.Itoa(index+1) + " " + catalogItemNameArray[index]
		}
		errorMessage := strings.Join(catalogItemNameArray, "\n")
		punctuation := "is"
		if len(catalogItemNameArray) > 1 {
			punctuation = "are"
		}
		return "", fmt.Errorf("There %s total %d catalog(s) present with same name.\n%s\n"+
			"Please select from above.", punctuation, len(catalogItemNameArray), errorMessage)
	}
	return catalogItemID, nil
}

// GetBusinessGroupID retrieves business group id from business group name
func (c *APIClient) GetBusinessGroupID(businessGroupName string, tenant string) (string, error) {

	path := Tenants + "/" + tenant + "/subtenants"

	log.Info("Fetching business group id from name..GET %s ", path)

	url := c.BuildEncodedURL(path, nil)

	resp, respErr := c.Get(url, nil)
	if respErr != nil {
		return "", respErr
	}

	var businessGroups BusinessGroups
	unmarshallErr := utils.UnmarshalJSON(resp.Body, &businessGroups)
	if unmarshallErr != nil {
		return "", unmarshallErr
	}
	// BusinessGroups array will contain only one BusinessGroup element containing the BG
	// with the name businessGroupName.
	// Fetch the id of that
	for _, businessGroup := range businessGroups.Content {
		if businessGroup.Name == businessGroupName {
			log.Info("Found the business group id of the group %s: %s ", businessGroupName, businessGroup.ID)
			return businessGroup.ID, nil
		}
	}
	log.Errorf("No business group found with name: %s ", businessGroupName)
	return "", fmt.Errorf("No business group found with name: %s ", businessGroupName)
}

//GetRequestStatus - To read request status of resource
// which is used to show information to user post create call.
func (c *APIClient) GetRequestStatus(requestID string) (*RequestStatusView, error) {
	//Form a URL to read request status
	path := fmt.Sprintf(ConsumerRequests+"/"+"%s", requestID)
	url := c.BuildEncodedURL(path, nil)
	resp, respErr := c.Get(url, nil)
	if respErr != nil {
		return nil, respErr
	}

	var response RequestStatusView
	unmarshallErr := utils.UnmarshalJSON(resp.Body, &response)
	if unmarshallErr != nil {
		return nil, unmarshallErr
	}
	return &response, nil
}

// GetDeploymentState - Read the state of a vRA7 Deployment
func (c *APIClient) GetDeploymentState(CatalogRequestID string) (*ResourceView, error) {
	//Form an URL to fetch resource list view
	path := fmt.Sprintf(GetRequestResourceViewAPI, CatalogRequestID)
	url := c.BuildEncodedURL(path, nil)
	resp, respErr := c.Get(url, nil)
	if respErr != nil {
		return nil, respErr
	}

	var response ResourceView
	unmarshallErr := utils.UnmarshalJSON(resp.Body, &response)
	if unmarshallErr != nil {
		return nil, unmarshallErr
	}
	return &response, nil
}

// GetRequestResourceView retrieves the resources that were provisioned as a result of a given request.
func (c *APIClient) GetRequestResourceView(catalogRequestID string) (*RequestResourceView, error) {
	path := fmt.Sprintf(GetRequestResourceViewAPI, catalogRequestID)
	url := c.BuildEncodedURL(path, nil)
	resp, respErr := c.Get(url, nil)
	if respErr != nil {
		return nil, respErr
	}

	var response RequestResourceView
	unmarshallErr := utils.UnmarshalJSON(resp.Body, &response)
	if unmarshallErr != nil {
		return nil, unmarshallErr
	}
	return &response, nil
}

//RequestCatalogItem - Make a catalog request.
func (c *APIClient) RequestCatalogItem(requestTemplate *CatalogItemRequestTemplate) (*CatalogRequest, error) {
	//Form a path to set a REST call to create a machine
	path := fmt.Sprintf(EntitledCatalogItems+"/"+"%s"+
		"/requests", requestTemplate.CatalogItemID)

	buffer, _ := utils.MarshalToJSON(requestTemplate)
	url := c.BuildEncodedURL(path, nil)
	resp, respErr := c.Post(url, buffer, nil)
	if respErr != nil {
		return nil, respErr
	}

	var response CatalogRequest
	unmarshallErr := utils.UnmarshalJSON(resp.Body, &response)
	if unmarshallErr != nil {
		return nil, unmarshallErr
	}
	return &response, nil
}

// GetResourceActions get the resource actions allowed for a resource
func (c *APIClient) GetResourceActions(catalogItemRequestID string) (*ResourceActions, error) {
	path := fmt.Sprintf(GetResourceAPI, catalogItemRequestID)

	url := c.BuildEncodedURL(path, nil)
	resp, respErr := c.Get(url, nil)
	if respErr != nil {
		return nil, respErr
	}

	var resourceActions ResourceActions
	unmarshallErr := utils.UnmarshalJSON(resp.Body, &resourceActions)
	if unmarshallErr != nil {
		return nil, unmarshallErr
	}
	return &resourceActions, nil
}

// GetResourceActionTemplate get the action template corresponding to the action id
func (c *APIClient) GetResourceActionTemplate(resourceID, actionID string) (*ResourceActionTemplate, error) {
	getActionTemplatePath := fmt.Sprintf(GetActionTemplateAPI, resourceID, actionID)
	log.Info("Call GET to fetch the reconfigure action template %v ", getActionTemplatePath)
	url := c.BuildEncodedURL(getActionTemplatePath, nil)
	resp, respErr := c.Get(url, nil)
	if respErr != nil {
		return nil, respErr
	}

	var resourceActionTemplate ResourceActionTemplate
	unmarshallErr := utils.UnmarshalJSON(resp.Body, &resourceActionTemplate)
	if unmarshallErr != nil {
		return nil, unmarshallErr
	}
	return &resourceActionTemplate, nil
}

// PostResourceAction updates the resource
func (c *APIClient) PostResourceAction(resourceID, actionID string, resourceActionTemplate *ResourceActionTemplate) error {

	postActionTemplatePath := fmt.Sprintf(PostActionTemplateAPI, resourceID, actionID)
	buffer, _ := utils.MarshalToJSON(resourceActionTemplate)
	url := c.BuildEncodedURL(postActionTemplatePath, nil)
	resp, respErr := c.Post(url, buffer, nil)
	if respErr != nil || resp.StatusCode != 201 {
		return respErr
	}
	return nil
}
