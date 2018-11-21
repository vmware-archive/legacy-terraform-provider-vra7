package vrealize

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/terraform-provider-vra7/utils"
)

//ResourceActionTemplate - is used to store information
//related to resource action template information.
type ResourceActionTemplate struct {
	Type        string                 `json:"type,omitempty"`
	ResourceID  string                 `json:"resourceId,omitempty"`
	ActionID    string                 `json:"actionId,omitempty"`
	Description string                 `json:"description,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

//ResourceView - is used to store information
//related to resource template information.
type ResourceView struct {
	Content []interface {
	} `json:"content"`
	Links []interface{} `json:"links"`
}

//RequestStatusView - used to store REST response of
//request triggered against any resource.
type RequestStatusView struct {
	RequestCompletion struct {
		RequestCompletionState string `json:"requestCompletionState"`
		CompletionDetails      string `json:"CompletionDetails"`
	} `json:"requestCompletion"`
	Phase string `json:"phase"`
}

type BusinessGroups struct {
	Content []BusinessGroup `json:"content,omitempty"`
}

type BusinessGroup struct {
	Name string `json:"name,omitempty"`
	Id   string `json:"id,omitempty"`
}

// Resource View of a provisioned request
type RequestResourceView struct {
	Content []DeploymentResource `json:"content,omitempty"`
	Links   []interface{}        `json:"links,omitempty"`
}

type DeploymentResource struct {
	RequestState    string                 `json:"requestState,omitempty"`
	Description     string                 `json:"description,omitempty"`
	LastUpdated     string                 `json:"lastUpdated,omitempty"`
	TenantId        string                 `json:"tenantId,omitempty"`
	Name            string                 `json:"name,omitempty"`
	BusinessGroupId string                 `json:"businessGroupId,omitempty"`
	DateCreated     string                 `json:"dateCreated,omitempty"`
	Status          string                 `json:"status,omitempty"`
	RequestId       string                 `json:"requestId,omitempty"`
	ResourceId      string                 `json:"resourceId,omitempty"`
	ResourceType    string                 `json:"resourceType,omitempty"`
	ResourcesData   DeploymentResourceData `json:"data,omitempty"`
}

type DeploymentResourceData struct {
	Memory                      int    `json:"MachineMemory,omitempty"`
	Cpu                         int    `json:"MachineCPU,omitempty"`
	IpAddress                   string `json:"ip_address,omitempty"`
	Storage                     int    `json:"MachineStorage,omitempty"`
	MachineInterfaceType        string `json:"MachineInterfaceType,omitempty"`
	MachineName                 string `json:"MachineName,omitempty"`
	MachineGuestOperatingSystem string `json:"MachineGuestOperatingSystem,omitempty"`
	MachineDestructionDate      string `json:"MachineDestructionDate,omitempty"`
	MachineGroupName            string `json:"MachineGroupName,omitempty"`
	MachineBlueprintName        string `json:"MachineBlueprintName,omitempty"`
	MachineReservationName      string `json:"MachineReservationName,omitempty"`
	MachineType                 string `json:"MachineType,omitempty"`
	MachineId                   string `json:"machineId,omitempty"`
	MachineExpirationDate       string `json:"MachineExpirationDate,omitempty"`
	Component                   string `json:"Component,omitempty"`
	Expire                      bool   `json:"Expire,omitempty"`
	Reconfigure                 bool   `json:"Reconfigure,omitempty"`
	Reset                       bool   `json:"Reset,omitempty"`
	Reboot                      bool   `json:"Reboot,omitempty"`
	PowerOff                    bool   `json:"PowerOff,omitempty"`
	Destroy                     bool   `json:"Destroy,omitempty"`
	Shutdown                    bool   `json:"Shutdown,omitempty"`
	Suspend                     bool   `json:"Suspend,omitempty"`
	Reprovision                 bool   `json:"Reprovision,omitempty"`
	ChangeLease                 bool   `json:"ChangeLease,omitempty"`
	ChangeOwner                 bool   `json:"ChangeOwner,omitempty"`
	CreateSnapshot              bool   `json:"CreateSnapshot,omitempty"`
}

// Retrieves the resources that were provisioned as a result of a given request.
// Also returns the actions allowed on the resources and their templates
type ResourceActions struct {
	Links   []interface{}           `json:"links,omitempty"`
	Content []ResourceActionContent `json:"content,omitempty"`
}

type ResourceActionContent struct {
	Id              string          `json:"id,omitempty"`
	Name            string          `json:"name,omitempty"`
	ResourceTypeRef ResourceTypeRef `json:"resourceTypeRef,omitempty"`
	Status          string          `json:"status,omitempty"`
	RequestId       string          `json:"requestId,omitempty"`
	RequestState    string          `json:"requestState,omitempty"`
	Operations      []Operation     `json:"operations,omitempty"`
	ResourceData    ResourceDataMap `json:"resourceData,omitempty"`
}

type ResourceTypeRef struct {
	Id    string `json:"id,omitempty"`
	Label string `json:"label,omitempty"`
}

type Operation struct {
	Name        string `json:"name,omitempty"`
	OperationId string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
}

type ResourceDataMap struct {
	Entries []ResourceDataEntry `json:"entries,omitempty"`
}
type ResourceDataEntry struct {
	Key   string                 `json:"key,omitempty"`
	Value map[string]interface{} `json:"value,omitempty"`
}

//CatalogRequest - A structure that captures a vRA catalog request.
type CatalogRequest struct {
	ID           string      `json:"id"`
	IconID       string      `json:"iconId"`
	Version      int         `json:"version"`
	State        string      `json:"state"`
	Description  string      `json:"description"`
	Reasons      interface{} `json:"reasons"`
	RequestedFor string      `json:"requestedFor"`
	RequestedBy  string      `json:"requestedBy"`
	Organization struct {
		TenantRef      string `json:"tenantRef"`
		TenantLabel    string `json:"tenantLabel"`
		SubtenantRef   string `json:"subtenantRef"`
		SubtenantLabel string `json:"subtenantLabel"`
	} `json:"organization"`

	RequestorEntitlementID   string                 `json:"requestorEntitlementId"`
	PreApprovalID            string                 `json:"preApprovalId"`
	PostApprovalID           string                 `json:"postApprovalId"`
	DateCreated              time.Time              `json:"dateCreated"`
	LastUpdated              time.Time              `json:"lastUpdated"`
	DateSubmitted            time.Time              `json:"dateSubmitted"`
	DateApproved             time.Time              `json:"dateApproved"`
	DateCompleted            time.Time              `json:"dateCompleted"`
	Quote                    interface{}            `json:"quote"`
	RequestData              map[string]interface{} `json:"requestData"`
	RequestCompletion        string                 `json:"requestCompletion"`
	RetriesRemaining         int                    `json:"retriesRemaining"`
	RequestedItemName        string                 `json:"requestedItemName"`
	RequestedItemDescription string                 `json:"requestedItemDescription"`
	Components               string                 `json:"components"`
	StateName                string                 `json:"stateName"`

	CatalogItemProviderBinding struct {
		BindingID   string `json:"bindingId"`
		ProviderRef struct {
			ID    string `json:"id"`
			Label string `json:"label"`
		} `json:"providerRef"`
	} `json:"catalogItemProviderBinding"`

	Phase           string `json:"phase"`
	ApprovalStatus  string `json:"approvalStatus"`
	ExecutionStatus string `json:"executionStatus"`
	WaitingStatus   string `json:"waitingStatus"`
	CatalogItemRef  struct {
		ID    string `json:"id"`
		Label string `json:"label"`
	} `json:"catalogItemRef"`
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
		utils.CATALOG_NAME: {
			Type:     schema.TypeString,
			Optional: true,
		},
		utils.CATALOG_ID: {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		utils.BUSINESS_GROUP_ID: {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		utils.BUSINESS_GROUP_NAME: {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		utils.WAIT_TIME_OUT: {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  15,
		},
		utils.REQUEST_STATUS: {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
		utils.FAILED_MESSAGE: {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
			Optional: true,
		},
		utils.DEPLOYMENT_CONFIGURATION: {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
		},
		utils.RESOURCE_CONFIGURATION: {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
		},
		utils.CATALOG_CONFIGURATION: {
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
