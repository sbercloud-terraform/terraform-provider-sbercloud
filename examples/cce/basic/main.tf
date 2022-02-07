# Get the VPC where CCE cluster will be created
data "sbercloud_vpc" "vpc_01" {
  name = "put_here_the_name_of_your_existing_vpc"
}

# Get the subnet where CCE cluster will be created
data "sbercloud_vpc_subnet" "subnet_01" {
  name = "put_here_the_name_of_your_existing_subnet"
}

# Create CCE cluster
resource "sbercloud_cce_cluster" "cce_01" {
  name                   = "demo-cluster"
  flavor_id              = "cce.s2.small"
  container_network_type = "overlay_l2"
  multi_az               = true
  vpc_id                 = data.sbercloud_vpc.vpc_01.id
  subnet_id              = data.sbercloud_vpc_subnet.subnet_01.id
}

# Create CCE worker node(s)
resource "sbercloud_cce_node" "cce_01_node" {
  cluster_id        = sbercloud_cce_cluster.cce_01.id
  name              = "cce-worker"
  flavor_id         = "s6.large.2"
  availability_zone = "ru-moscow-1a"
  os                = "CentOS 7.6"
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
