# Get the VPC where NAT gateway will be created
data "sbercloud_vpc" "vpc_01" {
  name = "put_here_the_name_of_your_existing_VPC"
}

# Get the subnet where NAT gateway will be created
data "sbercloud_vpc_subnet" "subnet_01" {
  name = "put_here_the_name_of_your_existing_subnet"
}

# Create EIP for NAT Gateway
resource "sbercloud_vpc_eip" "nat_eip" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "nat_bandwidth"
    size        = 4
    share_type  = "PER"
    charge_mode = "bandwidth"
  }
}

# Create NAT Gateway
resource "sbercloud_nat_gateway" "nat_01" {
  name        = "nat-terraform"
  description = "Demo NAT Gateway"
  spec        = "1"
  vpc_id      = data.sbercloud_vpc.vpc_01.id
  subnet_id   = data.sbercloud_vpc_subnet.subnet_01.id
}

# Create SNAT rule for your subnet
resource "sbercloud_nat_snat_rule" "snat_subnet_01" {
  nat_gateway_id = sbercloud_nat_gateway.nat_01.id
  subnet_id      = data.sbercloud_vpc_subnet.subnet_01.id
  floating_ip_id = sbercloud_vpc_eip.nat_eip.id
}
