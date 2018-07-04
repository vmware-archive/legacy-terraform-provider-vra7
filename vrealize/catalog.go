package vrealize

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

//CatalogItemRequestTemplate - A structure that captures a catalog request template, to be filled in and POSTED.
type CatalogItemRequestTemplate struct {
	Type            string                 `json:"type"`
	CatalogItemID   string                 `json:"catalogItemId"`
	RequestedFor    string                 `json:"requestedFor"`
	BusinessGroupID string                 `json:"businessGroupId"`
	Description     string                 `json:"description"`
	Reasons         string                 `json:"reasons"`
	Data            map[string]interface{} `json:"data"`
}

//catalogName - This struct holds catalog name from json response.
type catalogName struct {
	Name string `json:"name"`
}

//CatalogItem - This struct holds the value of response of catalog item list
type CatalogItem struct {
	CatalogItem catalogName `json:"catalogItem"`
}

//GetCatalogItemRequestTemplate - Call to retrieve a request template for a catalog item.
func (c *APIClient) GetCatalogItemRequestTemplate(catalogItemId string) (*CatalogItemRequestTemplate, error) {
	//Form a path to read catalog request template via REST call
	path := fmt.Sprintf("/catalog-service/api/consumer/entitledCatalogItems/"+
		"%s/requests/template",
		catalogItemId)

	log.Printf("GetCatalogItemRequestTemplate->path %v\n", path)

	requestTemplate := new(CatalogItemRequestTemplate)
	apiError := new(APIError)
	//Make the REST call to get catalog request template
	_, err := c.HTTPClient.New().Get(path).Receive(requestTemplate, apiError)

	if err != nil {
		return nil, err
	}

	if !apiError.isEmpty() {
		return nil, apiError
	}
	//Return catalog item template
	log.Printf("GetCatalogItemRequestTemplate->requestTemplate %v\n", requestTemplate)
	return requestTemplate, nil
}

type entitledCatalogItemViews struct {
	Links    interface{} `json:"links"`
	Content  interface{} `json:"content"`
	Metadata Metadata    `json:"metadata"`
}

// Metadata - Metadata  used to store metadata of resource list response
type Metadata struct {
	TotalElements int `json:"totalElements"`
}

// readCatalogItemNameByID - This function returns the catalog item name using catalog item ID
func (c *APIClient) readCatalogItemNameByID(catalogItemID string) (interface{}, error) {
	//Form a path to read catalog template via REST call
	path := fmt.Sprintf("/catalog-service/api/consumer/entitledCatalogItems/"+
		"%s", catalogItemID)

	template := new(CatalogItem)
	apiError := new(APIError)
	//Make a REST call to get catalog template
	_, err := c.HTTPClient.New().Get(path).Receive(template, apiError)

	if err != nil {
		return nil, err
	}

	if !apiError.isEmpty() {
		return nil, apiError
	}
	//Return catalog Name
	return template.CatalogItem.Name, nil
}

//readCatalogItemIdByName - To read id of catalog from vRA using catalog_name
func (c *APIClient) readCatalogItemIDByName(catalogName string) (interface{}, error) {
	var catalogItemID string

	log.Printf("readCatalogItemIdByName->catalog_name %v\n", catalogName)

	//Set a call to read number of catalogs from vRA
	path := fmt.Sprintf("catalog-service/api/consumer/entitledCatalogItemViews")

	template := new(entitledCatalogItemViews)
	apiError := new(APIError)

	_, preErr := c.HTTPClient.New().Get(path).Receive(template, apiError)

	if preErr != nil {
		return nil, preErr
	}

	if !apiError.isEmpty() {
		return nil, apiError
	}

	//Fetch all catalogs from vRA
	path = fmt.Sprintf("catalog-service/api/consumer/entitledCatalogItemViews?page=1&"+
		"limit=%d", template.Metadata.TotalElements)
	resp, errResp := c.HTTPClient.New().Get(path).Receive(template, apiError)

	if !apiError.isEmpty() {
		return nil, apiError
	}

	if resp.StatusCode != 200 {
		return nil, errResp
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
				return catalogItem["catalogItemId"].(interface{}), nil
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
		fmt.Println(errorMessage)
		punctuation := "is"
		if len(catalogItemNameArray) > 1 {
			punctuation = "are"
		}
		return nil, fmt.Errorf("There %s total %d catalog(s) present with same name.\n%s\n"+
			"Please select from above.", punctuation, len(catalogItemNameArray), errorMessage)
	}

	if !apiError.isEmpty() {
		return nil, apiError
	}
	return catalogItemID, nil
}
