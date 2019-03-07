package vrealize

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	logging "github.com/op/go-logging"
	"github.com/vmware/terraform-provider-vra7/sdk"
	"github.com/vmware/terraform-provider-vra7/utils"
)

// error constants
const (
	ConfigInvalidError                = "The resource_configuration in the config file has invalid component name(s): %v "
	DestroyActionTemplateError        = "Error retrieving destroy action template for the deployment %v: %v "
	BusinessGroupIDNameNotMatchingErr = "The business group name %s and id %s does not belong to the same business group, provide either name or id"
	CatalogItemIDNameNotMatchingErr   = "The catalog item name %s and id %s does not belong to the same catalog item, provide either name or id"
)

var (
	log       = logging.MustGetLogger(utils.LoggerID)
	vraClient *sdk.APIClient
)

// ProviderSchema represents the information provided in the tf file
type ProviderSchema struct {
	CatalogItemName         string
	CatalogItemID           string
	BusinessGroupID         string
	BusinessGroupName       string
	WaitTimeout             int
	RequestStatus           string
	FailedMessage           string
	DeploymentConfiguration map[string]interface{}
	ResourceConfiguration   map[string]interface{}
	CatalogConfiguration    map[string]interface{}
}

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
	vraClient = meta.(*sdk.APIClient)
	// Get client handle
	p := readProviderConfiguration(d)

	requestTemplate, validityErr := p.checkConfigValuesValidity(d)
	if validityErr != nil {
		return validityErr
	}
	validityErr = p.checkResourceConfigValidity(requestTemplate)
	if validityErr != nil {
		return validityErr
	}

	for field1 := range p.CatalogConfiguration {
		requestTemplate.Data[field1] = p.CatalogConfiguration[field1]
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
	for configKey, configValue := range p.ResourceConfiguration {
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
	for depField, depValue := range p.DeploymentConfiguration {
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
	catalogRequest, err := vraClient.RequestCatalogItem(requestTemplate)

	if err != nil {
		log.Errorf("Resource Machine Request Failed: %v", err)
		return fmt.Errorf("Resource Machine Request Failed: %v", err)
	}
	//Set request ID
	d.SetId(catalogRequest.ID)
	//Set request status
	d.Set(utils.RequestStatus, sdk.Submitted)
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
	vraClient = meta.(*sdk.APIClient)
	// Get the ID of the catalog request that was used to provision this Deployment.
	catalogItemRequestID := d.Id()
	// Get client handle
	p := readProviderConfiguration(d)

	requestTemplate, validityErr := p.checkConfigValuesValidity(d)
	if validityErr != nil {
		return validityErr
	}
	validityErr = p.checkResourceConfigValidity(requestTemplate)
	if validityErr != nil {
		return validityErr
	}

	resourceActions, err := vraClient.GetResourceActions(catalogItemRequestID)
	if err != nil {
		log.Errorf("Error while reading resource actions for the request %v: %v ", catalogItemRequestID, err.Error())
		return fmt.Errorf("Error while reading resource actions for the request %v: %v  ", catalogItemRequestID, err.Error())
	}

	if d.HasChange(utils.DeploymentConfiguration) {

	}

	// If any change made in resource_configuration.
	if d.HasChange(utils.ResourceConfiguration) {
		for _, resources := range resourceActions.Content {
			if resources.ResourceTypeRef.ID == sdk.InfrastructureVirtual {
				var reconfigureEnabled bool
				var reconfigureActionID string
				for _, op := range resources.Operations {
					if op.Name == sdk.Reconfigure {
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
					if entry.Key == sdk.Component {
						entryValue := entry.Value
						var componentName string
						for k, v := range entryValue {
							if k == "value" {
								componentName = v.(string)
							}
						}
						log.Info("Retrieving reconfigure action template for the component: %v ", componentName)

						resourceActionTemplate, err := vraClient.GetResourceActionTemplate(resources.ID, reconfigureActionID)
						if err != nil {
							log.Errorf("Error retrieving reconfigure action template for the component %v: %v ", componentName, err.Error())
							return fmt.Errorf("Error retrieving reconfigure action template for the component %v: %v ", componentName, err.Error())
						}
						configChanged := false
						returnFlag := false
						for configKey := range p.ResourceConfiguration {
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
									p.ResourceConfiguration[configKey])
								if returnFlag == true {
									configChanged = true
								}
							}
						}
						// If template value got changed then set post call and update resource child
						if configChanged != false {
							err := vraClient.PostResourceAction(resources.ID, reconfigureActionID, resourceActionTemplate)
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
	vraClient = meta.(*sdk.APIClient)
	// Get the ID of the catalog request that was used to provision this Deployment.
	catalogItemRequestID := d.Id()

	log.Info("Calling read resource to get the current resource status of the request: %v ", catalogItemRequestID)
	// Get client handle
	// Get requested status
	resourceTemplate, errTemplate := vraClient.GetRequestStatus(catalogItemRequestID)

	if errTemplate != nil {
		log.Errorf("Resource view failed to load:  %v", errTemplate)
		return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
	}

	// Update resource request status in state file
	d.Set(utils.RequestStatus, resourceTemplate.Phase)
	// If request is failed then set failed message in state file
	if resourceTemplate.Phase == sdk.Failed {
		log.Errorf(resourceTemplate.RequestCompletion.CompletionDetails)
		d.Set(utils.FailedMessage, resourceTemplate.RequestCompletion.CompletionDetails)
	}

	requestResourceView, errTemplate := vraClient.GetRequestResourceView(catalogItemRequestID)
	if errTemplate != nil {
		return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
	}

	resourceDataMap := make(map[string]map[string]interface{})
	for _, resource := range requestResourceView.Content {
		if resource.ResourceType == sdk.InfrastructureVirtual {
			resourceData := resource.ResourcesData
			log.Info("The resource data map of the resource %v is: \n%v", resourceData.Component, resource.ResourcesData)
			dataVals := make(map[string]interface{})
			resourceDataMap[resourceData.Component] = dataVals
			dataVals[sdk.MachineCPU] = resourceData.CPU
			dataVals[sdk.MachineStorage] = resourceData.Storage
			dataVals[sdk.IPAddress] = resourceData.IPAddress
			dataVals[sdk.MachineMemory] = resourceData.Memory
			dataVals[sdk.MachineName] = resourceData.MachineName
			dataVals[sdk.MachineGuestOs] = resourceData.MachineGuestOperatingSystem
			dataVals[sdk.MachineBpName] = resourceData.MachineBlueprintName
			dataVals[sdk.MachineType] = resourceData.MachineType
			dataVals[sdk.MachineReservationName] = resourceData.MachineReservationName
			dataVals[sdk.MachineInterfaceType] = resourceData.MachineInterfaceType
			dataVals[sdk.MachineID] = resourceData.MachineID
			dataVals[sdk.MachineGroupName] = resourceData.MachineGroupName
			dataVals[sdk.MachineDestructionDate] = resourceData.MachineDestructionDate
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
	vraClient = meta.(*sdk.APIClient)
	//Get requester machine ID from schema.dataresource
	catalogItemRequestID := d.Id()
	// Throw an error if request ID has no value or empty value
	if len(d.Id()) == 0 {
		return fmt.Errorf("Resource not found")
	}
	log.Info("Calling delete resource for the request id %v ", catalogItemRequestID)

	resourceActions, err := vraClient.GetResourceActions(catalogItemRequestID)
	if err != nil {
		return err
	}

	for _, resources := range resourceActions.Content {
		if resources.ResourceTypeRef.ID == sdk.DeploymentResourceType {
			deploymentName := resources.Name
			var destroyEnabled bool
			var destroyActionID string
			for _, op := range resources.Operations {
				if op.Name == sdk.Destroy {
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
			resourceActionTemplate, err := vraClient.GetResourceActionTemplate(resources.ID, destroyActionID)
			if err != nil {
				log.Errorf(DestroyActionTemplateError, deploymentName, err.Error())
				return fmt.Errorf(DestroyActionTemplateError, deploymentName, err.Error())
			}
			err = vraClient.PostResourceAction(resources.ID, destroyActionID, resourceActionTemplate)
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
		resourceView, err := vraClient.GetRequestResourceView(catalogItemRequestID)
		if err != nil {
			return fmt.Errorf("Resource view failed to load:  %v", err)
		}
		// If resource create status is in_progress then skip delete call and throw an exception
		// Note: vRA API should return error on destroy action if the request is in progress. Filed a bug
		if d.Get(utils.RequestStatus).(string) == sdk.InProgress {
			return fmt.Errorf("Machine cannot be deleted while request is in-progress state. Please try again later. \nRun terraform refresh to get the latest state of your request")
		}

		if len(resourceView.Content) == 0 {
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
func (p *ProviderSchema) checkResourceConfigValidity(requestTemplate *sdk.CatalogItemRequestTemplate) error {
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
	for k := range p.ResourceConfiguration {
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
		return fmt.Errorf(ConfigInvalidError, strings.Join(invalidKeys, ", "))
	}
	return nil
}

// check if the values provided in the config file are valid and set
// them in the resource schema. Requires to call APIs
func (p *ProviderSchema) checkConfigValuesValidity(d *schema.ResourceData) (*sdk.CatalogItemRequestTemplate, error) {
	// 	// If catalog_name and catalog_id both not provided then return an error
	if len(p.CatalogItemName) <= 0 && len(p.CatalogItemID) <= 0 {
		return nil, fmt.Errorf("Either catalog_name or catalog_id should be present in given configuration")
	}

	var catalogItemIDFromName string
	var catalogItemNameFromID string
	var err error
	// if catalog item id is provided, fetch the catalog item name
	if len(p.CatalogItemName) > 0 {
		catalogItemIDFromName, err = vraClient.ReadCatalogItemByName(p.CatalogItemName, 1)
		if err != nil || catalogItemIDFromName == "" {
			return nil, fmt.Errorf("Error in finding catalog item id corresponding to the catlog item name %v: \n %v", p.CatalogItemName, err)
		}
		log.Info("The catalog item id provided in the config is %v\n", catalogItemIDFromName)
	}

	// if catalog item name is provided, fetch the catalog item id
	if len(p.CatalogItemID) > 0 { // else if both are provided and matches or just id is provided, use id
		catalogItemNameFromID, err = vraClient.ReadCatalogItemNameByID(p.CatalogItemID)
		if err != nil || catalogItemNameFromID == "" {
			return nil, fmt.Errorf("Error in finding catalog item name corresponding to the catlog item id %v: \n %v", p.CatalogItemID, err)
		}
		log.Info("The catalog item name corresponding to the catalog item id in the config is:  %v\n", catalogItemNameFromID)
	}

	// if both catalog item name and id are provided but does not belong to the same catalog item, throw an error
	if len(p.CatalogItemName) > 0 && len(p.CatalogItemID) > 0 && (catalogItemIDFromName != p.CatalogItemID || catalogItemNameFromID != p.CatalogItemName) {
		log.Error(CatalogItemIDNameNotMatchingErr, p.CatalogItemName, p.CatalogItemID)
		return nil, fmt.Errorf(CatalogItemIDNameNotMatchingErr, p.CatalogItemName, p.CatalogItemID)
	} else if len(p.CatalogItemID) > 0 { // else if both are provided and matches or just id is provided, use id
		d.Set(utils.CatalogID, p.CatalogItemID)
		d.Set(utils.CatalogName, catalogItemNameFromID)
	} else if len(p.CatalogItemName) > 0 { // else if name is provided, use the id fetched from the name
		d.Set(utils.CatalogID, catalogItemIDFromName)
		d.Set(utils.CatalogName, p.CatalogItemName)
	}

	// update the catalogItemID var with the updated id
	p.CatalogItemID = d.Get(utils.CatalogID).(string)

	// Get request template for catalog item.
	requestTemplate, err := vraClient.GetCatalogItemRequestTemplate(p.CatalogItemID)
	if err != nil {
		return nil, err
	}
	log.Info("The request template data corresponding to the catalog item %v is: \n %v\n", p.CatalogItemID, requestTemplate.Data)

	for field1 := range p.CatalogConfiguration {
		requestTemplate.Data[field1] = p.CatalogConfiguration[field1]

	}
	// get the business group id from name
	var businessGroupIDFromName string
	if len(p.BusinessGroupName) > 0 {
		businessGroupIDFromName, err = vraClient.GetBusinessGroupID(p.BusinessGroupName, vraClient.Tenant)
		if err != nil || businessGroupIDFromName == "" {
			return nil, err
		}
	}

	//if both business group name and id are provided but does not belong to the same business group, throw an error
	if len(p.BusinessGroupName) > 0 && len(p.BusinessGroupID) > 0 && businessGroupIDFromName != p.BusinessGroupID {
		log.Error(BusinessGroupIDNameNotMatchingErr, p.BusinessGroupName, p.BusinessGroupID)
		return nil, fmt.Errorf(BusinessGroupIDNameNotMatchingErr, p.BusinessGroupName, p.BusinessGroupID)
	} else if len(p.BusinessGroupID) > 0 { // else if both are provided and matches or just id is provided, use id
		log.Info("Setting business group id %s ", p.BusinessGroupID)
		requestTemplate.BusinessGroupID = p.BusinessGroupID
	} else if len(p.BusinessGroupName) > 0 { // else if name is provided, use the id fetched from the name
		log.Info("Setting business group id %s for the group %s ", businessGroupIDFromName, p.BusinessGroupName)
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
		if requestStatus == sdk.Successful {
			log.Info("Resource creation SUCCESSFUL.")
			return nil
		} else if requestStatus == sdk.Failed {
			log.Error("Request Failed with message %v ", d.Get(utils.FailedMessage))
			//If request is failed during the time then
			//unset resource details from state.
			d.SetId("")
			return fmt.Errorf("Request failed \n %v ", d.Get(utils.FailedMessage))
		} else if requestStatus == sdk.InProgress {
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
func readProviderConfiguration(d *schema.ResourceData) *ProviderSchema {

	log.Info("Reading the provider configuration data.....")
	providerSchema := ProviderSchema{
		CatalogItemName:         strings.TrimSpace(d.Get(utils.CatalogName).(string)),
		CatalogItemID:           strings.TrimSpace(d.Get(utils.CatalogID).(string)),
		BusinessGroupName:       strings.TrimSpace(d.Get(utils.BusinessGroupName).(string)),
		BusinessGroupID:         strings.TrimSpace(d.Get(utils.BusinessGroupID).(string)),
		WaitTimeout:             d.Get(utils.WaitTimeout).(int) * 60,
		FailedMessage:           strings.TrimSpace(d.Get(utils.FailedMessage).(string)),
		ResourceConfiguration:   d.Get(utils.ResourceConfiguration).(map[string]interface{}),
		DeploymentConfiguration: d.Get(utils.DeploymentConfiguration).(map[string]interface{}),
		CatalogConfiguration:    d.Get(utils.CatalogConfiguration).(map[string]interface{}),
	}

	log.Info("The values provided in the TF config file is: \n %v ", providerSchema)
	return &providerSchema
}
