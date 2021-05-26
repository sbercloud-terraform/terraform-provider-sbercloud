# Get the VPC where CCE cluster will be created
data "sbercloud_vpc" "vpc_01" {
  name = "put_here_the_name_of_your_existing_vpc"
}

# Get the subnet where CCE cluster will be created
data "sbercloud_vpc_subnet" "subnet_01" {
  name = "put_here_the_name_of_your_existing_subnet"
}

# Get the list of availability zones
data "sbercloud_availability_zones" "list_of_az" {}

# Set local variables: the number of worker nodes in CCE cluster and get the number of Availability Zones
locals {
  number_of_workers = 2
  number_of_az      = length(data.sbercloud_availability_zones.list_of_az.names)
}

# Create CCE cluster
resource "sbercloud_cce_cluster" "cce_01" {
  name                   = "demo-cluster"
  flavor_id              = "cce.s2.medium"
  container_network_type = "vpc-router"
  multi_az               = true
  kube_proxy_mode        = "ipvs"
  vpc_id                 = data.sbercloud_vpc.vpc_01.id
  subnet_id              = data.sbercloud_vpc_subnet.subnet_01.id
}

# Create CCE worker node(s)
resource "sbercloud_cce_node" "cce_01_node" {
  count             = local.number_of_workers
  cluster_id        = sbercloud_cce_cluster.cce_01.id
  name              = "cce-worker-${count.index}"
  flavor_id         = "s6.large.2"
  availability_zone = data.sbercloud_availability_zones.list_of_az.names[count.index % local.number_of_az]
  key_pair          = "put_here_the_name_of_your_existing_key_pair"

  root_volume {
    size       = 50
    volumetype = "SAS"
  }

  data_volumes {
    size       = 100
    volumetype = "SAS"
  }
}
