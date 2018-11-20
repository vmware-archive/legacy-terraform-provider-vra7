package utils

const (
	CATALOG_NAME             = "catalog_name"
	CATALOG_ID               = "catalog_id"
	BUSINESS_GROUP_ID        = "businessgroup_id"
	BUSINESS_GROUP_NAME      = "businessgroup_name"
	WAIT_TIME_OUT            = "wait_timeout"
	FAILED_MESSAGE           = "failed_message"
	DEPLOYMENT_CONFIGURATION = "deployment_configuration"
	RESOURCE_CONFIGURATION   = "resource_configuration"
	CATALOG_CONFIGURATION    = "catalog_configuration"
	REQUEST_STATUS           = "request_status"
	INFRASTRUCTURE_VIRTUAL   = "Infrastructure.Virtual"
	DEPLOYMENT_RESOURCE_TYPE = "composition.resource.type.deployment"
	COMPONENT                = "Component"
	RECONFIGURE              = "Reconfigure"
	DESTROY                  = "Destroy"

	// read resource machine constants
	MACHINE_CPU              = "cpu"
	MACHINE_STORAGE          = "storage"
	MACHINE_MEMORY           = "memory"
	IP_ADDRESS               = "ip_address"
	MACHINE_NAME             = "name"
	MACHINE_GUEST_OS         = "guest_operating_system"
	MACHINE_BP_NAME          = "blueprint_name"
	MACHINE_TYPE             = "type"
	MACHINE_RESERVATION_NAME = "reservation_name"
	MACHINE_INTERFACE_TYPE   = "interface_type"
	MACHINE_ID               = "id"
	MACHINE_GROUP_NAME       = "group_name"
	MACHINE_DESTRUCTION_DATE = "destruction_date"
	MACHINE_RECONFIGURE      = "reconfigure"
	MACHINE_POWER_OFF        = "power_off"

	// utility constants
	LOGGER_ID              = "terraform-provider-vra7"
	WINDOWS_PATH_SEPARATOR = "\\"
	UNIX_PATH_SEPARATOR    = "/"
	WINDOWS_OS             = "windows"
	IN_PROGRESS            = "IN_PROGRESS"
	SUCCESSFUL             = "SUCCESSFUL"
	FAILED                 = "FAILED"

	// error constants
	CONFIG_INVALID_ERROR          = "The resource_configuration in the config file has invalid component name(s): %v "
	DESTROY_ACTION_TEMPLATE_ERROR = "Error retrieving destroy action template for the deployment %v: %v "

	// api constants
	CATALOG_SERVICE               = "/catalog-service"
	CATALOG_SERVICE_API           = CATALOG_SERVICE + "/api"
	CONSUMER                      = CATALOG_SERVICE_API + "/consumer"
	CONSUMER_REQUESTS             = CONSUMER + "/requests"
	CONSUMER_RESOURCES            = CONSUMER + "/resources"
	GET_RESOURCE_API              = CONSUMER_REQUESTS + "/" + "%s" + "/resources"
	POST_ACTION_TEMPLATE_API      = CONSUMER_RESOURCES + "/" + "%s" + "/actions/" + "%s" + "/requests"
	GET_ACTION_TEMPLATE_API       = POST_ACTION_TEMPLATE_API + "/template"
	GET_REQUEST_RESOURCE_VIEW_API = CONSUMER_REQUESTS + "/" + "%s" + "/resourceViews"
)
