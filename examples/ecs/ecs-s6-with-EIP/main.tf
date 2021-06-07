# Get the latest Ubuntu image
data "sbercloud_images_image" "ubuntu_image" {
  name        = "Ubuntu 20.04 server 64bit"
  most_recent = true
}

# Get the list of availability zones
data "sbercloud_availability_zones" "list_of_az" {}

# Get the flavor name
data "sbercloud_compute_flavors" "flavor_name" {
  availability_zone = data.sbercloud_availability_zones.list_of_az.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 8
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
  availability_zone = data.sbercloud_availability_zones.list_of_az.names[0]
  key_pair          = "place_the_name_of_your_existing_key_pair_here"

  system_disk_type = "SAS"
  system_disk_size = 16

  network {
    uuid = data.sbercloud_vpc_subnet.subnet_01.id
  }
}

# Create EIP
resource "sbercloud_vpc_eip" "eip_01" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "ecs-bandwidth"
    size        = 4
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

# Attach the EIP to the ECS
resource "sbercloud_compute_eip_associate" "associated_01" {
  public_ip   = sbercloud_vpc_eip.eip_01.address
  instance_id = sbercloud_compute_instance.ecs_01.id
}
