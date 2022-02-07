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

# Get the Kafka product details
data "sbercloud_dms_product" "kafka" {
  engine            = "kafka"
  instance_type     = "cluster"
  version           = "2.3.0"
  storage           = 1200
  bandwidth         = "300MB"
  storage_spec_code = "dms.physical.storage.high"
}

# Create Kafka instance
resource "sbercloud_dms_instance" "kafka_01" {
  name = "kafka-tf"

  vpc_id            = data.sbercloud_vpc.vpc_01.id
  subnet_id         = data.sbercloud_vpc_subnet.subnet_01.id
  security_group_id = data.sbercloud_networking_secgroup.sg_01.id

  available_zones = data.sbercloud_availability_zones.list_of_az.names

  engine            = data.sbercloud_dms_product.kafka.engine
  specification     = data.sbercloud_dms_product.kafka.bandwidth
  product_id        = data.sbercloud_dms_product.kafka.id
  engine_version    = data.sbercloud_dms_product.kafka.version
  storage_space     = data.sbercloud_dms_product.kafka.storage
  storage_spec_code = data.sbercloud_dms_product.kafka.storage_spec_code
}

# Create topic
resource "sbercloud_dms_kafka_topic" "topic" {
  instance_id = sbercloud_dms_instance.kafka_01.id
  name        = "topic_01"
  partitions  = 16
}
