package vrealize

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
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

//Function use - To authenticate terraform provider
func providerConfig(r *schema.ResourceData) (interface{}, error) {
	//Create a client handle to perform REST calls for various operations upon the resource
	client := NewClient(r.Get("username").(string),
		r.Get("password").(string),
		r.Get("tenant").(string),
		r.Get("host").(string),
		r.Get("insecure").(bool),
	)

	//Authenticate user
	err := client.Authenticate()

	//Raise an error on authentication fail
	if err != nil {
		return nil, fmt.Errorf("Error: Unable to get auth token: %v", err)
	}

	//Return client handle on success
	return &client, nil
}

//Function use - set machine resource details based on machine type
func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"vra7_resource": ResourceMachine(),
	}
}
