provider  "vra7" {
  username = "${var.username}"
  password  = "${var.password}"
  tenant = "${var.tenant}"
  host = "${var.host}"
}

resource "vra7_resource" "resource_1" {
  count            = 1
  catalog_name = "CentOS_6.3"
  resource_configuration = {
    CentOS_6.3.cpu = "2"
  }
  scale_resource = {
    CentOS_6.3 = 2
  }
}