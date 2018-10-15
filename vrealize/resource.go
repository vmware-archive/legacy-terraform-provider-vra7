package vrealize

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/op/go-logging"
	"github.com/vmware/terraform-provider-vra7/utils"
)

var (
	log = logging.MustGetLogger(utils.LOGGER_ID)
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

type BusinessGroups struct {
	Content []BusinessGroup `json:"content,omitempty"`
}

type BusinessGroup struct {
	Name string `json:"name,omitempty"`
	Id   string `json:"id,omitempty"`
}

//CatalogRequest - A structure that captures a vRA catalog request.
type CatalogRequest struct {
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

//set_resource_schema - This function is used to update the catalog item template/blueprint
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
		"businessgroup_name": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"wait_timeout": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  15,
		},
		utils.REQUEST_STATUS: {
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

//Replace the value for a given key in a catalog request template.
func replaceValueInRequestTemplate(templateInterface map[string]interface{}, field string, value interface{}) (map[string]interface{}, bool) {
	var replaced bool
	//Iterate over the map to get field provided as an argument
	for key := range templateInterface {
		//If value type is map then set recursive call which will fiend field in one level down of map interface
		if reflect.ValueOf(templateInterface[key]).Kind() == reflect.Map {
			template, _ := templateInterface[key].(map[string]interface{})
			templateInterface[key], replaced = replaceValueInRequestTemplate(template, field, value)
			if replaced == true {
				return templateInterface, true
			}
		} else if key == field {
			//If value type is not map then compare field name with provided field name
			//If both matches then update field value with provided value
			templateInterface[key] = value
			return templateInterface, true
		}
	}
	//Return updated map interface type
	return templateInterface, replaced
}

//modeled after replaceValueInRequestTemplate, for values being added to template vs updating existing ones
func addValueToRequestTemplate(templateInterface map[string]interface{}, field string, value interface{}) map[string]interface{} {
	//simplest case is adding a simple value. Leaving as a func in case there's a need to do more complicated additions later
	//	templateInterface[data]
	for k, v := range templateInterface {
		if reflect.ValueOf(v).Kind() == reflect.Map && k == "data" {
			template, _ := v.(map[string]interface{})
			v = addValueToRequestTemplate(template, field, value)
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
	// Get client handle
	vRAClient := meta.(*APIClient)

	// If catalog_name and catalog_id both not provided then return an error
	if len(d.Get("catalog_name").(string)) <= 0 && len(d.Get("catalog_id").(string)) <= 0 {
		return fmt.Errorf("Either catalog_name or catalog_id should be present in given configuration")
	}

	// If catalog item name is provided then get catalog item ID using name for further process
	// else if catalog item id is provided then fetch catalog name
	if len(d.Get("catalog_name").(string)) > 0 {
		catalogItemID, returnErr := vRAClient.readCatalogItemIDByName(d.Get("catalog_name").(string))
		log.Info("createResource->catalog_id %v\n", catalogItemID)
		if returnErr != nil {
			return fmt.Errorf("%v", returnErr)
		}
		if catalogItemID == nil {
			return fmt.Errorf("No catalog item found with name %v", d.Get("catalog_name").(string))
		} else if catalogItemID == "" {
			return fmt.Errorf("No catalog item found with name %v", d.Get("catalog_name").(string))
		}
		d.Set("catalog_id", catalogItemID.(string))
	} else if len(d.Get("catalog_id").(string)) > 0 {
		CatalogItemName, nameError := vRAClient.readCatalogItemNameByID(d.Get("catalog_id").(string))
		if nameError != nil {
			return fmt.Errorf("%v", nameError)
		}
		if nameError != nil {
			d.Set("catalog_name", CatalogItemName.(string))
		}
	}
	//Get request template for catalog item.
	requestTemplate, err := vRAClient.GetCatalogItemRequestTemplate(d.Get("catalog_id").(string))
	log.Info("createResource->requestTemplate %v\n", requestTemplate)

	catalogConfiguration, _ := d.Get("catalog_configuration").(map[string]interface{})
	for field1 := range catalogConfiguration {
		requestTemplate.Data[field1] = catalogConfiguration[field1]

	}
	log.Info("createResource->requestTemplate.Data %v\n", requestTemplate.Data)

	business_group_id := strings.TrimSpace(d.Get("businessgroup_id").(string))
	business_group_name := strings.TrimSpace(d.Get("businessgroup_name").(string))

	// get the business group id from name
	var businessGroupIdFromName string
	if len(business_group_name) > 0 {
		businessGroupIdFromName, err = vRAClient.GetBusinessGroupId(business_group_name)
		if err != nil || businessGroupIdFromName == "" {
			return err
		}
	}

	//if both business group name and id are provided but does not belong to the same business group, throw an error
	if len(business_group_name) > 0 && len(business_group_id) > 0 && businessGroupIdFromName != business_group_id {
		log.Error("The business group name %s and id %s does not belong to the same business group. Provide either name or id.", business_group_name, business_group_id)
		return errors.New(fmt.Sprintf("The business group name %s and id %s does not belong to the same business group. Provide either name or id.", business_group_name, business_group_id))
	} else if len(business_group_id) > 0 { // else if both are provided and matches or just id is provided, use id
		log.Info("Setting business group id %s ", business_group_id)
		requestTemplate.BusinessGroupID = business_group_id
	} else if len(business_group_name) > 0 { // else if name is provided, use the id fetched from the name
		log.Info("Setting business group id %s for the group %s ", businessGroupIdFromName, business_group_name)
		requestTemplate.BusinessGroupID = businessGroupIdFromName
	}

	// Get all component names in the blueprint corresponding to the catalog item.
	var componentNameList []string
	for field := range requestTemplate.Data {
		if reflect.ValueOf(requestTemplate.Data[field]).Kind() == reflect.Map {
			componentNameList = append(componentNameList, field)
		}
	}
	log.Info("createResource->key_list %v\n", componentNameList)

	// Arrange component names in descending order of text length.
	// Component names are sorted this way because '.', which is used as a separator, may also occur within
	// component names. In these situations, the longest name match that includes '.'s should win.
	sort.Sort(byLength(componentNameList))

	resourceConfiguration, _ := d.Get("resource_configuration").(map[string]interface{})

	validityErr := checkConfigValidity(requestTemplate, resourceConfiguration)
	if validityErr != nil {
		return validityErr
	}

	//Update request template field values with values from user configuration.
	for configKey, configValue := range resourceConfiguration {
		for _, componentName := range componentNameList {
			// User-supplied resource configuration keys are expected to be of the form:
			//     <component name>.<property name>.
			// Extract the property names and values for each component in the blueprint, and add/update
			// them in the right location in the request template.
			if strings.HasPrefix(configKey, componentName) {
				propertyName := strings.TrimPrefix(configKey, componentName+".")
				if len(propertyName) == 0 {
					return fmt.Errorf(
						"resource_configuration key is not in correct format. Expected %s to start with %s",
						configKey, componentName+".")
				}
				// Function call which changes request template field values with user-supplied values
				requestTemplate.Data[componentName] = updateRequestTemplate(
					requestTemplate.Data[componentName].(map[string]interface{}),
					propertyName,
					configValue)
				break
			}
		}
	}

	//update template with deployment level config
	// limit to description and reasons as other things could get us into trouble
	deploymentConfiguration, _ := d.Get("deployment_configuration").(map[string]interface{})
	for depField, depValue := range deploymentConfiguration {
		fieldstr := fmt.Sprintf("%s", depField)
		switch fieldstr {
		case "description":
			requestTemplate.Description = depValue.(string)
		case "reasons":
			requestTemplate.Reasons = depValue.(string)
		default:
			log.Info("unknown option [%s] with value [%s] ignoring\n", depField, depValue)
		}
	}

	log.Info("Updated template - %v\n", requestTemplate.Data)

	if err != nil {
		return fmt.Errorf("Invalid CatalogItem ID %v", err)
	}

	//Fire off a catalog item request to create a deployment.
	catalogRequest, err := vRAClient.RequestCatalogItem(requestTemplate)

	if err != nil {
		return fmt.Errorf("Resource Machine Request Failed: %v", err)
	}

	//Set request ID
	d.SetId(catalogRequest.ID)
	//Set request status
	d.Set(utils.REQUEST_STATUS, "SUBMITTED")
	return waitForRequestCompletion(d, meta)
}

func updateRequestTemplate(templateInterface map[string]interface{}, field string, value interface{}) map[string]interface{} {
	var replaced bool
	templateInterface, replaced = replaceValueInRequestTemplate(templateInterface, field, value)

	if !replaced {
		templateInterface["data"] = addValueToRequestTemplate(templateInterface["data"].(map[string]interface{}), field, value)
	}
	return templateInterface
}

//Function use - to update centOS 6.3 machine present in state file
//Terraform call - terraform refresh
func readActionLink(resourceSpecificLinks []interface{}, reconfigGetLinkTitleRel string) string {
	var actionLink string
	for _, linkData := range resourceSpecificLinks {
		linkInterface := linkData.(map[string]interface{})
		if linkInterface["rel"] == reconfigGetLinkTitleRel {
			//Get resource reconfiguration template link
			actionLink = linkInterface["href"].(string)
			break
		}
	}
	return actionLink
}

func readVMReconfigActionUrls(GetDeploymentStateData *ResourceView) map[string]interface{} {

	var urlMap map[string]interface{}
	urlMap = map[string]interface{}{}
	const reconfigGetLinkTitleRel = "GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name." +
		"machine.Reconfigure}"
	const reconfigPostLinkTitleRel = "POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name." +
		"machine.Reconfigure}"

	for _, value := range GetDeploymentStateData.Content {
		resourceMap := value.(map[string]interface{})
		if resourceMap["resourceType"] == "Infrastructure.Virtual" {
			resourceSpecificData := resourceMap["data"].(map[string]interface{})
			resourceSpecificLinks := resourceMap["links"].([]interface{})
			componentName := resourceSpecificData["Component"].(string)

			reconfigGetLink := readActionLink(resourceSpecificLinks, reconfigGetLinkTitleRel)
			reconfigPostLink := readActionLink(resourceSpecificLinks, reconfigPostLinkTitleRel)
			urlMap[componentName] = []string{reconfigGetLink, reconfigPostLink}
		}
	}
	return urlMap
}

// Terraform call - terraform apply
// This function updates the state of a vRA 7 Deployment when changes to a Terraform file are applied.
// The update is performed on the Deployment using supported (day-2) actions.
func updateResource(d *schema.ResourceData, meta interface{}) error {
	//Get the ID of the catalog request that was used to provision this Deployment.
	catalogItemRequestID := d.Id()
	//Get client handle
	vRAClient := meta.(*APIClient)

	//Get request template for catalog item.
	requestTemplate, _ := vRAClient.GetCatalogItemRequestTemplate(d.Get("catalog_id").(string))
	log.Info("createResource->requestTemplate %v\n", requestTemplate)

	resourceConfiguration, _ := d.Get("resource_configuration").(map[string]interface{})

	validityErr := checkConfigValidity(requestTemplate, resourceConfiguration)
	if validityErr != nil {
		return validityErr
	}

	//If any change made in resource_configuration.
	if d.HasChange("resource_configuration") {

		//Read resource template
		GetDeploymentStateData, errTemplate := vRAClient.GetDeploymentState(catalogItemRequestID)
		if errTemplate != nil {
			return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
		}

		VMReconfigActionUrls := readVMReconfigActionUrls(GetDeploymentStateData)

		//Iterate over the resources in the deployment
		for _, value := range GetDeploymentStateData.Content {
			resourceMap := value.(map[string]interface{})
			if resourceMap["resourceType"] == "Infrastructure.Virtual" {
				resourceSpecificData := resourceMap["data"].(map[string]interface{})
				//resourceSpecificLinks := resourceMap["links"].([]interface{})
				componentName := resourceSpecificData["Component"].(string)
				resourceActionTemplate := new(ResourceActionTemplate)
				apiError := new(APIError)
				//Get reource child reconfiguration template json
				response, err := vRAClient.HTTPClient.New().Get(VMReconfigActionUrls[componentName].([]string)[0]).
					Receive(resourceActionTemplate, apiError)
				response.Close = true

				if !apiError.isEmpty() {
					return apiError
				}
				if err != nil {
					return err
				}
				configChanged := false
				returnFlag := false
				for configKey := range resourceConfiguration {
					//compare resource list (resource_name) with user configuration fields
					if strings.HasPrefix(configKey, componentName+".") {
						//If user_configuration contains resource_list element
						// then split user configuration key into resource_name and field_name
						nameList := strings.Split(configKey, componentName+".")
						//actionResponseInterface := actionResponse.(map[string]interface{})
						//Function call which changes the template field values with  user values
						//Replace existing values with new values in resource child template
						resourceActionTemplate.Data, returnFlag = replaceValueInRequestTemplate(
							resourceActionTemplate.Data,
							nameList[1],
							resourceConfiguration[configKey])
						if returnFlag == true {
							configChanged = true
						}

					}
					//delete used user configuration
					//delete(resourceConfiguration, configKey)
				}
				//If template value got changed then set post call and update resource child
				if configChanged != false {
					err := postResourceConfig(
						d,
						VMReconfigActionUrls[componentName].([]string)[1],
						resourceActionTemplate,
						meta)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return waitForRequestCompletion(d, meta)
}

func postResourceConfig(d *schema.ResourceData, reconfigPostLink string, resourceActionTemplate *ResourceActionTemplate, meta interface{}) error {
	vRAClient := meta.(*APIClient)
	resourceActionTemplate2 := new(ResourceActionTemplate)
	apiError2 := new(APIError)

	response2, _ := vRAClient.HTTPClient.New().Post(reconfigPostLink).
		BodyJSON(resourceActionTemplate).Receive(resourceActionTemplate2, apiError2)

	if response2.StatusCode != 201 {
		oldData, _ := d.GetChange("resource_configuration")
		d.Set("resource_configuration", oldData)
		return apiError2
	}
	response2.Close = true
	if !apiError2.isEmpty() {
		oldData, _ := d.GetChange("resource_configuration")
		d.Set("resource_configuration", oldData)
		return apiError2
		//panic(d)
	}
	return nil
}

// Terraform call - terraform refresh
// This function retrieves the latest state of a vRA 7 deployment. Terraform updates its state based on
// the information returned by this function.
func readResource(d *schema.ResourceData, meta interface{}) error {
	//Get the ID of the catalog request that was used to provision this Deployment.
	catalogItemRequestID := d.Id()
	//Get client handle
	vRAClient := meta.(*APIClient)
	//Get requested status
	resourceTemplate, errTemplate := vRAClient.GetRequestStatus(catalogItemRequestID)

	//Raise an exception if error occured while fetching request status
	if errTemplate != nil {
		return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
	}
	//Update resource request status in state file
	d.Set(utils.REQUEST_STATUS, resourceTemplate.Phase)
	//If request is failed then set failed message in state file
	if resourceTemplate.Phase == "FAILED" {
		d.Set("failed_message", resourceTemplate.RequestCompletion.CompletionDetails)
	}

	GetDeploymentStateData, errTemplate := vRAClient.GetDeploymentState(catalogItemRequestID)
	if errTemplate != nil {
		return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
	}

	const reconfigGetLinkTitleRel = "GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name." +
		"machine.Reconfigure}"

	var childConfig map[string]interface{}
	childConfig = map[string]interface{}{}

	for _, value := range GetDeploymentStateData.Content {
		resourceMap := value.(map[string]interface{})
		resourceSpecificData := resourceMap["data"].(map[string]interface{})
		resourceSpecificLinks := resourceMap["links"].([]interface{})
		if resourceSpecificData["Component"] != nil {
			componentName := resourceSpecificData["Component"].(string)
			reconfigGetLink := readActionLink(resourceSpecificLinks, reconfigGetLinkTitleRel)

			resourceActionTemplate, err := getResourceConfigTemplate(reconfigGetLink, d, meta)
			if err != nil {
				return err
			}
			childConfig[componentName] = resourceActionTemplate.Data

			// get IP address
			if ipAddress, ok := resourceSpecificData["ip_address"]; ok {
				childConfig[componentName].(map[string]interface{})["ip_address"] = ipAddress.(string)
			}
		}
	}

	resourceConfiguration, _ := d.Get("resource_configuration").(map[string]interface{})
	changed := false

	resourceConfiguration, changed = updateResourceConfigurationMap(resourceConfiguration, childConfig)

	if changed {
		setError := d.Set("resource_configuration", resourceConfiguration)
		if setError != nil {
			return fmt.Errorf(setError.Error())
		}
	}

	return nil
}

func getResourceConfigTemplate(reconfigGetLink string, d *schema.ResourceData, meta interface{}) (*ResourceActionTemplate, error) {
	vRAClient := meta.(*APIClient)
	resourceActionTemplate := new(ResourceActionTemplate)
	apiError := new(APIError)
	//Get reource child reconfiguration template json
	resp, err := vRAClient.HTTPClient.New().Get(reconfigGetLink).Receive(resourceActionTemplate, apiError)
	resp.Close = true
	if !apiError.isEmpty() {
		return nil, apiError
	}
	if err != nil {
		if err.Error() == "invalid character '<' looking for beginning of value" {
			d.Set(utils.REQUEST_STATUS, "IN_PROGRESS")
			return nil, fmt.Errorf("resource is not yet ready to show up")
		}
		return nil, err
	}
	return resourceActionTemplate, nil
}

// updateResourceConfigurationMap - updates the tf resource > resource_configuration type interface with given values
// Input:
// resourceConfiguration map[string]interface{} : tf resource_configuration
// childConfig map[string]interface{} : data of deployment VMs
// Output:
// resourceConfiguration : updated resource_configuration map[string]interface data
// changed : boolean for data got changed or not
func updateResourceConfigurationMap(
	resourceConfiguration map[string]interface{}, vmData map[string]interface{}) (map[string]interface{}, bool) {
	var changed bool
	for configKey1, configValue1 := range resourceConfiguration {
		for configKey2, configValue2 := range vmData {
			if strings.HasPrefix(configKey1, configKey2+".") {
				trimmedKey := strings.TrimPrefix(configKey1, configKey2+".")
				currentValue := configValue1
				updatedValue := getTemplateFieldValue(configValue2.(map[string]interface{}), trimmedKey)

				if updatedValue != nil && updatedValue != currentValue {
					resourceConfiguration[configKey1] = updatedValue
					changed = true
				}
			}
		}
	}
	return resourceConfiguration, changed
}

//getTemplateFieldValue is use to check and return value of argument key
func getTemplateFieldValue(template map[string]interface{}, key string) interface{} {
	for k, v := range template {
		//If value type is map then set recursive call which will fiend field in one level down of map interface
		if reflect.ValueOf(v).Kind() == reflect.Map {
			template, _ := v.(map[string]interface{})
			resp := getTemplateFieldValue(template, key)
			if resp != nil {
				return convertInterfaceToString(resp)
			}
		} else if k == key {
			//If value type is not map then compare field name with provided field name
			//If both matches then update field value with provided value
			return convertInterfaceToString(v)
		}
	}

	return nil
}

func convertInterfaceToString(interfaceData interface{}) string {
	var stringData string
	if reflect.ValueOf(interfaceData).Kind() == reflect.Float64 {
		stringData =
			strconv.FormatFloat(interfaceData.(float64), 'f', 0, 64)
	} else if reflect.ValueOf(interfaceData).Kind() == reflect.Float32 {
		stringData =
			strconv.FormatFloat(interfaceData.(float64), 'f', 0, 32)
	} else if reflect.ValueOf(interfaceData).Kind() == reflect.Int {
		stringData = strconv.FormatInt(interfaceData.(int64), 10)
	} else {
		stringData = interfaceData.(string)
	}
	return stringData
}

//Function use - To delete resources which are created by terraform and present in state file
//Terraform call - terraform destroy
func deleteResource(d *schema.ResourceData, meta interface{}) error {
	//Get requester machine ID from schema.dataresource
	catalogItemRequestID := d.Id()
	//Get client handle
	vRAClient := meta.(*APIClient)

	//Through an error if request ID has no value or empty value
	if len(d.Id()) == 0 {
		return fmt.Errorf("Resource not found")
	}
	//If resource create status is in_progress then skip delete call and through an exception
	if d.Get(utils.REQUEST_STATUS).(string) != "SUCCESSFUL" {
		if d.Get(utils.REQUEST_STATUS).(string) == "FAILED" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Machine cannot be deleted while in-progress state. Please try later")

	}
	//Fetch machine template
	GetDeploymentStateData, errTemplate := vRAClient.GetDeploymentState(catalogItemRequestID)

	if errTemplate != nil {
		return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
	}

	//Set a delete machine template function call.
	//Which will fetch and return the delete machine template from the given template
	DestroyMachineTemplate, resourceTemplate, errDestroyAction := vRAClient.GetDestroyActionTemplate(GetDeploymentStateData)
	if errDestroyAction != nil {
		if errDestroyAction.Error() == "resource is not created or not found" {
			d.SetId("")
			return fmt.Errorf("possibly resource got deleted outside terraform")
		}

		return fmt.Errorf("Destory Machine action template failed to load: %v", errDestroyAction)
	}

	//Set a destroy machine REST call
	_, errDestroyMachine := vRAClient.DestroyMachine(DestroyMachineTemplate, resourceTemplate)
	//Raise an exception if error got while deleting resource
	if errDestroyMachine != nil {
		return fmt.Errorf("Destory Machine machine operation failed: %v", errDestroyMachine)
	}

	waitTimeout := d.Get("wait_timeout").(int) * 60
	sleepFor := 30
	for i := 0; i < waitTimeout/sleepFor; i++ {
		time.Sleep(time.Duration(sleepFor) * time.Second)
		log.Info("Checking to see if resource is deleted.")
		deploymentStateData, err := vRAClient.GetDeploymentState(catalogItemRequestID)
		if err != nil {
			return fmt.Errorf("Resource view failed to load:  %v", err)
		}
		if len(deploymentStateData.Content) == 0 {
			//If resource got deleted then unset the resource ID from state file
			d.SetId("")
			break
		}
	}
	if d.Id() != "" {
		d.SetId("")
		return fmt.Errorf("resource still being deleted after %v minutes", d.Get("wait_timeout"))
	}
	return nil
}

//DestroyMachine - To set resource destroy call
func (vRAClient *APIClient) DestroyMachine(destroyTemplate *ActionTemplate, resourceViewTemplate *ResourceView) (*ActionResponseTemplate, error) {
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
	resp, err := vRAClient.HTTPClient.New().Post(destroyactionURL).
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
func (vRAClient *APIClient) PowerOffMachine(powerOffTemplate *ActionTemplate, resourceViewTemplate *ResourceView) (*ActionResponseTemplate, error) {
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
	response, err := vRAClient.HTTPClient.New().Post(powerOffMachineactionURL).
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
func (vRAClient *APIClient) GetRequestStatus(ResourceID string) (*RequestStatusView, error) {
	//Form a URL to read request status
	path := fmt.Sprintf("catalog-service/api/consumer/requests/%s", ResourceID)
	RequestStatusViewTemplate := new(RequestStatusView)
	apiError := new(APIError)
	//Set a REST call and fetch a resource request status
	_, err := vRAClient.HTTPClient.New().Get(path).Receive(RequestStatusViewTemplate, apiError)
	if err != nil {
		return nil, err
	}
	if !apiError.isEmpty() {
		return nil, apiError
	}
	return RequestStatusViewTemplate, nil
}

// GetDeploymentState - Read the state of a vRA7 Deployment
func (vRAClient *APIClient) GetDeploymentState(CatalogRequestId string) (*ResourceView, error) {
	//Form an URL to fetch resource list view
	path := fmt.Sprintf("catalog-service/api/consumer/requests/%s"+
		"/resourceViews", CatalogRequestId)
	ResourceView := new(ResourceView)
	apiError := new(APIError)
	//Set a REST call to fetch resource view data
	_, err := vRAClient.HTTPClient.New().Get(path).Receive(ResourceView, apiError)
	if err != nil {
		return nil, err
	}
	if !apiError.isEmpty() {
		return nil, apiError
	}
	return ResourceView, nil
}

//RequestCatalogItem - Make a catalog request.
func (vRAClient *APIClient) RequestCatalogItem(requestTemplate *CatalogItemRequestTemplate) (*CatalogRequest, error) {
	//Form a path to set a REST call to create a machine
	path := fmt.Sprintf("/catalog-service/api/consumer/entitledCatalogItems/%s"+
		"/requests", requestTemplate.CatalogItemID)

	catalogRequest := new(CatalogRequest)
	apiError := new(APIError)

	jsonBody, jErr := json.Marshal(requestTemplate)
	if jErr != nil {
		log.Error("Error marshalling request templat as JSON")
		return nil, jErr
	}

	log.Info("JSON Request Info: %s", jsonBody)
	//Set a REST call to create a machine
	_, err := vRAClient.HTTPClient.New().Post(path).BodyJSON(requestTemplate).
		Receive(catalogRequest, apiError)

	if err != nil {
		return nil, err
	}

	if !apiError.isEmpty() {
		return nil, apiError
	}

	return catalogRequest, nil
}

// check if the resource configuration is valid in the terraform config file
func checkConfigValidity(requestTemplate *CatalogItemRequestTemplate, resourceConfiguration map[string]interface{}) error {
	log.Info("Checking if the terraform config file is valid")

	// Get all component names in the blueprint corresponding to the catalog item.
	componentSet := make(map[string]bool)
	for field := range requestTemplate.Data {
		if reflect.ValueOf(requestTemplate.Data[field]).Kind() == reflect.Map {
			componentSet[field] = true
		}
	}
	log.Info("The component name(s) in the blueprint corresponding to the catalog item: %v\n", componentSet)

	var invalidKeys []string
	// check if the keys in the resourceConfiguration map exists in the componentSet
	// if the key in config is machine1.vsphere.custom.location, match every string after each dot
	// until a matching string is found in componentSet.
	// If found, it's a valid key else the component name is invalid
	for k := range resourceConfiguration {
		var key = k
		var isValid bool
		for strings.LastIndex(key, ".") != -1 {
			lastIndex := strings.LastIndex(key, ".")
			key = key[0:lastIndex]
			if _, ok := componentSet[key]; ok {
				log.Info("The component name %s in the terraform config file is valid ", key)
				isValid = true
				break
			}
		}
		if !isValid {
			invalidKeys = append(invalidKeys, k)
		}
	}
	// there are invalid resource config keys in the terraform config file, abort and throw an error
	if len(invalidKeys) > 0 {
		log.Error("The resource_configuration in the config file has invalid component name(s): %v ", strings.Join(invalidKeys, ", "))
		return fmt.Errorf(utils.CONFIG_INVALID_ERROR, strings.Join(invalidKeys, ", "))
	}
	return nil
}

// check the request status on apply and update
func waitForRequestCompletion(d *schema.ResourceData, meta interface{}) error {

	waitTimeout := d.Get("wait_timeout").(int) * 60
	sleepFor := 30
	request_status := ""
	for i := 0; i < waitTimeout/sleepFor; i++ {
		log.Info("Waiting for %d seconds before checking request status.", sleepFor)
		time.Sleep(time.Duration(sleepFor) * time.Second)

		readResource(d, meta)

		request_status = d.Get(utils.REQUEST_STATUS).(string)
		log.Info("Checking to see if resource is created. Status: %s.", request_status)
		if request_status == "SUCCESSFUL" {
			log.Info("Resource creation SUCCESSFUL.")
			return nil
		} else if request_status == "FAILED" {
			//If request is failed during the time then
			//unset resource details from state.
			d.SetId("")
			return fmt.Errorf("Resource creation FAILED.")
		} else if request_status == "IN_PROGRESS" {
			log.Info("Resource creation is still IN PROGRESS.")
		} else {
			log.Info("Resource creation request status: %s.", request_status)
		}
	}

	// The execution has timed out while still IN PROGRESS.
	// The user will need to use 'terraform refresh' at a later point to resolve this.
	return fmt.Errorf("Resource creation has timed out !!")
}

// Retrieve business group id from business group name
func (vRAClient *APIClient) GetBusinessGroupId(businessGroupName string) (string, error) {

	path := "/identity/api/tenants/" + vRAClient.Tenant + "/subtenants?%24filter=name+eq+'" + businessGroupName + "'"
	log.Info("Fetching business group id from name..GET %s ", path)
	BusinessGroups := new(BusinessGroups)
	apiError := new(APIError)
	_, err := vRAClient.HTTPClient.New().Get(path).Receive(BusinessGroups, apiError)
	if err != nil {
		return "", err
	}
	if !apiError.isEmpty() {
		return "", apiError
	}
	// BusinessGroups array will contain only one BusinessGroup element containing the BG
	// with the name businessGroupName.
	// Fetch the id of that BG
	return BusinessGroups.Content[0].Id, nil
}
