package vrealize

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//ResourceActionTemplate - is used to store information
//related to resource action template information.
type ResourceActionTemplate struct {
	Type        string                 `json:"type"`
	ResourceID  string                 `json:"resourceId"`
	ActionID    string                 `json:"actionId"`
	Description string                 `json:"description"`
	Data        map[string]interface{} `json:"data"`
}

//ResourceView - is used to store information
//related to resource template information.
type ResourceView struct {
	Content []interface {
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
		"businessgroup_id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"wait_timeout": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  15,
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
		"deployment_configuration": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
		},
		"resource_configuration": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
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
func changeTemplateValue(templateInterface map[string]interface{}, field string, value interface{}) (map[string]interface{}, bool) {
	var replaced bool
	//Iterate over the map to get field provided as an argument
	for i := range templateInterface {
		//If value type is map then set recursive call which will fiend field in one level down of map interface
		if reflect.ValueOf(templateInterface[i]).Kind() == reflect.Map {
			template, _ := templateInterface[i].(map[string]interface{})
			templateInterface[i], replaced = changeTemplateValue(template, field, value)
		} else if i == field {
			//If value type is not map then compare field name with provided field name
			//If both matches then update field value with provided value
			templateInterface[i] = value
			return templateInterface, true
		}
	}
	//Return updated map interface type
	return templateInterface, replaced
}

//modeled after changeTemplateValue, for values being added to template vs updating existing ones
func addTemplateValue(templateInterface map[string]interface{}, field string, value interface{}) map[string]interface{} {
	//simplest case is adding a simple value. Leaving as a func in case there's a need to do more complicated additions later
	//	templateInterface[data]
	for i := range templateInterface {
		if reflect.ValueOf(templateInterface[i]).Kind() == reflect.Map && i == "data" {
			template, _ := templateInterface[i].(map[string]interface{})
			templateInterface[i] = addTemplateValue(template, field, value)
		} else { //if i == "data" {
			templateInterface[field] = value
		}
	}
	//Return updated map interface type
	return templateInterface
}

// Terraform call - terraform apply
// This function creates a new vRA 7 Deployment using configuration in a user's Terraform file.
// The Deployment is produced by invoking a catalog item that is specified in the configuration.
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
	if len(d.Get("catalog_name").(string)) > 0 {
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
		templateCatalogItem.Data[field1] = catalogConfiguration[field1]

	}
	log.Printf("createResource->templateCatalogItem.Data %v\n", templateCatalogItem.Data)

	if len(d.Get("businessgroup_id").(string)) > 0 {
		templateCatalogItem.BusinessGroupID = d.Get("businessgroup_id").(string)
	}

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

	//array to keep track of resource values that have been used
	var usedConfigKeys []string
	var replaced bool

	//Update template field values with user configuration
	resourceConfiguration, _ := d.Get("resource_configuration").(map[string]interface{})
	for configKey := range resourceConfiguration {
		for dataKey := range keyList {
			//compare resource list (resource_name) with user configuration fields (resource_name+field_name)
			if strings.Contains(configKey, keyList[dataKey]) {
				//If user_configuration contains resource_list element
				// then split user configuration key into resource_name and field_name
				splitedArray := strings.SplitN(configKey, keyList[dataKey]+".", 2)
				if len(splitedArray) != 2 {
					return fmt.Errorf("resource_configuration key is not in correct format. Expected %s to start with %s", configKey, keyList[dataKey]+".")
				}
				//Function call which changes the template field values with  user values
				templateCatalogItem.Data[keyList[dataKey]], replaced = changeTemplateValue(
					templateCatalogItem.Data[keyList[dataKey]].(map[string]interface{}),
					splitedArray[1],
					resourceConfiguration[configKey])
				if replaced {
					usedConfigKeys = append(usedConfigKeys, configKey)
				} else {
					log.Printf("%s was not replaced", configKey)
				}
			}
		}
	}

	//Add remaining keys to template vs updating values
	// first clean out used values
	for usedKey := range usedConfigKeys {
		delete(resourceConfiguration, usedConfigKeys[usedKey])
	}
	log.Println("Entering Add Loop")
	for configKey2 := range resourceConfiguration {
		for dataKey := range keyList {
			log.Printf("Add Loop: configKey2=[%s] keyList[%d] =[%v]", configKey2, dataKey, keyList[dataKey])
			if strings.Contains(configKey2, keyList[dataKey]) {
				splitArray := strings.Split(configKey2, keyList[dataKey]+".")
				log.Printf("Add Loop Contains %+v", splitArray[1])
				resourceItem := templateCatalogItem.Data[keyList[dataKey]].(map[string]interface{})
				resourceItem = addTemplateValue(
					resourceItem["data"].(map[string]interface{}),
					splitArray[1],
					resourceConfiguration[configKey2])
			}
		}
	}
	//update template with deployment level config
	// limit to description and reasons as other things could get us into trouble
	deploymentConfiguration, _ := d.Get("deployment_configuration").(map[string]interface{})
	for depField := range deploymentConfiguration {
		fieldstr := fmt.Sprintf("%s", depField)
		switch fieldstr {
		case "description":
			templateCatalogItem.Description = deploymentConfiguration[depField].(string)
		case "reasons":
			templateCatalogItem.Reasons = deploymentConfiguration[depField].(string)
		default:
			log.Printf("unknown option [%s] with value [%s] ignoring\n", depField, deploymentConfiguration[depField])
		}
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

	waitTimeout := d.Get("wait_timeout").(int) * 60

	for i := 0; i < waitTimeout/30; i++ {
		time.Sleep(3e+10)
		readResource(d, meta)

		if d.Get("request_status") == "SUCCESSFUL" {
			return nil
		}
		if d.Get("request_status") == "FAILED" {
			//If request is failed during the time then
			//unset resource details from state.
			d.SetId("")
			return fmt.Errorf("instance got failed while creating." +
				" kindly check detail for more information")
		}
	}
	if d.Get("request_status") == "IN_PROGRESS" {
		//If request is in_progress state during the time then
		//keep resource details in state files and throw an error
		//so that the child resource won't go for create call.
		//If execution gets timed-out and status is in progress
		//then dependent machine won't be get created in this iteration.
		//A user needs to ensure that the status should be a success state
		//using terraform refresh command and hit terraform apply again.
		return fmt.Errorf("resource is still being created")
	}

	return nil
}

// Terraform call - terraform apply
// This function updates the state of a vRA 7 Deployment when changes to a Terraform file are applied.
// The update is performed on the Deployment using supported (day-2) actions.
func updateResource(d *schema.ResourceData, meta interface{}) error {
	//Get requester machine ID from schema.dataresource
	catalogItemRequestID := d.Id()
	//Get client handle
	client := meta.(*APIClient)

	//If any change made in resource_configuration.
	if d.HasChange("resource_configuration") {
		//Read resource template
		resourceView, errTemplate := client.GetResourceViews(catalogItemRequestID)
		if errTemplate != nil {
			return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
		}
		//Set URL relation variables.
		const reconfigGetLinkTitleRel = "GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name." +
			"machine.Reconfigure}"
		const reconfigPostLinkTitleRel = "POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name." +
			"machine.Reconfigure}"

		resourceConfiguration, _ := d.Get("resource_configuration").(map[string]interface{})

		//Iterate over the resources in the deployment
		for item := range resourceView.Content {
			resourceMap := resourceView.Content[item].(map[string]interface{})
			//If component is not empty
			if resourceMap["Component"] != "Infrastructure.Virtual" {
				resourceSpecificData := resourceMap["data"].(map[string]interface{})
				if resourceSpecificData["Component"] == nil {
					continue
				}
				resourceSpecificLinks := resourceMap["links"].([]interface{})
				componentName := resourceSpecificData["Component"].(string)
				var reconfigGetLink string
				var reconfigPostLink string
				//Iterate over the links present in content elements
				for link := range resourceSpecificLinks {
					linkInterface := resourceSpecificLinks[link].(map[string]interface{})
					if linkInterface["rel"] == reconfigGetLinkTitleRel {
						//Get resource reconfiguration template link
						reconfigGetLink = linkInterface["href"].(string)
					}
					if linkInterface["rel"] == reconfigPostLinkTitleRel {
						//Get resourcce reconfiguration request post link
						reconfigPostLink = linkInterface["href"].(string)
					}
					if reconfigPostLink != "" && reconfigGetLink != "" {
						//If both links are found then get out of the loop
						break
					}
				}

				resourceAction := new(ResourceActionTemplate)
				apiError := new(APIError)
				//Get reource child reconfiguration template json
				response, err := client.HTTPClient.New().Get(reconfigGetLink).
					Receive(resourceAction, apiError)
				response.Close = true

				if !apiError.isEmpty() {
					return apiError
				}
				if err != nil {
					return err
				}
				for configKey := range resourceConfiguration {
					//compare resource list (resource_name) with user configuration fields
					if strings.HasPrefix(configKey, componentName+".") {
						//If user_configuration contains resource_list element
						// then split user configuration key into resource_name and field_name
						splitedArray := strings.Split(configKey, componentName+".")
						//actionResponseInterface := actionResponse.(map[string]interface{})
						//Function call which changes the template field values with  user values
						//Replace existing values with new values in resource child template
						resourceAction.Data, _ = changeTemplateValue(
							resourceAction.Data,
							splitedArray[1],
							resourceConfiguration[configKey])

					}
					//delete used user configuration
					//delete(resourceConfiguration, configKey)
				}
				//If template value got changed then set post call and update resource child
				resourceAction2 := new(ResourceActionTemplate)
				apiError2 := new(APIError)
				response2, _ := client.HTTPClient.New().Post(reconfigPostLink).
					BodyJSON(resourceAction).Receive(resourceAction2, apiError2)

				if response2.StatusCode != 201 {
					oldData, _ := d.GetChange("resource_configuration")
					d.Set("resource_configuration", oldData)
					return apiError2
				}
				response2.Close = true
				if !apiError2.isEmpty() {
					oldData, _ := d.GetChange("resource_configuration")
					d.Set("resource_configuration", oldData)
					panic(d)
				}
			}
		}
	}
	return nil
}

// Terraform call - terraform refresh
// This function retrieves the latest state of a vRA 7 deployment. Terraform updates its state based on
// the information returned by this function.
func readResource(d *schema.ResourceData, meta interface{}) error {
	//Get requester machine ID from schema.dataresource
	catalogItemRequestID := d.Id()
	//Get client handle
	client := meta.(*APIClient)
	//Get requested status
	resourceTemplate, errTemplate := client.GetRequestStatus(catalogItemRequestID)

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

	resourceView, errTemplate := client.GetResourceViews(catalogItemRequestID)
	if errTemplate != nil {
		return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
	}

	const reconfigGetLinkTitleRel = "GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name." +
		"machine.Reconfigure}"

	var childConfig map[string]interface{}
	childConfig = map[string]interface{}{}

	for item := range resourceView.Content {
		resourceMap := resourceView.Content[item].(map[string]interface{})
		resourceSpecificData := resourceMap["data"].(map[string]interface{})
		resourceSpecificLinks := resourceMap["links"].([]interface{})
		if resourceSpecificData["Component"] != nil {
			componentName := resourceSpecificData["Component"].(string)
			var reconfigGetLink string
			for link := range resourceSpecificLinks {
				linkInterface := resourceSpecificLinks[link].(map[string]interface{})
				if linkInterface["rel"] == reconfigGetLinkTitleRel {
					//Get resource reconfiguration template link
					reconfigGetLink = linkInterface["href"].(string)
					break
				}
			}
			resourceAction := new(ResourceActionTemplate)
			apiError := new(APIError)
			//Get reource child reconfiguration template json
			resp, err := client.HTTPClient.New().Get(reconfigGetLink).Receive(resourceAction, apiError)
			resp.Close = true
			if !apiError.isEmpty() {
				return apiError
			}
			if err != nil {
				return err
			}
			childConfig[componentName] = resourceAction.Data
		}
	}

	resourceConfiguration, _ := d.Get("resource_configuration").(map[string]interface{})
	changed := false

	resourceConfiguration, changed = updateResourceConfigurationMap(resourceConfiguration, childConfig, changed)

	if changed {
		setError := d.Set("resource_configuration", resourceConfiguration)
		if setError != nil {
			return fmt.Errorf(setError.Error())
		}
	}

	return nil
}

//updateResourceConfigurationMap updates the resource_configuration schema map data
//and returns boolean value for whether map got updated or not
func updateResourceConfigurationMap(resourceConfiguration map[string]interface{}, childConfig map[string]interface{}, changed bool) (map[string]interface{}, bool) {
	for index1 := range resourceConfiguration {
		for index2 := range childConfig {
			if strings.Contains(index1, index2+".") {
				splitedArray := strings.Split(index1, index2+".")
				currentValue := resourceConfiguration[index1]
				updatedValue := getTemplateFieldValue(childConfig[index2].(map[string]interface{}), splitedArray[1])
				if updatedValue != currentValue {
					if reflect.ValueOf(updatedValue).Kind() == reflect.Float64 {
						resourceConfiguration[index1] = strconv.FormatFloat(updatedValue.(float64), 'f', 0, 64)
					} else if reflect.ValueOf(updatedValue).Kind() == reflect.Float32 {
						resourceConfiguration[index1] = strconv.FormatFloat(updatedValue.(float64), 'f', 0, 32)
					} else if reflect.ValueOf(updatedValue).Kind() == reflect.Int {
						resourceConfiguration[index1] = strconv.FormatInt(updatedValue.(int64), 10)
					} else {
						resourceConfiguration[index1] = updatedValue
					}
					changed = true
				}
			}
		}
	}
	return resourceConfiguration, changed
}

//getTemplateFieldValue is use to check and return value of argument key
func getTemplateFieldValue(template map[string]interface{}, key string) interface{} {
	for i := range template {
		//If value type is map then set recursive call which will fiend field in one level down of map interface
		if reflect.ValueOf(template[i]).Kind() == reflect.Map {
			template, _ := template[i].(map[string]interface{})
			resp := getTemplateFieldValue(template, key)
			if resp != nil {
				return resp
			}
		} else if i == key {
			//If value type is not map then compare field name with provided field name
			//If both matches then update field value with provided value
			return template[i]
		}
	}

	return nil
}

//Function use - To delete resources which are created by terraform and present in state file
//Terraform call - terraform destroy
func deleteResource(d *schema.ResourceData, meta interface{}) error {
	//Get requester machine ID from schema.dataresource
	catalogItemRequestID := d.Id()
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
	resourceView, errTemplate := client.GetResourceViews(catalogItemRequestID)

	if errTemplate != nil {
		return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
	}

	//Set a delete machine template function call.
	//Which will fetch and return the delete machine template from the given template
	DestroyMachineTemplate, resourceTemplate, errDestroyAction := client.GetDestroyActionTemplate(resourceView)
	if errDestroyAction != nil {
		if errDestroyAction.Error() == "resource is not created or not found" {
			d.SetId("")
			return fmt.Errorf("possibly resource got deleted outside terraform")
		}

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
func (c *APIClient) DestroyMachine(destroyTemplate *ActionTemplate, resourceViewTemplate *ResourceView) (*ActionResponseTemplate, error) {
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
func (c *APIClient) PowerOffMachine(powerOffTemplate *ActionTemplate, resourceViewTemplate *ResourceView) (*ActionResponseTemplate, error) {
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
func (c *APIClient) GetResourceViews(ResourceID string) (*ResourceView, error) {
	//Form an URL to fetch resource list view
	path := fmt.Sprintf("catalog-service/api/consumer/requests/%s"+
		"/resourceViews", ResourceID)
	ResourceView := new(ResourceView)
	apiError := new(APIError)
	//Set a REST call to fetch resource view data
	_, err := c.HTTPClient.New().Get(path).Receive(ResourceView, apiError)
	if err != nil {
		return nil, err
	}
	if !apiError.isEmpty() {
		return nil, apiError
	}
	return ResourceView, nil
}

//RequestMachine - To set create resource REST call
func (c *APIClient) RequestMachine(template *CatalogItemTemplate) (*RequestMachineResponse, error) {
	//Form a path to set a REST call to create a machine
	path := fmt.Sprintf("/catalog-service/api/consumer/entitledCatalogItems/%s"+
		"/requests", template.CatalogItemID)

	requestMachineRes := new(RequestMachineResponse)
	apiError := new(APIError)

	jsonBody, jErr := json.Marshal(template)
	if jErr != nil {
		log.Printf("Error marshalling template as JSON")
		return nil, jErr
	}

	log.Printf("JSON Request Info: %s", jsonBody)
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
