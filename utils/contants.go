package utils

// terraform provider constants
const (
	CatalogName             = "catalog_name"
	CatalogID               = "catalog_id"
	BusinessGroupID         = "businessgroup_id"
	BusinessGroupName       = "businessgroup_name"
	WaitTimeout             = "wait_timeout"
	FailedMessage           = "failed_message"
	DeploymentConfiguration = "deployment_configuration"
	ResourceConfiguration   = "resource_configuration"
	CatalogConfiguration    = "catalog_configuration"
	RequestStatus           = "request_status"

	// read resource machine constants

	MachineCPU             = "cpu"
	MachineStorage         = "storage"
	MachineMemory          = "memory"
	IPAddress              = "ip_address"
	MachineName            = "name"
	MachineGuestOs         = "guest_operating_system"
	MachineBpName          = "blueprint_name"
	MachineType            = "type"
	MachineReservationName = "reservation_name"
	MachineInterfaceType   = "interface_type"
	MachineID              = "id"
	MachineGroupName       = "group_name"
	MachineDestructionDate = "destruction_date"
	MachineReconfigure     = "reconfigure"
	MachinePowerOff        = "power_off"

	// utility constants

	LoggerID               = "terraform-provider-vra7"
	WindowsPathSeparator   = "\\"
	UnixPathSeparator      = "/"
	WindowsOs              = "windows"
	InProgress             = "IN_PROGRESS"
	Successful             = "SUCCESSFUL"
	Failed                 = "FAILED"
	Submitted              = "SUBMITTED"
	InfrastructureVirtual  = "Infrastructure.Virtual"
	DeploymentResourceType = "composition.resource.type.deployment"
	Component              = "Component"
	Reconfigure            = "Reconfigure"
	Destroy                = "Destroy"

	// error constants

	ConfigInvalidError                = "The resource_configuration in the config file has invalid component name(s): %v "
	DestroyActionTemplateError        = "Error retrieving destroy action template for the deployment %v: %v "
	BusinessGroupIDNameNotMatchingErr = "The business group name %s and id %s does not belong to the same business group, provide either name or id"
	CatalogItemIDNameNotMatchingErr   = "The catalog item name %s and id %s does not belong to the same catalog item, provide either name or id"
	// api constants

	CatalogService            = "/catalog-service"
	CatalogServiceAPI         = CatalogService + "/api"
	Consumer                  = CatalogServiceAPI + "/consumer"
	ConsumerRequests          = Consumer + "/requests"
	ConsumerResources         = Consumer + "/resources"
	GetResourceAPI            = ConsumerRequests + "/" + "%s" + "/resources"
	PostActionTemplateAPI     = ConsumerResources + "/" + "%s" + "/actions/" + "%s" + "/requests"
	GetActionTemplateAPI      = PostActionTemplateAPI + "/template"
	GetRequestResourceViewAPI = ConsumerRequests + "/" + "%s" + "/resourceViews"
)
