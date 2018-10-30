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
	log                     = logging.MustGetLogger(utils.LOGGER_ID)
	catalogItemName         string
	catalogItemID           string
	businessGroupId         string
	businessGroupName       string
	waitTimeout             int
	requestStatus           string
	failedMessage           string
	deploymentConfiguration map[string]interface{}
	resourceConfiguration   map[string]interface{}
	catalogConfiguration    map[string]interface{}
)

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
	readProviderConfiguration(d)

	requestTemplate, validityErr := checkConfigValuesValidity(vRAClient, d)
	if validityErr != nil {
		return validityErr
	}
	validityErr = checkResourceConfigValidity(requestTemplate)
	if validityErr != nil {
		return validityErr
	}

	for field1 := range catalogConfiguration {
		requestTemplate.Data[field1] = catalogConfiguration[field1]
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

// Terraform call - terraform apply
// This function updates the state of a vRA 7 Deployment when changes to a Terraform file are applied.
// The update is performed on the Deployment using supported (day-2) actions.
func updateResource(d *schema.ResourceData, meta interface{}) error {

	// Get the ID of the catalog request that was used to provision this Deployment.
	catalogItemRequestID := d.Id()
	// Get client handle
	vRAClient := meta.(*APIClient)
	readProviderConfiguration(d)

	requestTemplate, validityErr := checkConfigValuesValidity(vRAClient, d)
	if validityErr != nil {
		return validityErr
	}
	validityErr = checkResourceConfigValidity(requestTemplate)
	if validityErr != nil {
		return validityErr
	}

	ResourceActions := new(ResourceActions)
	apiError := new(APIError)

	path := fmt.Sprintf(utils.GET_RESOURCE_API, catalogItemRequestID)
	_, err := vRAClient.HTTPClient.New().Get(path).
		Receive(ResourceActions, apiError)

	if err != nil {
		log.Errorf("Error while reading resource actions for the request %v: %v ", catalogItemRequestID, err.Error())
		return fmt.Errorf("Error while reading resource actions for the request %v: %v  ", catalogItemRequestID, err.Error())
	}
	if apiError != nil && !apiError.isEmpty() {
		log.Errorf("Error while reading resource actions for the request %v: %v ", catalogItemRequestID, apiError.Errors)
		return fmt.Errorf("Error while reading resource actions for the request %v: %v  ", catalogItemRequestID, apiError.Errors)
	}

	// If any change made in resource_configuration.
	if d.HasChange(utils.RESOURCE_CONFIGURATION) {
		for _, resources := range ResourceActions.Content {
			if resources.ResourceTypeRef.Id == utils.INFRASTRUCTURE_VIRTUAL {
				var reconfigureEnabled bool
				var reconfigureActionId string
				for _, op := range resources.Operations {
					if op.Name == utils.RECONFIGURE {
						reconfigureEnabled = true
						reconfigureActionId = op.OperationId
						break
					}
				}
				// if reconfigure action is not available for any resource of the deployment
				// return with an error message
				if !reconfigureEnabled {
					return fmt.Errorf("Update is not allowed for resource %v, your entitlement has no Reconfigure action enabled", resources.Id)
				} else {
					resourceData := resources.ResourceData
					entries := resourceData.Entries
					for _, entry := range entries {
						if entry.Key == utils.COMPONENT {
							entryValue := entry.Value
							var componentName string
							for k, v := range entryValue {
								if k == "value" {
									componentName = v.(string)
								}
							}
							resourceActionTemplate := new(ResourceActionTemplate)
							apiError := new(APIError)
							log.Info("Retrieving reconfigure action template for the component: %v ", componentName)
							getActionTemplatePath := fmt.Sprintf(utils.GET_ACTION_TEMPLATE_API, resources.Id, reconfigureActionId)
							log.Info("Call GET to fetch the reconfigure action template %v ", getActionTemplatePath)
							response, err := vRAClient.HTTPClient.New().Get(getActionTemplatePath).
								Receive(resourceActionTemplate, apiError)
							response.Close = true
							if !apiError.isEmpty() {
								log.Errorf("Error retrieving reconfigure action template for the component %v: %v ", componentName, apiError.Error())
								return fmt.Errorf("Error retrieving reconfigure action template for the component %v: %v ", componentName, apiError.Error())
							}
							if err != nil {
								log.Errorf("Error retrieving reconfigure action template for the component %v: %v ", componentName, err.Error())
								return fmt.Errorf("Error retrieving reconfigure action template for the component %v: %v ", componentName, err.Error())
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
							}
							// If template value got changed then set post call and update resource child
							if configChanged != false {
								postActionTemplatePath := fmt.Sprintf(utils.POST_ACTION_TEMPLATE_API, resources.Id, reconfigureActionId)
								err := postResourceConfig(d, postActionTemplatePath, resourceActionTemplate, meta)
								if err != nil {
									return err
								}
							}
						}
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
		oldData, _ := d.GetChange(utils.RESOURCE_CONFIGURATION)
		d.Set(utils.RESOURCE_CONFIGURATION, oldData)
		return apiError2
	}
	response2.Close = true
	if !apiError2.isEmpty() {
		oldData, _ := d.GetChange(utils.RESOURCE_CONFIGURATION)
		d.Set(utils.RESOURCE_CONFIGURATION, oldData)
		return apiError2
	}
	return nil
}

// Terraform call - terraform refresh
// This function retrieves the latest state of a vRA 7 deployment. Terraform updates its state based on
// the information returned by this function.
func readResource(d *schema.ResourceData, meta interface{}) error {
	// Get the ID of the catalog request that was used to provision this Deployment.
	catalogItemRequestID := d.Id()

	log.Info("Calling read resource to get the current resource status of the request: %v ", catalogItemRequestID)
	// Get client handle
	vRAClient := meta.(*APIClient)
	// Get requested status
	resourceTemplate, errTemplate := vRAClient.GetRequestStatus(catalogItemRequestID)

	if errTemplate != nil {
		log.Errorf("Resource view failed to load:  %v", errTemplate)
		return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
	}

	// Update resource request status in state file
	d.Set(utils.REQUEST_STATUS, resourceTemplate.Phase)
	// If request is failed then set failed message in state file
	if resourceTemplate.Phase == "FAILED" {
		log.Errorf(resourceTemplate.RequestCompletion.CompletionDetails)
		d.Set(utils.FAILED_MESSAGE, resourceTemplate.RequestCompletion.CompletionDetails)
	}

	requestResourceView, errTemplate := vRAClient.GetRequestResourceView(catalogItemRequestID)
	if errTemplate != nil {
		return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
	}

	resourceDataMap := make(map[string]map[string]interface{})
	for _, resource := range requestResourceView.Content {
		if resource.ResourceType == utils.INFRASTRUCTURE_VIRTUAL {
			resourceData := resource.ResourcesData
			log.Info("The resource data map of the resource %v is: \n%v", resourceData.Component, resource.ResourcesData)
			dataVals := make(map[string]interface{})
			resourceDataMap[resourceData.Component] = dataVals
			dataVals[utils.MACHINE_CPU] = resourceData.Cpu
			dataVals[utils.MACHINE_STORAGE] = resourceData.Storage
			dataVals[utils.IP_ADDRESS] = resourceData.IpAddress
			dataVals[utils.MACHINE_MEMORY] = resourceData.Memory
			dataVals[utils.MACHINE_NAME] = resourceData.MachineName
			dataVals[utils.MACHINE_GUEST_OS] = resourceData.MachineGuestOperatingSystem
			dataVals[utils.MACHINE_BP_NAME] = resourceData.MachineBlueprintName
			dataVals[utils.MACHINE_TYPE] = resourceData.MachineType
			dataVals[utils.MACHINE_RESERVATION_NAME] = resourceData.MachineReservationName
			dataVals[utils.MACHINE_INTERFACE_TYPE] = resourceData.MachineInterfaceType
			dataVals[utils.MACHINE_ID] = resourceData.MachineId
			dataVals[utils.MACHINE_GROUP_NAME] = resourceData.MachineGroupName
			dataVals[utils.MACHINE_DESTRUCTION_DATE] = resourceData.MachineDestructionDate
		}
	}
	resourceConfiguration, _ := d.Get(utils.RESOURCE_CONFIGURATION).(map[string]interface{})
	changed := false

	resourceConfiguration, changed = updateResourceConfigurationMap(resourceConfiguration, resourceDataMap)

	if changed {
		setError := d.Set(utils.RESOURCE_CONFIGURATION, resourceConfiguration)
		if setError != nil {
			return fmt.Errorf(setError.Error())
		}
	}
	return nil
}

// update the resource configuration with the deployment resource data.
// if there is difference between the config data and deployment data, return true
func updateResourceConfigurationMap(
	resourceConfiguration map[string]interface{}, vmData map[string]map[string]interface{}) (map[string]interface{}, bool) {
	log.Info("Updating resource configuration with the request resource view data...")
	var changed bool
	for configKey1, configValue1 := range resourceConfiguration {
		for configKey2, configValue2 := range vmData {
			if strings.HasPrefix(configKey1, configKey2+".") {
				trimmedKey := strings.TrimPrefix(configKey1, configKey2+".")
				currentValue := configValue1
				updatedValue := convertInterfaceToString(configValue2[trimmedKey])

				if updatedValue != "" && updatedValue != currentValue {
					resourceConfiguration[configKey1] = updatedValue
					changed = true
				}
			}
		}
	}
	return resourceConfiguration, changed
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
		stringData = strconv.Itoa(interfaceData.(int))
	} else if reflect.ValueOf(interfaceData).Kind() == reflect.String {
		stringData = interfaceData.(string)
	} else if reflect.ValueOf(interfaceData).Kind() == reflect.Bool {
		stringData = strconv.FormatBool(interfaceData.(bool))
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

	log.Info("Calling delete resource for the request id %v ", catalogItemRequestID)
	ResourceActions := new(ResourceActions)
	apiError := new(APIError)

	_, err := vRAClient.HTTPClient.New().Get("/catalog-service/api/consumer/resources/"+d.Id()+"/actions").
		Receive(ResourceActions, apiError)

	if err != nil {
		log.Errorf("error while reading resource actions for the request %v ", d.Id())
	}
	if apiError != nil {
		log.Errorf("API error while reading resource actions for the request %v ", d.Id())
	}

	log.Info("the resource action struct is %v", ResourceActions)

	//If resource create status is in_progress then skip delete call and through an exception
	if d.Get(utils.REQUEST_STATUS).(string) != "SUCCESSFUL" {
		if d.Get(utils.REQUEST_STATUS).(string) == "FAILED" {
			log.Info("The status of the request is FAILED, setting the id to null")
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
			log.Errorf("The destroy action template cannot be retrieved with error, resource is not created or not found")
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

	waitTimeout := d.Get(utils.WAIT_TIME_OUT).(int) * 60
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
		return fmt.Errorf("resource still being deleted after %v minutes", d.Get(utils.WAIT_TIME_OUT))
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
func (vRAClient *APIClient) GetRequestStatus(requestId string) (*RequestStatusView, error) {
	//Form a URL to read request status
	path := fmt.Sprintf("catalog-service/api/consumer/requests/%s", requestId)
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

// Retrieves the resources that were provisioned as a result of a given request.
func (vRAClient *APIClient) GetRequestResourceView(catalogRequestId string) (*RequestResourceView, error) {
	path := fmt.Sprintf(utils.GET_REQUEST_RESOURCE_VIEW_API, catalogRequestId)
	requestResourceView := new(RequestResourceView)
	apiError := new(APIError)
	_, err := vRAClient.HTTPClient.New().Get(path).Receive(requestResourceView, apiError)
	if err != nil {
		return nil, err
	}
	if !apiError.isEmpty() {
		return nil, apiError
	}
	return requestResourceView, nil
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
func checkResourceConfigValidity(requestTemplate *CatalogItemRequestTemplate) error {
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

// check if the values provided in the config file are valid and set
// them in the resource schema. Requires to call APIs
func checkConfigValuesValidity(vRAClient *APIClient, d *schema.ResourceData) (*CatalogItemRequestTemplate, error) {
	// 	// If catalog_name and catalog_id both not provided then return an error
	if len(catalogItemName) <= 0 && len(catalogItemID) <= 0 {
		return nil, fmt.Errorf("Either catalog_name or catalog_id should be present in given configuration")
	}

	var catalogItemIdFromName string
	var catalogItemNameFromId string
	var err error
	// if catalog item id is provided, fetch the catalog item name
	if len(catalogItemName) > 0 {
		catalogItemIdFromName, err = vRAClient.readCatalogItemIDByName(catalogItemName)
		if err != nil || catalogItemIdFromName == "" {
			return nil, fmt.Errorf("Error in finding catalog item id corresponding to the catlog item name %v: \n %v", catalogItemName, err)
		}
		log.Info("The catalog item id provided in the config is %v\n", catalogItemIdFromName)
	}

	// if catalog item name is provided, fetch the catalog item id
	if len(catalogItemID) > 0 { // else if both are provided and matches or just id is provided, use id
		catalogItemNameFromId, err = vRAClient.readCatalogItemNameByID(catalogItemID)
		if err != nil || catalogItemNameFromId == "" {
			return nil, fmt.Errorf("Error in finding catalog item name corresponding to the catlog item id %v: \n %v", catalogItemID, err)
		}
		log.Info("The catalog item name corresponding to the catalog item id in the config is:  %v\n", catalogItemNameFromId)
	}

	// if both catalog item name and id are provided but does not belong to the same catalog item, throw an error
	if len(catalogItemName) > 0 && len(catalogItemID) > 0 && (catalogItemIdFromName != catalogItemID || catalogItemNameFromId != catalogItemName) {
		log.Error("The catalog item name %s and id %s does not belong to the same catalog item. Provide either name or id.")
		return nil, errors.New("The catalog item name %s and id %s does not belong to the same catalog item. Provide either name or id.")
	} else if len(catalogItemID) > 0 { // else if both are provided and matches or just id is provided, use id
		d.Set(utils.CATALOG_ID, catalogItemID)
		d.Set(utils.CATALOG_NAME, catalogItemNameFromId)
	} else if len(catalogItemName) > 0 { // else if name is provided, use the id fetched from the name
		d.Set(utils.CATALOG_ID, catalogItemIdFromName)
		d.Set(utils.CATALOG_NAME, catalogItemName)
	}

	// update the catalogItemID var with the updated id
	catalogItemID = d.Get(utils.CATALOG_ID).(string)

	// Get request template for catalog item.
	requestTemplate, err := vRAClient.GetCatalogItemRequestTemplate(catalogItemID)
	if err != nil {
		return nil, err
	}
	log.Info("The request template data corresponding to the catalog item %v is: \n %v\n", catalogItemID, requestTemplate.Data)

	for field1 := range catalogConfiguration {
		requestTemplate.Data[field1] = catalogConfiguration[field1]

	}
	// get the business group id from name
	var businessGroupIdFromName string
	if len(businessGroupName) > 0 {
		businessGroupIdFromName, err := vRAClient.GetBusinessGroupId(businessGroupName)
		if err != nil || businessGroupIdFromName == "" {
			return nil, err
		}
	}

	//if both business group name and id are provided but does not belong to the same business group, throw an error
	if len(businessGroupName) > 0 && len(businessGroupId) > 0 && businessGroupIdFromName != businessGroupId {
		log.Error("The business group name %s and id %s does not belong to the same business group. Provide either name or id.", businessGroupName, businessGroupId)
		return nil, errors.New(fmt.Sprintf("The business group name %s and id %s does not belong to the same business group. Provide either name or id.", businessGroupName, businessGroupId))
	} else if len(businessGroupId) > 0 { // else if both are provided and matches or just id is provided, use id
		log.Info("Setting business group id %s ", businessGroupId)
		requestTemplate.BusinessGroupID = businessGroupId
	} else if len(businessGroupName) > 0 { // else if name is provided, use the id fetched from the name
		log.Info("Setting business group id %s for the group %s ", businessGroupIdFromName, businessGroupName)
		requestTemplate.BusinessGroupID = businessGroupIdFromName
	}
	return requestTemplate, nil
}

// check the request status on apply and update
func waitForRequestCompletion(d *schema.ResourceData, meta interface{}) error {

	waitTimeout := d.Get(utils.WAIT_TIME_OUT).(int) * 60
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
			return fmt.Errorf("Request failed \n %v ", d.Get(utils.FAILED_MESSAGE))
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

// read the config file
func readProviderConfiguration(d *schema.ResourceData) {

	log.Info("Reading the provider configuration data.....")
	catalogItemName = strings.TrimSpace(d.Get(utils.CATALOG_NAME).(string))
	log.Info("Catalog item name: %v ", catalogItemName)
	catalogItemID = strings.TrimSpace(d.Get(utils.CATALOG_ID).(string))
	log.Info("Catalog item ID: %v", catalogItemID)
	businessGroupName = strings.TrimSpace(d.Get(utils.BUSINESS_GROUP_NAME).(string))
	log.Info("Business Group name: %v", businessGroupName)
	businessGroupId = strings.TrimSpace(d.Get(utils.BUSINESS_GROUP_ID).(string))
	log.Info("Business Group Id: %v", businessGroupId)
	waitTimeout = d.Get(utils.WAIT_TIME_OUT).(int) * 60
	log.Info("Wait time out: %v ", waitTimeout)
	failedMessage = strings.TrimSpace(d.Get(utils.FAILED_MESSAGE).(string))
	log.Info("Failed message: %v ", failedMessage)
	resourceConfiguration = d.Get(utils.RESOURCE_CONFIGURATION).(map[string]interface{})
	log.Info("Resource Configuration: %v ", resourceConfiguration)
	deploymentConfiguration = d.Get(utils.DEPLOYMENT_CONFIGURATION).(map[string]interface{})
	log.Info("Deployment Configuration: %v ", deploymentConfiguration)
	catalogConfiguration = d.Get(utils.CATALOG_CONFIGURATION).(map[string]interface{})
	log.Info("Catalog configuration: %v ", catalogConfiguration)
}
