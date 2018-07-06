provider "vra7" {
  username = "${var.username}"
  password = "${var.password}"
  tenant   = "${var.tenant}"
  host     = "${var.host}"
}

resource "vra7_resource" "vm" {
  catalog_name = "${var.catalog_name}"

  count = "${var.count}"

  catalog_configuration = {
    VirtualMachine.Disk1.Size = "${var.extra_disk}"
  }

  resource_configuration {
    Machine.description = "${var.description}"
    Machine.cpu         = "${var.cpu}"
    Machine.memory      = "${var.ram}"

    Machine.ip_address = ""
  }

  deployment_configuration = {
    reasons     = "${var.description}"
    description = "${var.description}"
  }

  wait_timeout = "${var.wait_timeout}"

  // Connection settings
  connection {
    host     = "${self.resource_configuration.Machine.ip_address}"
    user     = "${var.ssh_user}"
    password = "${var.ssh_password}"
  }

  // Extend volume to second disk
  provisioner "remote-exec" {
    inline = [
      "pvcreate /dev/sdb",
      "vgextend VolGroup00 /dev/sdb",
      "lvextend -l +100%FREE /dev/mapper/VolGroup00-rootLV",
      "resize2fs /dev/mapper/VolGroup00-rootLV",
    ]
  }
}
