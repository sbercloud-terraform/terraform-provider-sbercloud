# Create VPC
resource "sbercloud_vpc" "vpc_01" {
  name = "vpc-terraform"
  cidr = "10.33.0.0/16"
}

# Create two subnets
resource "sbercloud_vpc_subnet" "subnet_01" {
  name       = "subnet-one"
  cidr       = "10.33.10.0/24"
  gateway_ip = "10.33.10.1"

  primary_dns   = "100.125.13.59"
  secondary_dns = "8.8.8.8"

  vpc_id = sbercloud_vpc.vpc_01.id
}

resource "sbercloud_vpc_subnet" "subnet_02" {
  name       = "subnet-two"
  cidr       = "10.33.20.0/24"
  gateway_ip = "10.33.20.1"

  primary_dns   = "100.125.13.59"
  secondary_dns = "8.8.8.8"

  vpc_id = sbercloud_vpc.vpc_01.id
}
