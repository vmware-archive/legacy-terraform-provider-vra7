provider  "vra7" {
    username = "${var.username}"
    password  = "${var.password}"
    tenant = "${var.tenant}"
    host = "${var.host}"
}

resource "vra7_resource" "machine1" {
  count            = 1
  catalog_name = "CentOS 7.0 x64"
}

resource "vra7_resource" "machine2" {
  count            = 1
  catalog_name = "CentOS 7.0 x64"
  depends_on = ["vra7_resource.machine1"]
}
