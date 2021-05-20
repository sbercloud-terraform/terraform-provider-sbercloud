# Get the latest Ubuntu image
data "sbercloud_images_image" "ubuntu_image" {
  name        = "Ubuntu 20.04 server 64bit"
  most_recent = true
}

# Get the list of availability zones
data "sbercloud_availability_zones" "list_of_az" {}

# Get the flavor name
data "sbercloud_compute_flavors" "flavor_name" {
  availability_zone = data.sbercloud_availability_zones.list_of_az.names[1]
  performance_type  = "computingv3"
  cpu_core_count    = 4
  memory_size       = 8
}

# Get first subnet where ECS will be attached to
data "sbercloud_vpc_subnet" "subnet_01" {
  name = "place_the_name_of_your_existing_subnet_number_one"
}

# Get second subnet where ECS will be attached to
data "sbercloud_vpc_subnet" "subnet_02" {
  name = "place_the_name_of_your_existing_subnet_number_two"
}

# Get third subnet where ECS will be attached to
data "sbercloud_vpc_subnet" "subnet_03" {
  name = "place_the_name_of_your_existing_subnet_number_three"
}

# Create ECS
resource "sbercloud_compute_instance" "ecs_01" {
  name              = "terraform-ecs"
  image_id          = data.sbercloud_images_image.ubuntu_image.id
  flavor_id         = data.sbercloud_compute_flavors.flavor_name.ids[0]
  security_groups   = ["default"]
  availability_zone = data.sbercloud_availability_zones.list_of_az.names[1]
  key_pair          = "place_the_name_of_your_existing_key_pair_here"

  system_disk_type = "SAS"
  system_disk_size = 16

  network {
    uuid = data.sbercloud_vpc_subnet.subnet_01.id
  }
 
  network {
    uuid = data.sbercloud_vpc_subnet.subnet_02.id
  }
}

# Attach the ECS to third subnet
resource "sbercloud_compute_interface_attach" "attached_01" {
  instance_id = sbercloud_compute_instance.ecs_01.id
  network_id  = data.sbercloud_vpc_subnet.subnet_03.id
}
