# Get the latest Ubuntu image
data "sbercloud_images_image" "ubuntu_image" {
  name        = "Ubuntu 20.04 server 64bit"
  most_recent = true
}

# Get the flavor name
data "sbercloud_compute_flavors" "flavor_name" {
  availability_zone = "ru-moscow-1a"
  performance_type  = "highmem"
  cpu_core_count    = 4
  memory_size       = 32
}

# Get the subnet where ECS will be created
data "sbercloud_vpc_subnet" "subnet_01" {
  name = "place_the_name_of_your_existing_subnet_here"
}

# Create ECS
resource "sbercloud_compute_instance" "ecs_01" {
  name              = "terraform-ecs"
  image_id          = data.sbercloud_images_image.ubuntu_image.id
  flavor_id         = data.sbercloud_compute_flavors.flavor_name.ids[0]
  security_groups   = ["default"]
  availability_zone = "ru-moscow-1a"
  key_pair          = "place_the_name_of_your_existing_key_pair_here"

  system_disk_type = "SAS"
  system_disk_size = 16

  network {
    uuid = data.sbercloud_vpc_subnet.subnet_01.id
  }
}

# Create additional EVS volume
resource "sbercloud_evs_volume" "volume_01" {
  name              = "additional-volume"
  description       = "Additional volume"
  volume_type       = "SSD"
  size              = 64
  availability_zone = "ru-moscow-1a"

  tags = {
    created_by = "terraform"
  }
}

# Attach the disk to the ECS
resource "sbercloud_compute_volume_attach" "attached_01" {
  instance_id = sbercloud_compute_instance.ecs_01.id
  volume_id   = sbercloud_evs_volume.volume_01.id
}
