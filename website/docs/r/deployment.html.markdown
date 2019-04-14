---
layout: "vra7"
page_title: "VMware vRA7: vra7_deployment"
sidebar_current: "docs-vra7-resource-deployment"
description: |-
  Provides a VMware vRA7 deployment resource. This can be used to deploy vRA7 catalog items.
---

# vra7\_deployment

Provides a VMware vRA7 deployment resource. This can be used to deploy vRA7 catalog items.

## Example Usages

**Simple deployment of a CentOS Linux host with 2 CPU's:**

```hcl
resource "vra7_deployment" "machine" {
  count            = 1
  catalog_item_name = "CentOS 7.0 x64"
  resource_configuration = {
    Linux.cpu = "2"
  }
}
```

**Catalog "multi_machine_catalog" contains Linux, Windows and http (apache) designs:**

```hcl
resource "vra7_deployment" "resource_1" {
  count            = 1
  catalog_item_name = "multi_machine_catalog"
  resource_configuration = {
    Windows.cpu = "2"                //Windows Machine CPU
    Linux.cpu = "2"                  //Linux Machine CPU
    http.hostname = "xyz.com"        //HTTP (apache) hostname
    http.network_mode = "bridge"     //HTTP (apache) network mode
  }
}
```

**Showing the use of `depends_on` between two deployments:**

```hcl
resource "vra7_deployment" "machine1" {
  count            = 1
  catalog_item_name = "CentOS 7.0 x64"
}

resource "vra7_deployment" "machine2" {
  count            = 1
  catalog_item_name = "CentOS 7.0 x64"
  depends_on = ["vra7_deployment.machine1"]
}
```

## Argument Reference

The following arguments are supported:

* `businessgroup_id` - (Optional) The id of the vRA business group to use for this deployment.
* `businessgroup_name` - (Optional) The name of the vRA business group to use for this deployment.
* `catalog_item_id` - (Optional) The id of the catalog item to deploy into vRA.
* `catalog_item_name` - (Optional) The name of the catalog item to deploy into vRA.
* `description` - (Optional) Description of the deployment
* `resource_configuration` - (Optional) The configuration of the individual components from the catalog item
