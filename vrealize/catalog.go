package vrealize

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vmware/terraform-provider-vra7/client"
	"github.com/vmware/terraform-provider-vra7/utils"
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
	ID   string `json:"catalogItemId"`
}

//CatalogItem - This struct holds the value of response of catalog item list
type CatalogItem struct {
	CatalogItem catalogName `json:"catalogItem"`
}

//GetCatalogItemRequestTemplate - Call to retrieve a request template for a catalog item.
func GetCatalogItemRequestTemplate(catalogItemID string) (*CatalogItemRequestTemplate, error) {

	//Form a path to read catalog request template via REST call
	path := fmt.Sprintf("/catalog-service/api/consumer/entitledCatalogItems/"+
		"%s/requests/template",
		catalogItemID)
	url := client.BuildEncodedURL(path, nil)
	respBody, respErr := client.Get(url, nil)
	if respErr != nil {
		return nil, respErr
	}

	var requestTemplate CatalogItemRequestTemplate
	unmarshallErr := utils.UnmarshalJSON(respBody.Body, &requestTemplate)
	if unmarshallErr != nil {
		return nil, unmarshallErr
	}
	return &requestTemplate, nil
}

type EntitledCatalogItemViews struct {
	Links    interface{} `json:"links"`
	Content  interface{} `json:"content"`
	Metadata Metadata    `json:"metadata"`
}

// Metadata - Metadata  used to store metadata of resource list response
type Metadata struct {
	TotalElements int `json:"totalElements"`
}

// readCatalogItemNameByID - This function returns the catalog item name using catalog item ID
func ReadCatalogItemNameByID(catalogItemID string) (string, error) {

	path := fmt.Sprintf("/catalog-service/api/consumer/entitledCatalogItems/"+
		"%s", catalogItemID)
	url := client.BuildEncodedURL(path, nil)
	respBody, respErr := client.Get(url, nil)
	if respErr != nil {
		return "", respErr
	}

	var response CatalogItem
	unmarshallErr := utils.UnmarshalJSON(respBody.Body, &response)
	if unmarshallErr != nil {
		return "", unmarshallErr
	}
	return response.CatalogItem.Name, nil
}

// ReadCatalogItemIdByName - To read id of catalog from vRA using catalog_name
// func ReadCatalogItemByName(catalogName string) (string, error) {
//
// 	path := "/catalog-service/api/consumer/entitledCatalogItemViews"
// 	log.Info("Fetching business group id from name..GET %s ", path)
// 	uri := client.BuildEncodedURL(path+"?$filter=name eq "+catalogName, nil)
// 	customURL := strings.Replace(uri, "%3F", "?", -1)
//
// 	respBody, respErr := client.Get(customURL, nil)
// 	if respErr != nil {
// 		return "", respErr
// 	}
//
// 	var response CatalogItem
// 	unmarshallErr := utils.UnmarshalJSON(respBody, &response)
// 	if unmarshallErr != nil {
// 		return "", unmarshallErr
// 	}
// 	log.Info("the catalog item is %v ", response.CatalogItem.ID)
// 	return response.CatalogItem.ID, nil
// }

//readCatalogItemIdByName - To read id of catalog from vRA using catalog_name
func ReadCatalogItemByName(catalogName string) (string, error) {
	var catalogItemID string

	log.Info("readCatalogItemIdByName->catalog_name %v\n", catalogName)

	//Set a call to read number of catalogs from vRA
	path := fmt.Sprintf("catalog-service/api/consumer/entitledCatalogItemViews")

	url := client.BuildEncodedURL(path, nil)
	respBody, respErr := client.Get(url, nil)
	if respErr != nil {
		return "", respErr
	}
	if respBody.StatusCode != 200 {
		return "", fmt.Errorf("Error with status code %v", respBody.StatusCode)
	}

	var template EntitledCatalogItemViews
	unmarshallErr := utils.UnmarshalJSON(respBody.Body, &template)
	if unmarshallErr != nil {
		return "", unmarshallErr
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
				return catalogItem["catalogItemId"].(string), nil
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
		return "", fmt.Errorf("There %s total %d catalog(s) present with same name.\n%s\n"+
			"Please select from above.", punctuation, len(catalogItemNameArray), errorMessage)
	}
	return catalogItemID, nil
}
