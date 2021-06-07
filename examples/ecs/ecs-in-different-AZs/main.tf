# Get the latest Ubuntu image
data "sbercloud_images_image" "ubuntu_image" {
  name        = "Ubuntu 20.04 server 64bit"
  most_recent = true
}

# Get the list of availability zones
data "sbercloud_availability_zones" "list_of_az" {}

# Define local variables
locals {
  number_of_az  = length(data.sbercloud_availability_zones.list_of_az.names)
}

# Get subnet where ECS will be attached to
data "sbercloud_vpc_subnet" "subnet_01" {
  name = "place_the_name_of_your_existing_subnet_here"
}

# Create ECS
resource "sbercloud_compute_instance" "ecs_01" {
  count = 4

  name              = "terraform-ecs-${count.index}"
  image_id          = data.sbercloud_images_image.ubuntu_image.id
  flavor_id         = "s6.large.2"
  security_groups   = ["default"]
  availability_zone = data.sbercloud_availability_zones.list_of_az.names[count.index % local.number_of_az]
  key_pair          = "place_the_name_of_your_existing_key_pair_here"

  system_disk_type = "SAS"
  system_disk_size = 16

  network {
    uuid = data.sbercloud_vpc_subnet.subnet_01.id
  }
}
