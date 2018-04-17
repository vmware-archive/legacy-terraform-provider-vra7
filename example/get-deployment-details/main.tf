provider  "vra7" {
    username = "${var.username}"
    password  = "${var.password}"
    tenant = "${var.tenant}"
    host = "${var.host}"
}

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
