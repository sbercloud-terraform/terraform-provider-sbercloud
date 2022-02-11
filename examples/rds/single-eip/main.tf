# Get the VPC where RDS instance will be created
data "sbercloud_vpc" "vpc_01" {
  name = "put_here_name_of_your_existing_VPC"
}

# Get the subnet where RDS instance will be created
data "sbercloud_vpc_subnet" "subnet_01" {
  name = "put_here_name_of_your_existing_subnet"
}

# Get the security group for RDS instance
data "sbercloud_networking_secgroup" "sg_01" {
  name = "put_here_name_of_your_existing_security_group"
}

# Get RDS flavors
data "sbercloud_rds_flavors" "rds_flavors" {
  db_type       = "PostgreSQL"
  db_version    = "13"
  instance_mode = "single"
}

locals {
  rds_flavor = compact([
    for item in data.sbercloud_rds_flavors.rds_flavors.flavors :
    item["vcpus"] == 2 && item["memory"] == 4 ? item["name"] : ""
  ])[0]
}
  
# Create RDS instance
resource "sbercloud_rds_instance" "rds_01" {
  name                  = "terraform-pg-single"
  flavor                = local.rds_flavor
  vpc_id                = data.sbercloud_vpc.vpc_01.id
  subnet_id             = data.sbercloud_vpc_subnet.subnet_01.id
  security_group_id     = data.sbercloud_networking_secgroup.sg_01.id
  availability_zone     = ["ru-moscow-1b"]

  db {
    type     = "PostgreSQL"
    version  = "13"
    password = "put_here_password_of_database_root_user"
  }

  volume {
    type = "ULTRAHIGH"
    size = 100
  }
}

# Create EIP for RDS
resource "sbercloud_vpc_eip" "eip_rds" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    share_type  = "PER"
    name        = "eip-for-rds"
    size        = 4
    charge_mode = "bandwidth"
  }
}

# Get the port of the RDS instance by private_ip
data "sbercloud_networking_port" "rds_port" {
  network_id = data.sbercloud_vpc_subnet.subnet_01.id
  fixed_ip   = sbercloud_rds_instance.rds_01.private_ips[0]
}

# Attach EIP to the RDS instance
resource "sbercloud_networking_eip_associate" "associated" {
  public_ip = sbercloud_vpc_eip.eip_rds.address
  port_id   = data.sbercloud_networking_port.rds_port.id
}
