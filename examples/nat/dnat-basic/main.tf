data "sbercloud_availability_zones" "newAZ_Example" {}

data "sbercloud_images_image" "newIMS_Example" {
  name        = var.ims_name
  visibility  = "public"
  most_recent = true
}


resource "sbercloud_compute_instance" "newCompute_Example" {
  name              = var.ecs_name
  image_id          = data.sbercloud_images_image.newIMS_Example.id
  flavor_id         = "s6.small.1"
  security_groups   = [sbercloud_networking_secgroup.newSecgroup_Example.name]
  admin_pass        = var.password
  availability_zone = data.sbercloud_availability_zones.newAZ_Example.names[0]

  system_disk_type = "SSD"
  system_disk_size = 40

  network {
    uuid = sbercloud_vpc_subnet.newSubnet_Example.id
  }
}

resource "sbercloud_vpc_eip" "newEIP_Example" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = var.bandwidth_name
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "sbercloud_vpc" "newVPC_Example" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "sbercloud_vpc_subnet" "newSubnet_Example" {
  name          = var.subnet_name
  cidr          = var.subnet_cidr
  gateway_ip    = var.subnet_gateway_ip
  vpc_id        = sbercloud_vpc.newVPC_Example.id
  primary_dns   = "100.125.13.59"
  secondary_dns = "100.125.65.14"
}

resource "sbercloud_networking_secgroup" "newSecgroup_Example" {
  name        = var.secgroup_name
  description = "This is a security group"
}

resource "sbercloud_networking_secgroup_rule" "newSecgroupRule_Example" {
  count = length(var.security_group_rule)

  direction         = lookup(var.security_group_rule[count.index], "direction", null)
  ethertype         = lookup(var.security_group_rule[count.index], "ethertype", null)
  protocol          = lookup(var.security_group_rule[count.index], "protocol", null)
  port_range_min    = lookup(var.security_group_rule[count.index], "port_range_min", null)
  port_range_max    = lookup(var.security_group_rule[count.index], "port_range_max", null)
  remote_ip_prefix  = lookup(var.security_group_rule[count.index], "remote_ip_prefix", null)
  security_group_id = sbercloud_networking_secgroup.newSecgroup_Example.id
}

resource "sbercloud_nat_gateway" "newNet_gateway_Example" {
  name                = var.net_gateway_name
  description         = "example for net test"
  spec                = "1"
  vpc_id           = sbercloud_vpc.newVPC_Example.id
  subnet_id = sbercloud_vpc_subnet.newSubnet_Example.id
}

resource "sbercloud_nat_snat_rule" "newSNATRule_Example" {
  nat_gateway_id = sbercloud_nat_gateway.newNet_gateway_Example.id
  subnet_id     = sbercloud_vpc_subnet.newSubnet_Example.id
  floating_ip_id = sbercloud_vpc_eip.newEIP_Example.id
}

resource "sbercloud_nat_dnat_rule" "newDNATRule_Example" {
  count = length(var.example_dnat_rule)

  floating_ip_id = sbercloud_vpc_eip.newEIP_Example.id
  nat_gateway_id = sbercloud_nat_gateway.newNet_gateway_Example.id
  port_id        = sbercloud_compute_instance.newCompute_Example.network[0].port

  internal_service_port = lookup(var.example_dnat_rule[count.index], "internal_service_port", null)
  protocol              = lookup(var.example_dnat_rule[count.index], "protocol", null)
  external_service_port = lookup(var.example_dnat_rule[count.index], "external_service_port", null)
}