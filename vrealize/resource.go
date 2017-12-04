package vrealize

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

//ResourceViewsTemplate - is used to store information
//related to resource template information.
type ResourceViewsTemplate struct {
	Content []struct {
		ResourceID   string `json:"resourceId"`
		RequestState string `json:"requestState"`
		Links        []struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
		} `json:"links"`
	} `json:"content"`
	Links []interface{} `json:"links"`
}

//RequestStatusView - used to store REST response of
//request triggered against any resource.
type RequestStatusView struct {
	RequestCompletion struct {
		RequestCompletionState string `json:"requestCompletionState"`
		CompletionDetails      string `json:"CompletionDetails"`
	} `json:"requestCompletion"`
	Phase string `json:"phase"`
}

//RequestMachineResponse - used to store response of request
//created against machine provision.
type RequestMachineResponse struct {
	ID           string      `json:"id"`
	IconID       string      `json:"iconId"`
	Version      int         `json:"version"`
	State        string      `json:"state"`
	Description  string      `json:"description"`
	Reasons      interface{} `json:"reasons"`
	RequestedFor string      `json:"requestedFor"`
	RequestedBy  string      `json:"requestedBy"`
	Organization struct {
		TenantRef      string `json:"tenantRef"`
		TenantLabel    string `json:"tenantLabel"`
		SubtenantRef   string `json:"subtenantRef"`
		SubtenantLabel string `json:"subtenantLabel"`
	} `json:"organization"`

	RequestorEntitlementID   string                 `json:"requestorEntitlementId"`
	PreApprovalID            string                 `json:"preApprovalId"`
	PostApprovalID           string                 `json:"postApprovalId"`
	DateCreated              time.Time              `json:"dateCreated"`
	LastUpdated              time.Time              `json:"lastUpdated"`
	DateSubmitted            time.Time              `json:"dateSubmitted"`
	DateApproved             time.Time              `json:"dateApproved"`
	DateCompleted            time.Time              `json:"dateCompleted"`
	Quote                    interface{}            `json:"quote"`
	RequestData              map[string]interface{} `json:"requestData"`
	RequestCompletion        string                 `json:"requestCompletion"`
	RetriesRemaining         int                    `json:"retriesRemaining"`
	RequestedItemName        string                 `json:"requestedItemName"`
	RequestedItemDescription string                 `json:"requestedItemDescription"`
	Components               string                 `json:"components"`
	StateName                string                 `json:"stateName"`

	CatalogItemProviderBinding struct {
		BindingID   string `json:"bindingId"`
		ProviderRef struct {
			ID    string `json:"id"`
			Label string `json:"label"`
		} `json:"providerRef"`
	} `json:"catalogItemProviderBinding"`

	Phase           string `json:"phase"`
	ApprovalStatus  string `json:"approvalStatus"`
	ExecutionStatus string `json:"executionStatus"`
	WaitingStatus   string `json:"waitingStatus"`
	CatalogItemRef  struct {
		ID    string `json:"id"`
		Label string `json:"label"`
	} `json:"catalogItemRef"`
}

//ResourceMachine - use to set resource fields
func ResourceMachine() *schema.Resource {
	return &schema.Resource{
		Create: createResource,
		Read:   readResource,
		Update: updateResource,
		Delete: deleteResource,
		Schema: setResourceSchema(),
	}
}

//set_resource_schema - This function is used to update the catalog template/blueprint
//and replace the values with user defined values added in .tf file.
func setResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"catalog_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"catalog_id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"request_status": {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
		"failed_message": {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
			Optional: true,
		},
		"resource_configuration": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
		},
		"catalog_configuration": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
		},
	}
}

//Function use - to create machine
//Terraform call - terraform apply
func changeTemplateValue(templateInterface map[string]interface{}, field string, value interface{}) map[string]interface{} {
	//Iterate over the map to get field provided as an argument
	for i := range templateInterface {
		//If value type is map then set recursive call which will fiend field in one level down of map interface
		if reflect.ValueOf(templateInterface[i]).Kind() == reflect.Map {
			template, _ := templateInterface[i].(map[string]interface{})
			templateInterface[i] = changeTemplateValue(template, field, value)
		} else if i == field {
			//If value type is not map then compare field name with provided field name
			//If both matches then update field value with provided value
			templateInterface[i] = value
		}
	}
	//Return updated map interface type
	return templateInterface
}

//Function use - to set a create resource call
//Terraform call - terraform apply
func createResource(d *schema.ResourceData, meta interface{}) error {
	//Log file handler to generate logs for debugging purpose
	//Get client handle
	client := meta.(*APIClient)

	//If catalog_name and catalog_id both not provided then throw an error
	if len(d.Get("catalog_name").(string)) <= 0 && len(d.Get("catalog_id").(string)) <= 0 {
		return fmt.Errorf("Either catalog_name or catalog_id should be present in given configuration")
	}

	//If catalog name is provided then get catalog ID using name for further process
	//else if catalog id is provided then fetch catalog name
	log.Println("print before block")
	if len(d.Get("catalog_name").(string)) > 0 {
		log.Println("print in block")
		catalogID, returnErr := client.readCatalogIDByName(d.Get("catalog_name").(string))
		log.Printf("createResource->catalog_id %v\n", catalogID)
		if returnErr != nil {
			return fmt.Errorf("%v", returnErr)
		}
		if catalogID == nil {
			return fmt.Errorf("No catalog found with name %v", d.Get("catalog_name").(string))
		} else if catalogID == "" {
			return fmt.Errorf("No catalog found with name %v", d.Get("catalog_name").(string))
		}
		d.Set("catalog_id", catalogID.(string))
	} else if len(d.Get("catalog_id").(string)) > 0 {
		CatalogName, nameError := client.readCatalogNameByID(d.Get("catalog_id").(string))
		if nameError != nil {
			return fmt.Errorf("%v", nameError)
		}
		if nameError != nil {
			d.Set("catalog_name", CatalogName.(string))
		}
	}
	//Get catalog blueprint
	templateCatalogItem, err := client.GetCatalogItem(d.Get("catalog_id").(string))
	log.Printf("createResource->templateCatalogItem %v\n", templateCatalogItem)

	catalogConfiguration, _ := d.Get("catalog_configuration").(map[string]interface{})
	for field1 := range catalogConfiguration {
		if templateCatalogItem.Data[field1] != nil {
			templateCatalogItem.Data[field1] = catalogConfiguration[field1]
		}else{
			return fmt.Errorf(field1+" is not present in catalog configuration")
		}
	}
	log.Printf("createResource->templateCatalogItem.Data %v\n", templateCatalogItem.Data)

	//Get all resource keys from blueprint in array
	var keyList []string
	for field := range templateCatalogItem.Data {
		if reflect.ValueOf(templateCatalogItem.Data[field]).Kind() == reflect.Map {
			keyList = append(keyList, field)
		}
	}
	log.Printf("createResource->key_list %v\n", keyList)

	//Arrange keys in descending order of text length
	for field1 := range keyList {
		for field2 := range keyList {
			if len(keyList[field1]) > len(keyList[field2]) {
				temp := keyList[field1]
				keyList[field1], keyList[field2] = keyList[field2], temp
			}
		}
	}

	//Update template field values with user configuration
	resourceConfiguration, _ := d.Get("resource_configuration").(map[string]interface{})
	for configKey := range resourceConfiguration {
		for dataKey := range keyList {
			//compare resource list (resource_name) with user configuration fields (resource_name+field_name)
			if strings.Contains(configKey, keyList[dataKey]) {
				//If user_configuration contains resource_list element
				// then split user configuration key into resource_name and field_name
				splitedArray := strings.Split(configKey, keyList[dataKey]+".")
				//Function call which changes the template field values with  user values
				templateCatalogItem.Data[keyList[dataKey]] = changeTemplateValue(
					templateCatalogItem.Data[keyList[dataKey]].(map[string]interface{}),
					splitedArray[1],
					resourceConfiguration[configKey])
			}
		}
		//delete used user configuration
		delete(resourceConfiguration, configKey)
	}
	//Log print of template after values updated
	log.Printf("Updated template - %v\n", templateCatalogItem.Data)

	//Through an exception if there is any error while getting catalog template
	if err != nil {
		return fmt.Errorf("Invalid CatalogItem ID %v", err)
	}

	//Set a  create machine function call
	requestMachine, err := client.RequestMachine(templateCatalogItem)

	//Check if error got while create machine call
	//If Error is occured, through an exception with an error message
	if err != nil {
		return fmt.Errorf("Resource Machine Request Failed: %v", err)
	}

	//Set request ID
	d.SetId(requestMachine.ID)
	//Set request status
	d.Set("request_status", "SUBMITTED")
	return nil
}

//Function use - to update centOS 6.3 machine present in state file
//Terraform call - terraform refresh
func updateResource(d *schema.ResourceData, meta interface{}) error {
	log.Println(d)
	return nil
}

//Function use - To read configuration of centOS 6.3 machine present in state file
//Terraform call - terraform refresh
func readResource(d *schema.ResourceData, meta interface{}) error {
	//Get requester machine ID from schema.dataresource
	requestMachineID := d.Id()
	//Get client handle
	client := meta.(*APIClient)
	//Get requested status
	resourceTemplate, errTemplate := client.GetRequestStatus(requestMachineID)

	//Raise an exception if error occured while fetching request status
	if errTemplate != nil {
		return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
	}

	//Update resource request status in state file
	d.Set("request_status", resourceTemplate.Phase)
	//If request is failed then set failed message in state file
	if resourceTemplate.Phase == "FAILED" {
		d.Set("failed_message", resourceTemplate.RequestCompletion.CompletionDetails)
	}
	return nil
}

//Function use - To delete resources which are created by terraform and present in state file
//Terraform call - terraform destroy
func deleteResource(d *schema.ResourceData, meta interface{}) error {
	//Get requester machine ID from schema.dataresource
	requestMachineID := d.Id()
	//Get client handle
	client := meta.(*APIClient)

	//Through an error if request ID has no value or empty value
	if len(d.Id()) == 0 {
		return fmt.Errorf("Resource not found")
	}
	//If resource create status is in_progress then skip delete call and through an exception
	if d.Get("request_status").(string) != "SUCCESSFUL" {
		if d.Get("request_status").(string) == "FAILED" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Machine cannot be deleted while in-progress state. Please try later")

	}
	//Fetch machine template
	templateResources, errTemplate := client.GetResourceViews(requestMachineID)

	if errTemplate != nil {
		return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
	}

	//Set a delete machine template function call.
	//Which will fetch and return the delete machine template from the given template
	DestroyMachineTemplate, resourceTemplate, errDestroyAction := client.GetDestroyActionTemplate(templateResources)
	if errDestroyAction != nil {
		return fmt.Errorf("Destory Machine action template failed to load: %v", errDestroyAction)
	}
	//Set a destroy machine REST call
	_, errDestroyMachine := client.DestroyMachine(DestroyMachineTemplate, resourceTemplate)
	//Raise an exception if error got while deleting resource
	if errDestroyMachine != nil {
		return fmt.Errorf("Destory Machine machine operation failed: %v", errDestroyMachine)
	}
	//If resource got deleted then unset the resource ID from state file
	d.SetId("")
	return nil
}

//DestroyMachine - To set resource destroy call
func (c *APIClient) DestroyMachine(destroyTemplate *ActionTemplate, resourceViewTemplate *ResourceViewsTemplate) (*ActionResponseTemplate, error) {
	//Get a destroy template URL from given resource template
	var destroyactionURL string
	destroyactionURL = getactionURL(resourceViewTemplate, "POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.destroy.name}")
	//Raise an error if any exception raised while fetching delete resource URL
	if len(destroyactionURL) == 0 {
		return nil, fmt.Errorf("Resource is not created or not found")
	}

	actionResponse := new(ActionResponseTemplate)
	apiError := new(APIError)

	//Set a REST call with delete resource request and delete resource template as a data
	resp, err := c.HTTPClient.New().Post(destroyactionURL).
		BodyJSON(destroyTemplate).Receive(actionResponse, apiError)

	if resp.StatusCode != 201 {
		return nil, err
	}

	if !apiError.isEmpty() {
		return nil, apiError
	}

	return actionResponse, nil
}

//PowerOffMachine - To set resource power-off call
func (c *APIClient) PowerOffMachine(powerOffTemplate *ActionTemplate, resourceViewTemplate *ResourceViewsTemplate) (*ActionResponseTemplate, error) {
	//Get power-off resource URL from given template
	var powerOffMachineactionURL string
	powerOffMachineactionURL = getactionURL(resourceViewTemplate, "POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.PowerOff}")
	//Raise an exception if error got while fetching URL
	if len(powerOffMachineactionURL) == 0 {
		return nil, fmt.Errorf("resource is not created or not found")
	}

	actionResponse := new(ActionResponseTemplate)
	apiError := new(APIError)

	//Set a rest call to power-off the resource with resource power-off template as a data
	response, err := c.HTTPClient.New().Post(powerOffMachineactionURL).
		BodyJSON(powerOffTemplate).Receive(actionResponse, apiError)

	response.Close = true
	if response.StatusCode == 201 {
		return actionResponse, nil
	}

	if !apiError.isEmpty() {
		return nil, apiError
	}

	if err != nil {
		return nil, err
	}

	return nil, err
}

//GetRequestStatus - To read request status of resource
// which is used to show information to user post create call.
func (c *APIClient) GetRequestStatus(ResourceID string) (*RequestStatusView, error) {
	//Form a URL to read request status
	path := fmt.Sprintf("catalog-service/api/consumer/requests/%s", ResourceID)
	RequestStatusViewTemplate := new(RequestStatusView)
	apiError := new(APIError)
	//Set a REST call and fetch a resource request status
	_, err := c.HTTPClient.New().Get(path).Receive(RequestStatusViewTemplate, apiError)
	if err != nil {
		return nil, err
	}
	if !apiError.isEmpty() {
		return nil, apiError
	}
	return RequestStatusViewTemplate, nil
}

//GetResourceViews - To read resource configuration
func (c *APIClient) GetResourceViews(ResourceID string) (*ResourceViewsTemplate, error) {
	//Form an URL to fetch resource list view
	path := fmt.Sprintf("catalog-service/api/consumer/requests/%s"+
		"/resourceViews", ResourceID)
	resourceViewsTemplate := new(ResourceViewsTemplate)
	apiError := new(APIError)
	//Set a REST call to fetch resource view data
	_, err := c.HTTPClient.New().Get(path).Receive(resourceViewsTemplate, apiError)
	if err != nil {
		return nil, err
	}
	if !apiError.isEmpty() {
		return nil, apiError
	}
	return resourceViewsTemplate, nil
}

//RequestMachine - To set create resource REST call
func (c *APIClient) RequestMachine(template *CatalogItemTemplate) (*RequestMachineResponse, error) {
	//Form a path to set a REST call to create a machine
	path := fmt.Sprintf("/catalog-service/api/consumer/entitledCatalogItems/%s"+
		"/requests", template.CatalogItemID)

	requestMachineRes := new(RequestMachineResponse)
	apiError := new(APIError)
	//Set a REST call to create a machine
	_, err := c.HTTPClient.New().Post(path).BodyJSON(template).
		Receive(requestMachineRes, apiError)

	if err != nil {
		return nil, err
	}

	if !apiError.isEmpty() {
		return nil, apiError
	}

	return requestMachineRes, nil
}
