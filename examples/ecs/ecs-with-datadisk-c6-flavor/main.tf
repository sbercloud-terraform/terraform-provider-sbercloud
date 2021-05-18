# Get the latest Ubuntu image
data "sbercloud_images_image" "ubuntu_image" {
  name        = "Ubuntu 20.04 server 64bit"
  most_recent = true
}

# Get the flavor name
data "sbercloud_compute_flavors" "flavor_name" {
  availability_zone = "ru-moscow-1b"
  performance_type  = "computingv3"
  cpu_core_count    = 2
  memory_size       = 4
}

# Get the subnet where ECS will be created
data "sbercloud_vpc_subnet" "subnet_01" {
  name = "public-01"
}

# Create ECS
resource "sbercloud_compute_instance" "ecs_01" {
  name              = "terraform-ecs"
  image_id          = data.sbercloud_images_image.ubuntu_image.id
  flavor_id         = data.sbercloud_compute_flavors.flavor_name.ids[0]
  security_groups   = ["default", "sg-ssh"]
  availability_zone = "ru-moscow-1b"
  key_pair          = "KeyPair-01-SBC"

  system_disk_type = "SAS"
  system_disk_size = 16

  data_disks {
    type = "SAS"
    size = "32"
  }

  network {
    uuid = data.sbercloud_vpc_subnet.subnet_01.id
  }
}
