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

// byLength yype to sort component name list by it's name length
type byLength []string

//GetActionTemplate - set call for read template/blueprint
func (c *APIClient) GetActionTemplate(resourceViewsTemplate *ResourceView, actionURLString string) (*ActionTemplate, *ResourceView, error) {
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
func getactionURL(template *ResourceView, relationVal string) (templateactionURL string) {
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
func (c *APIClient) GetPowerOffActionTemplate(resourceData *ResourceView) (*ActionTemplate, *ResourceView, error) {
	//Set resource power-off URL label
	actionURL := "GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.PowerOff}"
	//Set get action URL function call
	return c.GetActionTemplate(resourceData, actionURL)
}

//GetDestroyActionTemplate - To read destroy resource action template from provided resource configuration
func (c *APIClient) GetDestroyActionTemplate(resourceData *ResourceView) (*ActionTemplate, *ResourceView, error) {
	//Set destroy resource URL label
	actionURL := "GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.destroy.name}"
	//Set get action URL function call
	return c.GetActionTemplate(resourceData, actionURL)
}

func (s byLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}
func (s byLength) Len() int {
	return len(s)
}
func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// getDeploymentLeaseActionLinks returns a MAP of GET/POST action links of change lease
func (vRAClient *APIClient) getDeploymentLeaseActionLinks(resourceData *ResourceView) (map[string]interface{}) {
	//Set resource power-off URL label
	const templateRel = "GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changelease.name}"
	const postRel = "POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changelease.name}"
	//Set get action URL function call
	actionData := vRAClient.GetDeploymentActionLinks(resourceData, templateRel, postRel)
	return actionData
}

// GetDeploymentActionLinks determines change lease action links from resourceView template
func (vRAClient *APIClient) GetDeploymentActionLinks(resourceData *ResourceView, templateRel string, postRel string) map[string]interface{}  {
	data := map[string]interface{}{}
	for _, value := range resourceData.Content {
		resourceMap := value.(map[string]interface{})
		if resourceMap["resourceType"] == "composition.resource.type.deployment" {
			resourceSpecificLinks := resourceMap["links"].([]interface{})
			data["template_url"] = readActionLink(resourceSpecificLinks, templateRel)
			data["post_url"] = readActionLink(resourceSpecificLinks, postRel)
			break
		}

	}
	return data
}

