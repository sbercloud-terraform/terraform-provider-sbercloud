# Get the VPC where RDS instance will be created
data "sbercloud_vpc" "vpc_01" {
  name = "put_here_name_of_your_existing_VPC"
}

# Get the subnet where RDS instance will be created
data "sbercloud_vpc_subnet" "subnet_01" {
  name = "put_here_name_of_your_existing_subnet"
}

# Get the security group for RDS instance
data "sbercloud_networking_secgroup" "sg_01" {
  name = "put_here_name_of_your_existing_security_group"
}

# Get the list of Availability Zones (AZ)
data "sbercloud_availability_zones" "list_of_az" {}

# Get the RabbitMQ product details
data "sbercloud_dms_product" "rabbitmq" {
  engine            = "rabbitmq"
  instance_type     = "cluster"
  version           = "3.7.17"
  storage           = 1000
  storage_spec_code = "dms.physical.storage.ultra"
}

# Create RabbitMQ instance
resource "sbercloud_dms_instance" "rabbitmq_01" {
  name = "rabbitmq-tf"

  vpc_id            = data.sbercloud_vpc.vpc_01.id
  subnet_id         = data.sbercloud_vpc_subnet.subnet_01.id
  security_group_id = data.sbercloud_networking_secgroup.sg_01.id

  available_zones = data.sbercloud_availability_zones.list_of_az.names

  access_user       = "admin"
  password          = "put_here_password_of_rabbitmq_user"

  engine            = data.sbercloud_dms_product.rabbitmq.engine
  product_id        = data.sbercloud_dms_product.rabbitmq.id
  engine_version    = data.sbercloud_dms_product.rabbitmq.version
  storage_space     = data.sbercloud_dms_product.rabbitmq.storage
  storage_spec_code = data.sbercloud_dms_product.rabbitmq.storage_spec_code
}
