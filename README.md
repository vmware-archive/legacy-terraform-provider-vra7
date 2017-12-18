# VMware Terraform provider for vRealize Automation 7
[![Build Status](https://travis-ci.org/vmware/terraform-provider-vra7.svg?branch=master)](https://travis-ci.org/vmware/terraform-provider-vra7)

A self-contained deployable integration between Terraform and vRealize Automation (vRA) which allows Terraform users to request/provision entitled vRA catalog items using Terraform. Supports Terraform destroying vRA provisioned resources.

## Getting Started

These instructions will get you a copy of the project up and run on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## Prerequisites

To make plugin up and running you need following things.
* [Terraform 9.8 or above](https://www.terraform.io/downloads.html)
* [Go Language 1.9.2 or above](https://golang.org/dl/)
* [dep - new dependency management tool for Go](https://github.com/golang/dep)

## Project Setup

Setup a GOLang project structure

```
|-/home/<USER>/TerraformPluginProject
    |-bin
    |-pkg
    |-pkg

```

## Environment Setup

Set following environment variables

**Linux Users**

*GOROOT is a golang library path*
```
export GOROOT=/usr/local/go
```

*GOPATH is a path pointing toward the source code directory*
```
export GOPATH=/home/<user>/TerraformPluginProject
```

**Windows Users**

*GOROOT is a golang library path*
```
set GOROOT=C:\Go
```

*GOPATH is a path pointing toward the source code directory*
```
set GOPATH=C:\TerraformPluginProject
```

## Set terraform provider

**Linux Users**

Create *~/.terraformrc* and put following content in it.
```
    providers {
         vra7 = "/home/<USER>/TerraformPluginProject/bin/terraform-provider-vra7"
    }
```

**Windows Users**

Create *%APPDATA%/terraform.rc* and put following content in it.
```
    providers {
         vra7 = "C:\\TerraformPluginProject\\bin\\terraform-provider-vra7.exe"
    }
```


## Installation
Clone repo code into go project using *go get*
```
    go get github.com/vmware/terraform-provider-vra7

```

## Create Binary

**Linux and MacOS Users**

Navigate to */home/<USER>/TerraformPluginProject/src/github.com/vmware/terraform-provider-vra7* and run go build command to generate plugin binary

```
    dep ensure
    go build -o /home/<USER>/TerraformPluginProject/bin/terraform-provider-vra7

```

**Windows Users**

Navigate to *C:\TerraformPluginProject\src\github.com\vmware\terraform-provider-vra7* and run go build command to generate plugin binary

```
    dep ensure
    go build -o C:\TerraformPluginProject\bin\terraform-provider-vra7.exe

```

## Create Terraform Configuration file

In VMware vRA terraform configuration file contains two objects

### Provider

This part contains service provider details.

**Configure Provider**

Provider block contains four mandatory fields

* **username** - *vRA portal username*
* **password** - *vRA portal password*
* **tenant** - *vRA portal tenant*
* **host** - *End point of REST API*

Example

```
    provider "vra7" {
      username = "vRAUser1@vsphere.local"
      password = "password123!"
      tenant = "corp.local.tenant"
      host = "http://myvra.example.com/"
    }

```


### Resource

This part contains any resource that can be deployed on that service provider.
For example, in our case machine blueprint, software blueprint, complex blueprint, network, etc.

**Configure Resource**

Syntax

```
resource "vra7" "<resource_name1>" {
}
```

Resource block contains two mandatory and three optional fields as follows

* **catalog_name** - *catalog_name is a mandatory field which contains valid catalog name from your vRA*

* **catalog_id** - *catalog_id is also a mandatory field but optional to catalog_name which contains valid catalog_id from your vRA. You either include catalog_name or catalo_id but one field should be present in resource configuration.*

* **businessgroup_id** - *This is an optional field. You can specify a different Business Group ID from what provided by default in the template reques, provided that your account is allowed to do it*

* **resource_configuration** - *This is optional field. If blueprint properties have default values or no mandatory property value is required then you can skip this field from terraform configuration file. This field contains user inputs to catalog services. Value of this field is in key value pair. Key is service.field_name and value is any valid user input to the respective field.*

* **catalog_configuration** - *This is an optional field. If catalog properties have default values or no mandatory user input required for catalog service then you can skip this field from terraform configuration file. This field contains user inputs to catalog services. Value of this field is in key value pair. Key is any field name of catalog and value is any valid user input to the respective field.*

* **count** - *This field is used to create replicas of resources. If count is not provided then it will be considered as 1 by default.*

Example 1

```
resource "vra7_resource" "example_machine1" {
  catalog_name = "CentOS 6.3"
   resource_configuration = {
         Linux.cpu = "1"
         Windows2008R2SP1.cpu =  "2"
         Windows2012.cpu =  "4"
         Windows2016.cpu =  "2"
     }
     catalog_configuration = {
         lease_days = "5"
     }
     count = 3
}

```

Example 2

```
resource "vra7_resource" "example_machine2" {
  catalog_id = "e5dd4fba7f96239286be45ed"
   resource_configuration = {
         Linux.cpu = "1"
         Windows2008.cpu =  "2"
         Windows2012.cpu =  "4"
         Windows2016.cpu =  "2"
     }
     count = 4
}

```

Save this configuration in main.tf in a path where the binary is placed.

## Execution

These are the terraform commands that can be used on vRA plugin as follows.
* **terraform init** - *The init command is used to initialize a working directory containing Terraform configuration files.*

* **terraform plan** - *Plan command shows plan for resources like how many resources will be provisioned and how many will be destroyed.*

* **terraform apply** - *apply is responsible to execute actual calls for provision resources.*

* **terraform refresh** - *By using refresh command you can check status of request.*

* **terraform show** - *show will set a console output for resource configuration and request status.*

* **terraform destroy** - *destroy command will destroy all the  resources present in terraform configuration file.*

Navigate to the location where main.tf and binary are placed and use above commands as needed.

## Contributing

The terraform-provider-vra7 project team welcomes contributions from the community. Before you start working with terraform-provider-vra7, please read our [Developer Certificate of Origin](https://cla.vmware.com/dco). All contributions to this repository must be signed as described on that page. Your signature certifies that you wrote the patch or have the right to pass it on as an open-source patch. For more detailed information, refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License

terraform-provider-vra7 is available under the [MIT license](LICENSE).
