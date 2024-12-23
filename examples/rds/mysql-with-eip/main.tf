data "sbercloud_availability_zones" "myaz" {}

resource "sbercloud_vpc" "myvpc" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "sbercloud_vpc_subnet" "mysubnet" {
  vpc_id      = sbercloud_vpc.myvpc.id
  name        = var.subnet_name
  cidr        = var.subnet_cidr
  gateway_ip  = var.subnet_gateway
  primary_dns = var.primary_dns
}

resource "sbercloud_networking_secgroup" "mysecgroup" {
  name        = "mysecgroup"
  description = "a basic security group"
}

resource "sbercloud_networking_secgroup_rule" "allow_rds" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 3306
  port_range_max    = 3306
  remote_ip_prefix  = var.allow_cidr
  security_group_id = sbercloud_networking_secgroup.mysecgroup.id
}

resource "sbercloud_rds_instance" "myinstance" {
  name                = "mysql_instance"
  flavor              = "rds.mysql.c2.large.ha"
  ha_replication_mode = "async"
  vpc_id              = sbercloud_vpc.myvpc.id
  subnet_id           = sbercloud_vpc_subnet.mysubnet.id
  security_group_id   = sbercloud_networking_secgroup.mysecgroup.id
  availability_zone = [
    data.sbercloud_availability_zones.myaz.names[0],
    data.sbercloud_availability_zones.myaz.names[1]
  ]

  db {
    type     = "MySQL"
    version  = "8.0"
    password = var.password
  }
  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

resource "sbercloud_vpc_eip" "myeip" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

# get the port of rds instance by private_ip
data "sbercloud_networking_port" "rds_port" {
  network_id = sbercloud_vpc_subnet.mysubnet.id
  fixed_ip   = sbercloud_rds_instance.myinstance.private_ips[0]
}

resource "sbercloud_networking_eip_associate" "associated" {
  public_ip = sbercloud_vpc_eip.myeip.address
  port_id   = data.sbercloud_networking_port.rds_port.id
}
