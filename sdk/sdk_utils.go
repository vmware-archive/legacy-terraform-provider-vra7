package sdk

// API constants
const (
	CatalogService              = "/catalog-service"
	CatalogServiceAPI           = CatalogService + "/api"
	Consumer                    = CatalogServiceAPI + "/consumer"
	ConsumerRequests            = Consumer + "/requests"
	ConsumerResources           = Consumer + "/resources"
	EntitledCatalogItems        = Consumer + "/entitledCatalogItems"
	EntitledCatalogItemViewsAPI = Consumer + "/entitledCatalogItemViews"
	GetResourceAPI              = ConsumerRequests + "/" + "%s" + "/resources"
	PostActionTemplateAPI       = ConsumerResources + "/" + "%s" + "/actions/" + "%s" + "/requests"
	GetActionTemplateAPI        = PostActionTemplateAPI + "/template"
	GetRequestResourceViewAPI   = ConsumerRequests + "/" + "%s" + "/resourceViews"
	RequestTemplateAPI          = EntitledCatalogItems + "/" + "%s" + "/requests/template"
)
