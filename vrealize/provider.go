package vrealize

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/terraform-provider-vra7/sdk"
	"github.com/vmware/terraform-provider-vra7/utils"
)

//Provider - This function initializes the provider schema
//also the config function and resource mapping
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema:        providerSchema(),
		ConfigureFunc: providerConfig,
		ResourcesMap:  providerResources(),
	}
}

//providerSchema - To set provider fields
func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"username": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Tenant administrator username.",
		},
		"password": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Tenant administrator password.",
		},
		"tenant": {
			Type:     schema.TypeString,
			Required: true,
			Description: "Specifies the tenant URL token determined by the system administrator" +
				"when creating the tenant, for example, support.",
		},
		"host": {
			Type:     schema.TypeString,
			Required: true,
			Description: "host name.domain name of the vRealize Automation server, " +
				"for example, mycompany.mktg.mydomain.com.",
		},
		"insecure": {
			Type:        schema.TypeBool,
			Default:     false,
			Optional:    true,
			Description: "Specify whether to validate TLS certificates.",
		},
	}
}

//ResourceMachine - use to set resource fields
func ResourceMachine() *schema.Resource {
	return &schema.Resource{
		Create: createResource,
		Read:   readResource,
		Update: updateResource,
		Delete: deleteResource,
		Schema: resourceSchema(),
	}
}

//set_resource_schema - This function is used to update the catalog item template/blueprint
//and replace the values with user defined values added in .tf file.
func resourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		utils.CatalogName: {
			Type:     schema.TypeString,
			Optional: true,
		},
		utils.CatalogID: {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		utils.BusinessGroupID: {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		utils.BusinessGroupName: {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		utils.WaitTimeout: {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  15,
		},
		utils.RequestStatus: {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
		utils.FailedMessage: {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
			Optional: true,
		},
		utils.DeploymentConfiguration: {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
		},
		utils.ResourceConfiguration: {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
		},
		utils.CatalogConfiguration: {
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

//Function use - To authenticate terraform provider
func providerConfig(r *schema.ResourceData) (interface{}, error) {
	//Create a client handle to perform REST calls for various operations upon the resource

	user := r.Get("username").(string)
	password := r.Get("password").(string)
	tenant := r.Get("tenant").(string)
	baseURL := r.Get("host").(string)
	insecure := r.Get("insecure").(bool)
	vraClient := sdk.NewClient(user, password, tenant, baseURL, insecure)

	//Authenticate user
	err := vraClient.Authenticate()

	//Raise an error on authentication fail
	if err != nil {
		return nil, fmt.Errorf("Error: Unable to get auth token: %v", err)
	}

	//Return client handle on success
	return &vraClient, nil
}

//Function use - set machine resource details based on machine type
func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"vra7_resource": ResourceMachine(),
	}
}
