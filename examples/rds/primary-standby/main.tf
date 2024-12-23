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

# Get the list of availability zones
data "sbercloud_availability_zones" "list_of_az" {}

# Get RDS flavors
data "sbercloud_rds_flavors" "rds_flavors" {
  db_type       = "PostgreSQL"
  db_version    = "12"
  instance_mode = "ha"
}

locals {
  rds_flavor = compact([
    for item in data.sbercloud_rds_flavors.rds_flavors.flavors :
    item["vcpus"] == "2" && item["memory"] == 8 ? item["name"] : ""
  ])[0]
}

resource "sbercloud_rds_instance" "rds_01" {
  name                  = "terraform-pg-cluster"
  flavor                = local.rds_flavor
  vpc_id                = data.sbercloud_vpc.vpc_01.id
  subnet_id             = data.sbercloud_vpc_subnet.subnet_01.id
  security_group_id     = data.sbercloud_networking_secgroup.sg_01.id
  availability_zone     = [data.sbercloud_availability_zones.list_of_az.names[0], data.sbercloud_availability_zones.list_of_az.names[1]]
  ha_replication_mode   = "async"

  db {
    type     = "PostgreSQL"
    version  = "12"
    password = "put_here_password_of_database_root_user"
  }

  volume {
    type = "CLOUDSSD"
    size = 100
  }

  backup_strategy {
    start_time = "01:00-02:00"
    keep_days  = 3
  }

  tags = {
    "environment" = "stage"
  }
}
