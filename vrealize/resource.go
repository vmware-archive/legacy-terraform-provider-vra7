package vrealize

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	logging "github.com/op/go-logging"
	"github.com/vmware/terraform-provider-vra7/client"
	"github.com/vmware/terraform-provider-vra7/sdk"
	"github.com/vmware/terraform-provider-vra7/utils"
)

var (
	log                     = logging.MustGetLogger(utils.LoggerID)
	catalogItemName         string
	catalogItemID           string
	businessGroupID         string
	businessGroupName       string
	waitTimeout             int
	requestStatus           string
	failedMessage           string
	deploymentConfiguration map[string]interface{}
	resourceConfiguration   map[string]interface{}
	catalogConfiguration    map[string]interface{}
	vraClient               *client.APIClient
)

// byLength type to sort component name list by it's name length
type byLength []string

func (s byLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}
func (s byLength) Len() int {
	return len(s)
}
func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Terraform call - terraform apply
// This function creates a new vRA 7 Deployment using configuration in a user's Terraform file.
// The Deployment is produced by invoking a catalog item that is specified in the configuration.
func createResource(d *schema.ResourceData, meta interface{}) error {
	vraClient = meta.(*client.APIClient)
	// Get client handle
	readProviderConfiguration(d)

	requestTemplate, validityErr := checkConfigValuesValidity(d)
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
	catalogRequest, err := sdk.RequestCatalogItem(requestTemplate)

	if err != nil {
		log.Errorf("Resource Machine Request Failed: %v", err)
		return fmt.Errorf("Resource Machine Request Failed: %v", err)
	}
	//Set request ID
	d.SetId(catalogRequest.ID)
	//Set request status
	d.Set(utils.RequestStatus, utils.Submitted)
	return waitForRequestCompletion(d, meta)
}

func updateRequestTemplate(templateInterface map[string]interface{}, field string, value interface{}) map[string]interface{} {
	var replaced bool
	templateInterface, replaced = utils.ReplaceValueInRequestTemplate(templateInterface, field, value)

	if !replaced {
		templateInterface["data"] = utils.AddValueToRequestTemplate(templateInterface["data"].(map[string]interface{}), field, value)
	}
	return templateInterface
}

// Terraform call - terraform apply
// This function updates the state of a vRA 7 Deployment when changes to a Terraform file are applied.
// The update is performed on the Deployment using supported (day-2) actions.
func updateResource(d *schema.ResourceData, meta interface{}) error {
	vraClient = meta.(*client.APIClient)
	// Get the ID of the catalog request that was used to provision this Deployment.
	catalogItemRequestID := d.Id()
	// Get client handle
	readProviderConfiguration(d)

	requestTemplate, validityErr := checkConfigValuesValidity(d)
	if validityErr != nil {
		return validityErr
	}
	validityErr = checkResourceConfigValidity(requestTemplate)
	if validityErr != nil {
		return validityErr
	}

	resourceActions, err := sdk.GetResourceActions(catalogItemRequestID)
	if err != nil {
		log.Errorf("Error while reading resource actions for the request %v: %v ", catalogItemRequestID, err.Error())
		return fmt.Errorf("Error while reading resource actions for the request %v: %v  ", catalogItemRequestID, err.Error())
	}

	// If any change made in resource_configuration.
	if d.HasChange(utils.ResourceConfiguration) {
		for _, resources := range resourceActions.Content {
			if resources.ResourceTypeRef.ID == utils.InfrastructureVirtual {
				var reconfigureEnabled bool
				var reconfigureActionID string
				for _, op := range resources.Operations {
					if op.Name == utils.Reconfigure {
						reconfigureEnabled = true
						reconfigureActionID = op.OperationID
						break
					}
				}
				// if reconfigure action is not available for any resource of the deployment
				// return with an error message
				if !reconfigureEnabled {
					return fmt.Errorf("Update is not allowed for resource %v, your entitlement has no Reconfigure action enabled", resources.ID)
				}
				resourceData := resources.ResourceData
				entries := resourceData.Entries
				for _, entry := range entries {
					if entry.Key == utils.Component {
						entryValue := entry.Value
						var componentName string
						for k, v := range entryValue {
							if k == "value" {
								componentName = v.(string)
							}
						}
						log.Info("Retrieving reconfigure action template for the component: %v ", componentName)

						resourceActionTemplate, err := sdk.GetResourceActionTemplate(resources.ID, reconfigureActionID)
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
								resourceActionTemplate.Data, returnFlag = utils.ReplaceValueInRequestTemplate(
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
							err := sdk.PostResourceAction(resources.ID, reconfigureActionID, resourceActionTemplate)
							if err != nil {
								oldData, _ := d.GetChange(utils.ResourceConfiguration)
								d.Set(utils.ResourceConfiguration, oldData)
								log.Errorf("The update request failed with error: %v ", err)
								return err
							}
						}
					}
				}
			}
		}
	}
	return waitForRequestCompletion(d, meta)
}

// Terraform call - terraform refresh
// This function retrieves the latest state of a vRA 7 deployment. Terraform updates its state based on
// the information returned by this function.
func readResource(d *schema.ResourceData, meta interface{}) error {
	vraClient = meta.(*client.APIClient)
	// Get the ID of the catalog request that was used to provision this Deployment.
	catalogItemRequestID := d.Id()

	log.Info("Calling read resource to get the current resource status of the request: %v ", catalogItemRequestID)
	// Get client handle
	// Get requested status
	resourceTemplate, errTemplate := sdk.GetRequestStatus(catalogItemRequestID)

	if errTemplate != nil {
		log.Errorf("Resource view failed to load:  %v", errTemplate)
		return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
	}

	// Update resource request status in state file
	d.Set(utils.RequestStatus, resourceTemplate.Phase)
	// If request is failed then set failed message in state file
	if resourceTemplate.Phase == utils.Failed {
		log.Errorf(resourceTemplate.RequestCompletion.CompletionDetails)
		d.Set(utils.FailedMessage, resourceTemplate.RequestCompletion.CompletionDetails)
	}

	requestResourceView, errTemplate := sdk.GetRequestResourceView(catalogItemRequestID)
	if errTemplate != nil {
		return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
	}

	resourceDataMap := make(map[string]map[string]interface{})
	for _, resource := range requestResourceView.Content {
		if resource.ResourceType == utils.InfrastructureVirtual {
			resourceData := resource.ResourcesData
			log.Info("The resource data map of the resource %v is: \n%v", resourceData.Component, resource.ResourcesData)
			dataVals := make(map[string]interface{})
			resourceDataMap[resourceData.Component] = dataVals
			dataVals[utils.MachineCPU] = resourceData.CPU
			dataVals[utils.MachineStorage] = resourceData.Storage
			dataVals[utils.IPAddress] = resourceData.IPAddress
			dataVals[utils.MachineMemory] = resourceData.Memory
			dataVals[utils.MachineName] = resourceData.MachineName
			dataVals[utils.MachineGuestOs] = resourceData.MachineGuestOperatingSystem
			dataVals[utils.MachineBpName] = resourceData.MachineBlueprintName
			dataVals[utils.MachineType] = resourceData.MachineType
			dataVals[utils.MachineReservationName] = resourceData.MachineReservationName
			dataVals[utils.MachineInterfaceType] = resourceData.MachineInterfaceType
			dataVals[utils.MachineID] = resourceData.MachineID
			dataVals[utils.MachineGroupName] = resourceData.MachineGroupName
			dataVals[utils.MachineDestructionDate] = resourceData.MachineDestructionDate
		}
	}
	resourceConfiguration, _ := d.Get(utils.ResourceConfiguration).(map[string]interface{})
	changed := false

	resourceConfiguration, changed = utils.UpdateResourceConfigurationMap(resourceConfiguration, resourceDataMap)

	if changed {
		setError := d.Set(utils.ResourceConfiguration, resourceConfiguration)
		if setError != nil {
			return fmt.Errorf(setError.Error())
		}
	}
	return nil
}

//Function use - To delete resources which are created by terraform and present in state file
//Terraform call - terraform destroy
func deleteResource(d *schema.ResourceData, meta interface{}) error {
	vraClient = meta.(*client.APIClient)
	//Get requester machine ID from schema.dataresource
	catalogItemRequestID := d.Id()
	// Throw an error if request ID has no value or empty value
	if len(d.Id()) == 0 {
		return fmt.Errorf("Resource not found")
	}
	log.Info("Calling delete resource for the request id %v ", catalogItemRequestID)

	resourceActions, err := sdk.GetResourceActions(catalogItemRequestID)
	if err != nil {
		return err
	}

	for _, resources := range resourceActions.Content {
		if resources.ResourceTypeRef.ID == utils.DeploymentResourceType {
			deploymentName := resources.Name
			var destroyEnabled bool
			var destroyActionID string
			for _, op := range resources.Operations {
				if op.Name == utils.Destroy {
					destroyEnabled = true
					destroyActionID = op.OperationID
					break
				}
			}
			// if destroy deployment action is not available for the deployment
			// return with an error message
			if !destroyEnabled {
				return fmt.Errorf("The deployment %v cannot be destroyed, your entitlement has no Destroy Deployment action enabled", deploymentName)
			}
			resourceActionTemplate, err := sdk.GetResourceActionTemplate(resources.ID, destroyActionID)
			if err != nil {
				log.Errorf(utils.DestroyActionTemplateError, deploymentName, err.Error())
				return fmt.Errorf(utils.DestroyActionTemplateError, deploymentName, err.Error())
			}
			err = sdk.PostResourceAction(resources.ID, destroyActionID, resourceActionTemplate)
			if err != nil {
				log.Errorf("The destroy deployment request failed with error: %v ", err)
				return err
			}
		}
	}

	waitTimeout := d.Get(utils.WaitTimeout).(int) * 60
	sleepFor := 30
	for i := 0; i < waitTimeout/sleepFor; i++ {
		time.Sleep(time.Duration(sleepFor) * time.Second)
		log.Info("Checking to see if resource is deleted.")
		deploymentStateData, err := sdk.GetDeploymentState(catalogItemRequestID)
		if err != nil {
			return fmt.Errorf("Resource view failed to load:  %v", err)
		}
		// If resource create status is in_progress then skip delete call and throw an exception
		// Note: vRA API should return error on destroy action if the request is in progress. Filed a bug
		if d.Get(utils.RequestStatus).(string) == utils.InProgress {
			return fmt.Errorf("Machine cannot be deleted while request is in-progress state. Please try again later. \nRun terraform refresh to get the latest state of your request")
		}

		if len(deploymentStateData.Content) == 0 {
			//If resource got deleted then unset the resource ID from state file
			d.SetId("")
			break
		}
	}
	if d.Id() != "" {
		d.SetId("")
		return fmt.Errorf("Resource still being deleted after %v minutes", d.Get(utils.WaitTimeout))
	}
	return nil
}

// check if the resource configuration is valid in the terraform config file
func checkResourceConfigValidity(requestTemplate *utils.CatalogItemRequestTemplate) error {
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
		return fmt.Errorf(utils.ConfigInvalidError, strings.Join(invalidKeys, ", "))
	}
	return nil
}

// check if the values provided in the config file are valid and set
// them in the resource schema. Requires to call APIs
func checkConfigValuesValidity(d *schema.ResourceData) (*utils.CatalogItemRequestTemplate, error) {
	// 	// If catalog_name and catalog_id both not provided then return an error
	if len(catalogItemName) <= 0 && len(catalogItemID) <= 0 {
		return nil, fmt.Errorf("Either catalog_name or catalog_id should be present in given configuration")
	}

	var catalogItemIDFromName string
	var catalogItemNameFromID string
	var err error
	// if catalog item id is provided, fetch the catalog item name
	if len(catalogItemName) > 0 {
		catalogItemIDFromName, err = sdk.ReadCatalogItemByName(catalogItemName)
		if err != nil || catalogItemIDFromName == "" {
			return nil, fmt.Errorf("Error in finding catalog item id corresponding to the catlog item name %v: \n %v", catalogItemName, err)
		}
		log.Info("The catalog item id provided in the config is %v\n", catalogItemIDFromName)
	}

	// if catalog item name is provided, fetch the catalog item id
	if len(catalogItemID) > 0 { // else if both are provided and matches or just id is provided, use id
		catalogItemNameFromID, err = sdk.ReadCatalogItemNameByID(catalogItemID)
		if err != nil || catalogItemNameFromID == "" {
			return nil, fmt.Errorf("Error in finding catalog item name corresponding to the catlog item id %v: \n %v", catalogItemID, err)
		}
		log.Info("The catalog item name corresponding to the catalog item id in the config is:  %v\n", catalogItemNameFromID)
	}

	// if both catalog item name and id are provided but does not belong to the same catalog item, throw an error
	if len(catalogItemName) > 0 && len(catalogItemID) > 0 && (catalogItemIDFromName != catalogItemID || catalogItemNameFromID != catalogItemName) {
		log.Error(utils.CatalogItemIDNameNotMatchingErr, catalogItemName, catalogItemID)
		return nil, fmt.Errorf(utils.CatalogItemIDNameNotMatchingErr, catalogItemName, catalogItemID)
	} else if len(catalogItemID) > 0 { // else if both are provided and matches or just id is provided, use id
		d.Set(utils.CatalogID, catalogItemID)
		d.Set(utils.CatalogName, catalogItemNameFromID)
	} else if len(catalogItemName) > 0 { // else if name is provided, use the id fetched from the name
		d.Set(utils.CatalogID, catalogItemIDFromName)
		d.Set(utils.CatalogName, catalogItemName)
	}

	// update the catalogItemID var with the updated id
	catalogItemID = d.Get(utils.CatalogID).(string)

	// Get request template for catalog item.
	requestTemplate, err := sdk.GetCatalogItemRequestTemplate(catalogItemID)
	if err != nil {
		return nil, err
	}
	log.Info("The request template data corresponding to the catalog item %v is: \n %v\n", catalogItemID, requestTemplate.Data)

	for field1 := range catalogConfiguration {
		requestTemplate.Data[field1] = catalogConfiguration[field1]

	}
	// get the business group id from name
	var businessGroupIDFromName string
	if len(businessGroupName) > 0 {
		businessGroupIDFromName, err = sdk.GetBusinessGroupID(businessGroupName, vraClient.Tenant)
		if err != nil || businessGroupIDFromName == "" {
			return nil, err
		}
	}

	//if both business group name and id are provided but does not belong to the same business group, throw an error
	if len(businessGroupName) > 0 && len(businessGroupID) > 0 && businessGroupIDFromName != businessGroupID {
		log.Error(utils.BusinessGroupIDNameNotMatchingErr, businessGroupName, businessGroupID)
		return nil, fmt.Errorf(utils.BusinessGroupIDNameNotMatchingErr, businessGroupName, businessGroupID)
	} else if len(businessGroupID) > 0 { // else if both are provided and matches or just id is provided, use id
		log.Info("Setting business group id %s ", businessGroupID)
		requestTemplate.BusinessGroupID = businessGroupID
	} else if len(businessGroupName) > 0 { // else if name is provided, use the id fetched from the name
		log.Info("Setting business group id %s for the group %s ", businessGroupIDFromName, businessGroupName)
		requestTemplate.BusinessGroupID = businessGroupIDFromName
	}
	return requestTemplate, nil
}

// check the request status on apply and update
func waitForRequestCompletion(d *schema.ResourceData, meta interface{}) error {

	waitTimeout := d.Get(utils.WaitTimeout).(int) * 60
	sleepFor := 30
	requestStatus := ""
	for i := 0; i < waitTimeout/sleepFor; i++ {
		log.Info("Waiting for %d seconds before checking request status.", sleepFor)
		time.Sleep(time.Duration(sleepFor) * time.Second)

		readResource(d, meta)

		requestStatus = d.Get(utils.RequestStatus).(string)
		log.Info("Checking to see if resource is created. Status: %s.", requestStatus)
		if requestStatus == utils.Successful {
			log.Info("Resource creation SUCCESSFUL.")
			return nil
		} else if requestStatus == utils.Failed {
			log.Error("Request Failed with message %v ", d.Get(utils.FailedMessage))
			//If request is failed during the time then
			//unset resource details from state.
			d.SetId("")
			return fmt.Errorf("Request failed \n %v ", d.Get(utils.FailedMessage))
		} else if requestStatus == utils.InProgress {
			log.Info("Resource creation is still IN PROGRESS.")
		} else {
			log.Info("Resource creation request status: %s.", requestStatus)
		}
	}

	// The execution has timed out while still IN PROGRESS.
	// The user will need to use 'terraform refresh' at a later point to resolve this.
	return fmt.Errorf("Resource creation has timed out")
}

// read the config file
func readProviderConfiguration(d *schema.ResourceData) {

	log.Info("Reading the provider configuration data.....")
	catalogItemName = strings.TrimSpace(d.Get(utils.CatalogName).(string))
	log.Info("Catalog item name: %v ", catalogItemName)
	catalogItemID = strings.TrimSpace(d.Get(utils.CatalogID).(string))
	log.Info("Catalog item ID: %v", catalogItemID)
	businessGroupName = strings.TrimSpace(d.Get(utils.BusinessGroupName).(string))
	log.Info("Business Group name: %v", businessGroupName)
	businessGroupID = strings.TrimSpace(d.Get(utils.BusinessGroupID).(string))
	log.Info("Business Group Id: %v", businessGroupID)
	waitTimeout = d.Get(utils.WaitTimeout).(int) * 60
	log.Info("Wait time out: %v ", waitTimeout)
	failedMessage = strings.TrimSpace(d.Get(utils.FailedMessage).(string))
	log.Info("Failed message: %v ", failedMessage)
	resourceConfiguration = d.Get(utils.ResourceConfiguration).(map[string]interface{})
	log.Info("Resource Configuration: %v ", resourceConfiguration)
	deploymentConfiguration = d.Get(utils.DeploymentConfiguration).(map[string]interface{})
	log.Info("Deployment Configuration: %v ", deploymentConfiguration)
	catalogConfiguration = d.Get(utils.CatalogConfiguration).(map[string]interface{})
	log.Info("Catalog configuration: %v ", catalogConfiguration)
}
