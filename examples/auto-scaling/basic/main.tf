data "sbercloud_availability_zones" "myaz" {}

data "sbercloud_compute_flavors" "myflavor" {
  availability_zone = data.sbercloud_availability_zones.myaz.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "sbercloud_images_image" "myimage" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "sbercloud_vpc" "vpc_1" {
  name = var.vpc_name
}

data "sbercloud_vpc_subnet" "subnet_1" {
  name   = var.subnet_name
  vpc_id = data.sbercloud_vpc.vpc_1.id
}

data "sbercloud_networking_secgroup" "secgroup_1" {
  name = var.secgroup_name
}

resource "sbercloud_compute_keypair" "my_keypair" {
  name       = "my_keypair"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB"
}

resource "sbercloud_as_configuration" "my_as_config" {
  scaling_configuration_name = "my_as_config"

  instance_config {
    flavor   = data.sbercloud_compute_flavors.myflavor.ids[0]
    image    = data.sbercloud_images_image.myimage.id
    key_name = sbercloud_compute_keypair.my_keypair.name
    disk {
      size        = 40
      volume_type = "SSD"
      disk_type   = "SYS"
    }
  }
}

resource "sbercloud_as_group" "my_as_group" {
  scaling_group_name       = "my_as_group"
  scaling_configuration_id = sbercloud_as_configuration.my_as_config.id
  desire_instance_number   = 2
  min_instance_number      = 0
  max_instance_number      = 10
  vpc_id                   = data.sbercloud_vpc.vpc_1.id
  delete_publicip          = true
  delete_instances         = "yes"

  networks {
    id = data.sbercloud_vpc_subnet.subnet_1.id
  }
  security_groups {
    id = data.sbercloud_networking_secgroup.secgroup_1.id
  }
  tags = {
    owner = "AutoScaling"
  }
}
