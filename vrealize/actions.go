package vrealize

import (
	"fmt"
)

//ActionTemplate - is used to store action template
//which is further used to make REST POST call for any action
//for example - poweroff action, destroy action.
type ActionTemplate struct {
	ActionID    string      `json:"actionId"`
	Data        interface{} `json:"data"`
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
		lengthLinks := len(template.Content[i].Links)
		for j := 0; j < lengthLinks; j++ {
			//If template action URL matches with given URL then store it in actionURL var
			if template.Content[i].Links[j].Rel == relationVal {
				actionURL = template.Content[i].Links[j].Href
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

//getScaleInActionTemplate returns scale in resource action template from provided resource comfiguration template
func (c *APIClient) getScaleInActionTemplate(resourceViewsTemplate *ResourceViewsTemplate) (*ActionTemplate, *ResourceViewsTemplate, error) {
	//Set destroy resource URL label
	actionURL := "GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scalein.name}"
	//Set get action URL function call
	return c.GetActionTemplate(resourceViewsTemplate, actionURL)
}

//getScaleInActionPostURL returns scale in resource action URL
//Which is use to set actual POST call of scale in deployment
func (c *APIClient) getScaleInActionPostURL(resourceViewsTemplate *ResourceViewsTemplate) string {
	//Set destroy resource URL label
	actionURL := "POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scalein.name}"
	//Set get action URL function call
	return getactionURL(resourceViewsTemplate, actionURL)
}

//getScaleOutActionTemplate returns scale out resource action URL
//Which is use to set actual POST call of scale out deployment
func (c *APIClient) getScaleOutActionPostURL(resourceViewsTemplate *ResourceViewsTemplate) string {
	//Set destroy resource URL label
	actionURL := "POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scaleout.name}"
	//Set get action URL function call
	return getactionURL(resourceViewsTemplate, actionURL)
}

//getScaleOutActionPostURL returns scale out resource action template from provided resource comfiguration template
func (c *APIClient) getScaleOutActionTemplate(resourceViewsTemplate *ResourceViewsTemplate) (*ActionTemplate, *ResourceViewsTemplate, error) {
	//Set destroy resource URL label
	actionURL := "GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scaleout.name}"
	//Set get action URL function call
	return c.GetActionTemplate(resourceViewsTemplate, actionURL)
}
