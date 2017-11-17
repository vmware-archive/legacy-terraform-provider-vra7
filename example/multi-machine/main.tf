provider  "vra7" {
    username = "${var.username}"
    password  = "${var.password}"
    tenant = "${var.tenant}"
    host = "${var.host}"
}

//Catalog "composite_catalog" coltains Linux, Windows and http (apache) designs.
resource "vra7_resource" "machine" {
  count            = 1
  catalog_name = "composite_catalog"
  resource_configuration = {
    Windows.cpu = "2"                //Windows Machine CPU
    Linux.cpu = "2"                  //Linux Machine CPU
    http.hostname = "xyz.com"        //HTTP (apache) hostname
    http.network_mode = "bridge"     //HTTP (apache) network mode
  }
}

