provider  "vra7" {
    username = "${var.username}"
    password  = "${var.password}"
    tenant = "${var.tenant}"
    host = "${var.host}"
}

<<<<<<< HEAD
//TF resource with one centos deployment
resource "vra7_resource" "CentOS_machine" {
  count            = 1
  catalog_name = "CentOS 6.3"
  //Deployment level configuration
  resource_configuration = {
    //After successful deployment
    //CPU and IP address will get updated
    //with actual values of deployment
    //on `terraform refresh'
    CentOs.cpu = ""
    Win.ip_address = ""
  }
}
output "centOsMachineCPU" {
  value = "${vra7_resource.CentOS_machine.resource_configuration.CentOs.cpu}"
}

output "WinMachineAddress" {
  value = "${vra7_resource.CentOS_machine.resource_configuration.Win.ip_address}"
}
=======
resource "vra7_resource" "multi_machine_cl" {
  count            = 1
  catalog_name = "CentOS"
  resource_configuration = {
    Linux.0.cpu = ""
    Linux.0.memory = ""
    Linux.1.cpu = ""
    Linux.1.memory = ""
  }
}
>>>>>>> d2f1a81... Fetch deployment details example
