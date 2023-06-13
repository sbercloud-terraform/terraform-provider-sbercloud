data "sbercloud_availability_zones" "myaz" {}

data "sbercloud_compute_flavors" "myflavor" {
  availability_zone = data.sbercloud_availability_zones.myaz.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "sbercloud_images_image" "myimage" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "sbercloud_vpc" "vpc_1" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "sbercloud_vpc_subnet" "subnet_1" {
  vpc_id      = sbercloud_vpc.vpc_1.id
  name        = var.subnet_name
  cidr        = var.subnet_cidr
  gateway_ip  = var.subnet_gateway
  primary_dns = var.primary_dns
}

resource "sbercloud_compute_instance" "mycompute" {
  name              = "mycompute_${count.index}"
  image_id          = data.sbercloud_images_image.myimage.id
  flavor_id         = data.sbercloud_compute_flavors.myflavor.ids[0]
  security_groups   = ["default"]
  availability_zone = data.sbercloud_availability_zones.myaz.names[0]

  network {
    uuid = sbercloud_vpc_subnet.subnet_1.id
  }
  count = 2
}

resource "sbercloud_networking_vip" "vip_1" {
  network_id = sbercloud_vpc_subnet.subnet_1.id
}

# associate ports to the vip
resource "sbercloud_networking_vip_associate" "vip_associated" {
  vip_id   = sbercloud_networking_vip.vip_1.id
  port_ids = [
    sbercloud_compute_instance.mycompute[0].network[0].port,
    sbercloud_compute_instance.mycompute[1].network[0].port
  ]
}
