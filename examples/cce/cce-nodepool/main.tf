# Get CCE cluster
data "sbercloud_cce_cluster" "cce_01" {
  name   = "put_here_the_name_of_your_existing_CCE_cluster"
  status = "Available"
}

# Create CCE Node pool
resource "sbercloud_cce_node_pool" "node_pool_01" {
  cluster_id               = data.sbercloud_cce_cluster.cce_01.id
  name                     = "terraform-pool"
  flavor_id                = "s6.xlarge.4"
  availability_zone        = "ru-moscow-1a"
  key_pair                 = "put_here_the_name_of_your_existing_key_pair"
  scall_enable             = true
  min_node_count           = 2
  initial_node_count       = 2
  max_node_count           = 10
  scale_down_cooldown_time = 100
  priority                 = 1
  type                     = "vm"
  os                       = "CentOS 7.6"

  labels = {
    created_by = "Terraform"
    creation_date = "December2021"
  }

  root_volume {
    size       = 50
    volumetype = "SAS"
  }

  data_volumes {
    size       = 100
    volumetype = "SAS"
  }
}
