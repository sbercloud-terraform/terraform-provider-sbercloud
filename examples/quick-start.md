## Quick Start
Example for fast jump-start without deep dive into details.

Creates:
* **1 VPC** with **1 public IPv4**
* **1 VM** with Ubuntu 20.04 and binded public IPv4
* **1 Managed PostgreSQL 12** with access from VM by LAN
* ... and a little bit service entities

### Network / VPC

```terraform
resource "sbercloud_vpc" "vpc_quickstart" {
  name = "vpc_quickstart"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_subnet" "subnet_quickstart" {

  name = "subnet_quickstart"
  
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"

  primary_dns   = "100.125.13.59"
  secondary_dns = "1.1.1.1"

  vpc_id = sbercloud_vpc.vpc_quickstart.id

}
```

### Network / Public IP

```terraform
resource "sbercloud_vpc_eip" "eip_quickstart" {
  
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "elb_bandwidth"
    size        = 100
    share_type  = "PER"
    charge_mode = "bandwidth"
  }

}
```

### Security group and basic rules

```terraform
resource "sbercloud_networking_secgroup" "secgroup_quickstart" {
  name        = "secgroup_quickstart"
  description = "quickstart"
}

resource "sbercloud_networking_secgroup_rule" "secgroup_quickstart_rule_22" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = sbercloud_networking_secgroup.secgroup_quickstart.id
}

resource "sbercloud_networking_secgroup_rule" "secgroup_quickstart_rule_80" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 80
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = sbercloud_networking_secgroup.secgroup_quickstart.id
}

resource "sbercloud_networking_secgroup_rule" "secgroup_quickstart_rule_443" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 443
  port_range_max    = 443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = sbercloud_networking_secgroup.secgroup_quickstart.id
}

resource "sbercloud_networking_secgroup_rule" "secgroup_quickstart_rule_5432" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 5432
  port_range_max    = 5432
  remote_ip_prefix  = "192.168.0.0/16"
  security_group_id = sbercloud_networking_secgroup.secgroup_quickstart.id
}
```

### Compute / VM Image

```terraform
data "sbercloud_images_image" "ubuntu_image" {
  name        = "Ubuntu 20.04 server 64bit"
  most_recent = true
}
```

### Compute / SSH Public Key

```terraform
resource "sbercloud_compute_keypair" "quickstart" {
  name       = "terraform-keypair"
  public_key = "JUST_REPLACE_IT_YOURS_PUBLIC_KEY"
}
```

### Compute / VM

```terraform
resource "sbercloud_compute_instance" "ecs_quickstart" {

  name              = "ecs-quickstart"
  image_id          = data.sbercloud_images_image.ubuntu_image.id
  flavor_id         = "c6.6xlarge.2"
  security_groups   = ["default", sbercloud_networking_secgroup.secgroup_quickstart.name]
  availability_zone = "ru-moscow-1b"
  key_pair          = sbercloud_compute_keypair.quickstart.name
  system_disk_type = "SAS"
  system_disk_size = 40

  network {
    uuid = sbercloud_vpc_subnet.subnet_quickstart.id
  }

}
```

### Network / Bind public IP to VM

```terraform
resource "sbercloud_compute_eip_associate" "associated_eip_quickstart" {
  public_ip   = sbercloud_vpc_eip.eip_quickstart.address
  instance_id = sbercloud_compute_instance.ecs_quickstart.id
  fixed_ip    = sbercloud_compute_instance.ecs_quickstart.network.0.fixed_ip_v4
}
```

### Managed Database / Postgres 12

```terraform
resource "sbercloud_rds_instance" "rds_quickstart" {

  name              = "rds_quickstart"
  flavor            = "rds.pg.c6.xlarge.4"
  vpc_id            = sbercloud_vpc.vpc_quickstart.id
  subnet_id         = sbercloud_vpc_subnet.subnet_quickstart.id
  security_group_id = sbercloud_networking_secgroup.secgroup_quickstart.id
  availability_zone = ["ru-moscow-1b"]

  db {
    type     = "PostgreSQL"
    version  = "12"
    password = ""
  }

  volume {
    type = "HIGH"
    size = 10
  }

  backup_strategy {
    start_time = "04:00-05:00"
    keep_days  = 1
  }

}
```
