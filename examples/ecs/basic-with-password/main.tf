# Declare all required input variables
variable "root_password" {
  description = "Root password for ECS"
  sensitive   = true
}

# Get the latest Ubuntu image
data "sbercloud_images_image" "ubuntu_image" {
  name        = "Ubuntu 20.04 server 64bit"
  most_recent = true
}

# Get the subnet where ECS will be created
data "sbercloud_vpc_subnet" "subnet_01" {
  name = "place_the_name_of_your_existing_subnet_here"
}

# Create ECS
resource "sbercloud_compute_instance" "ecs_01" {
  name              = "terraform-ecs"
  image_id          = data.sbercloud_images_image.ubuntu_image.id
  flavor_id         = "s6.large.2"
  security_groups   = ["sg-ssh"]
  availability_zone = "ru-moscow-1a"
  admin_pass        = var.root_password

  system_disk_type = "SAS"
  system_disk_size = 16

  network {
    uuid = data.sbercloud_vpc_subnet.subnet_01.id
  }
}
