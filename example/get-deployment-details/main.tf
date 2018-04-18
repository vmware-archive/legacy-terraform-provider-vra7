provider  "vra7" {
    username = "${var.username}"
    password  = "${var.password}"
    tenant = "${var.tenant}"
    host = "${var.host}"
}

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
    CentOS_6.3.cpu = ""
    CentOS_6.3.ip_address = ""
  }
}
output "machine_ip" {
  value = "${vra7_resource.CentOS_machine.resource_configuration.CentOS_6.3.ip_address}"
}

output "machine_cpu" {
  value = "${vra7_resource.CentOS_machine.resource_configuration.CentOS_6.3.cpu}"
}