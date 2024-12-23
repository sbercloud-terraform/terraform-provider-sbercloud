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
    password = var.rds_password
  }
  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

resource "sbercloud_rds_read_replica_instance" "myreplica" {
  name                = "myreplica"
  flavor              = "rds.mysql.c2.large.rr"
  primary_instance_id = sbercloud_rds_instance.myinstance.id
  availability_zone   = data.sbercloud_availability_zones.myaz.names[1]
  volume {
    type = "CLOUDSSD"
  }

  tags = {
    type = "readonly"
  }
}
