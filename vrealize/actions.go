package vrealize

import (
	"fmt"
)

//ActionTemplate - is used to store action template
//which is further used to make REST POST call for any action
//for example - poweroff action, destroy action.
type ActionTemplate struct {
	ActionID string `json:"actionId"`
	Data     struct {
		Description  interface{} `json:"description"`
		Reasons      interface{} `json:"reasons"`
		ForceDestroy bool        `json:"ForceDestroy"`
	} `json:"data"`
	Description interface{} `json:"description"`
	ResourceID  string      `json:"resourceId"`
	Type        string      `json:"type"`
}

//GetActionTemplate - set call for read template/blueprint
func (c *APIClient) GetActionTemplate(resourceViewsTemplate *ResourceViewsTemplate, actionURLString string) (*ActionTemplate, *ResourceViewsTemplate, error) {
	//Fetch an action URL from given template
	actionURL := getactionURL(resourceViewsTemplate, actionURLString)

	//Raise an error if action URL not found
	if len(actionURL) == 0 {
		return nil, resourceViewsTemplate, fmt.Errorf("resource is not created or not found")
	}

	actionTemplate := new(ActionTemplate)
	apiError := new(APIError)

	//Set a REST call to perform an action on resource
	_, err := c.HTTPClient.New().Get(actionURL).Receive(actionTemplate, apiError)

	if err != nil {
		return nil, resourceViewsTemplate, err
	}

	if !apiError.isEmpty() {
		return nil, resourceViewsTemplate, apiError
	}

	return actionTemplate, resourceViewsTemplate, nil
}

//getactionURL - Read action URL from provided template of resource item
func getactionURL(template *ResourceViewsTemplate, relationVal string) (templateactionURL string) {
	var actionURL string
	l := len(template.Content)
	//Loop to iterate over the action URLs
	for i := 0; i < l; i++ {
		content := template.Content[i].(map[string]interface{})
		links := content["links"].([]interface{})
		lengthLinks := len(links)
		for j := 0; j < lengthLinks; j++ {
			linkObj := links[j].(map[string]interface{})
			//If template action URL matches with given URL then store it in actionURL var
			if linkObj["rel"] == relationVal {
				actionURL = linkObj["href"].(string)
			}

		}

	}
	//Return action URL
	return actionURL
}

//GetPowerOffActionTemplate - To read power-off action template from provided resource configuration
func (c *APIClient) GetPowerOffActionTemplate(resourceViewsTemplate *ResourceViewsTemplate) (*ActionTemplate, *ResourceViewsTemplate, error) {
	//Set resource power-off URL label
	actionURL := "GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.PowerOff}"
	//Set get action URL function call
	return c.GetActionTemplate(resourceViewsTemplate, actionURL)
}

//GetDestroyActionTemplate - To read destroy resource action template from provided resource configuration
func (c *APIClient) GetDestroyActionTemplate(resourceViewsTemplate *ResourceViewsTemplate) (*ActionTemplate, *ResourceViewsTemplate, error) {
	//Set destroy resource URL label
	actionURL := "GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.destroy.name}"
	//Set get action URL function call
	return c.GetActionTemplate(resourceViewsTemplate, actionURL)
}
