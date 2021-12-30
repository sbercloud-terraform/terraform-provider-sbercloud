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

resource "sbercloud_css_cluster" "css_prod" {
  expect_node_num = 1
  name            = "css-terraform"
  engine_version  = "7.9.3"

  node_config {
    flavor = "ess.spec-4u32g"
    network_info {
      security_group_id = data.sbercloud_networking_secgroup.sg_01.id
      subnet_id         = data.sbercloud_vpc_subnet.subnet_01.id
      vpc_id            = data.sbercloud_vpc.vpc_01.id
    }
    volume {
      volume_type = "ULTRAHIGH"
      size        = 80
    }
    availability_zone = "ru-moscow-1a"
  }

  backup_strategy {
    bucket      = "p-test-02"
    backup_path = "css_backups/css-terraform"
    agency      = "css_obs_agency"
    start_time  = "01:00 GMT+03:00"
    keep_days   = 4
  }

  tags = {
    "environment" = "prod"
  }
}
