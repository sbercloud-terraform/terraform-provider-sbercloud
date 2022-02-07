# Get the VPC where RDS instance will be created
data "sbercloud_vpc" "vpc_01" {
  name = "put_here_name_of_your_existing_VPC"
}

# Get the subnet where RDS instance will be created
data "sbercloud_vpc_subnet" "subnet_01" {
  name = "put_here_name_of_your_existing_subnet"
}

resource "sbercloud_dcs_instance" "redis_01" {
  name               = "redis-tf-single"
  engine             = "Redis"
  engine_version     = "5.0"
  capacity           = 4
  vpc_id             = data.sbercloud_vpc.vpc_01.id
  subnet_id          = data.sbercloud_vpc_subnet.subnet_01.id
  availability_zones = ["ru-moscow-1b"]
  flavor             = "redis.single.xu1.large.4"
}
