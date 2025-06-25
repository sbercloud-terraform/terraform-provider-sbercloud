variable "nat_gatway_name" {
  default = "nat-gateway-basic"
}

variable "vpc_name" {
  default = "vpc-basic"
}

variable "vpc_cidr" {
  default = "192.168.0.0/16"
}

variable "subnet_name" {
  default = "subent-basic"
}

variable "subnet_cidr" {
  default = "192.168.10.0/24"
}

variable "subnet_gateway" {
  default = "192.168.10.1"
}

variable "primary_dns" {
  default = "100.125.13.59"
}
