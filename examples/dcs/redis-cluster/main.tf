# Get the VPC where RDS instance will be created
data "sbercloud_vpc" "vpc_01" {
  name = "put_here_name_of_your_existing_VPN"
}

# Get the subnet where RDS instance will be created
data "sbercloud_vpc_subnet" "subnet_01" {
  name = "put_here_name_of_your_existing_subnet"
}

# Create Redis cluster
resource "sbercloud_dcs_instance" "redis_01" {
  name               = "redis-tf-cluster"
  engine             = "Redis"
  engine_version     = "5.0"
  capacity           = 8
  password           = "put_here_password_for_Redis"
  vpc_id             = data.sbercloud_vpc.vpc_01.id
  subnet_id          = data.sbercloud_vpc_subnet.subnet_01.id
  availability_zones = ["ru-moscow-1a", "ru-moscow-1b"]
  flavor             = "redis.cluster.xu1.large.r2.8"

  backup_policy {
    save_days   = 5
    backup_type = "auto"
    begin_at    = "02:00-03:00"
    period_type = "weekly"
    backup_at   = [2, 4, 6]
  }

  tags = {
    "environment" = "test"
  }
}
