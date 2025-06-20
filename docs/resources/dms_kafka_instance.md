---
subcategory: "Distributed Message Service (DMS)"
---

# sbercloud_dms_kafka_instance

Manage DMS Kafka instance resources within SberCloud.

## Example Usage

### Create a Kafka instance using flavor ID

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}

variable "availability_zones" {
  default = ["your_availability_zones_a", "your_availability_zones_b", "your_availability_zones_c"]
}
variable "flavor_id" {
  default = "your_flavor_id, such: c6.2u4g.cluster"
}
variable "storage_spec_code" {
  default = "your_storage_spec_code, such: dms.physical.storage.ultra.v2"
}

data "sbercloud_dms_kafka_flavors" "test" {
  type               = "cluster"
  flavor_id          = var.flavor_id
  availability_zones = var.availability_zones
  storage_spec_code  = var.storage_spec_code
}

resource "sbercloud_dms_kafka_instance" "test" {
  name              = "kafka_test"
  vpc_id            = var.vpc_id
  network_id        = var.subnet_id
  security_group_id = var.security_group_id

  flavor_id          = var.flavor_id
  storage_spec_code  = var.storage_spec_code
  availability_zones = var.availability_zones
  engine_version     = "2.7"
  storage_space      = 600
  broker_num         = 3

  access_user = "user"
  password    = "Kafka_%^&_Test"

  manager_user     = "kafka_manager"
  manager_password = "Kafka_Test^&*("

  depends_on = ["data.sbercloud_dms_kafka_flavors.test"]
}
```

-> **Why depend on "data.sbercloud_dms_kafka_flavors.test", it is not used?**
The specified `flavor_id` and `storage_spec_code` are not valid in all regions.
Before creating kafka, verify their validity through datasource to avoid creation errors.

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the DMS Kafka instances. If omitted, the
  provider-level region will be used. Changing this creates a new instance resource.

* `name` - (Required, String) Specifies the name of the DMS Kafka instance. An instance name starts with a letter,
  consists of 4 to 64 characters, and supports only letters, digits, hyphens (-) and underscores (_).

* `flavor_id` - (Optional, String) Specifies the Kafka `flavor ID`. This parameter and `product_id` are alternative.

  -> It is recommended to use `flavor_id` if the region supports it.\
  One of:
  * c6.2u4g.cluster
  * c6.4u8g.cluster
  * c6.8u16g.cluster
  * c6.12u24g.cluster
  * c6.16u32g.cluster \
  Or use data source sbercloud_dms_kafka_flavors

* `product_id` - (Optional, String) Specifies a product ID, which includes bandwidth, partition, broker and default
  storage capacity.

  -> **NOTE:** Change this to change the bandwidth, partition and broker of the Kafka instances. Please note that the
  broker changes may cause storage capacity changes. So, if you specify the value of `storage_space`, you need to
  manually modify the value of `storage_space` after changing the `product_id`.

* `engine_version` - (Required, String, ForceNew) Specifies the version of the Kafka engine,
  such as 1.1.0, 2.3.0, 2.7 or other supported versions. Changing this creates a new instance resource.

* `storage_spec_code` - (Required, String, ForceNew) Specifies the storage I/O specification. Value range:
  + When bandwidth is 100MB: dms.physical.storage.high or dms.physical.storage.ultra
  + When bandwidth is 300MB: dms.physical.storage.high or dms.physical.storage.ultra
  + When bandwidth is 600MB: dms.physical.storage.ultra
  + When bandwidth is 1200MB: dms.physical.storage.ultra

  If the instance is created with `product_id`, the valid values are as follows:
  + **dms.physical.storage.high**: Type of the disk that uses high I/O.
    The corresponding bandwidths are **100MB** and **300MB**.
  + **dms.physical.storage.ultra**: Type of the disk that uses ultra-high I/O.
    The corresponding bandwidths are **100MB**, **300MB**, **600MB** and **1,200MB**.

  Changing this creates a new instance resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of a VPC. Changing this creates a new instance resource.

* `network_id` - (Required, String, ForceNew) Specifies the ID of a subnet. Changing this creates a new instance
  resource.

* `security_group_id` - (Required, String) Specifies the ID of a security group.

* `availability_zones` - (Required, List, ForceNew) The names of the AZ where the Kafka instances reside.
  The parameter value can not be left blank or an empty array. Changing this creates a new instance resource.

  -> **NOTE:** Deploy one availability zone or at least three availability zones. Do not select two availability zones.
  Deploy to more availability zones, the better the reliability and SLA coverage.
  
  ~> The parameter behavior of `availability_zones` has been changed from `list` to `set`.

* `manager_user` - (Optional, String, ForceNew) Specifies the username for logging in to the Kafka Manager. The username
  consists of 4 to 64 characters and can contain letters, digits, hyphens (-), and underscores (_). Changing this
  creates a new instance resource.

* `manager_password` - (Optional, String, ForceNew) Specifies the password for logging in to the Kafka Manager. The
  password must meet the following complexity requirements: Must be 8 to 32 characters long. Must contain at least 2 of
  the following character types: lowercase letters, uppercase letters, digits, and special characters (`~!@#$%^&*()-_
  =+\\|[{}]:'",<.>/?). Changing this creates a new instance resource.

  -> **NOTE:** `manager_user` and `manager_password` are deprecated and will be deleted in future releases

* `storage_space` - (Optional, Int) Specifies the message storage capacity, the unit is GB.
  + When bandwidth is 100MB: 600–90000 GB
  + When bandwidth is 300MB: 1200–90000 GB
  + When bandwidth is 600MB: 2400–90000 GB
  + When bandwidth is 1200MB: 4800–90000 GB

  Changing this creates a new instance resource.

  It is required when creating an instance with `flavor_id`.

* `broker_num` - (Optional, Int) Specifies the broker numbers.
  It is required when creating an instance with `flavor_id`.

* `access_user` - (Optional, String, ForceNew) Specifies the username of SASL_SSL user. A username consists of 4
  to 64 characters and supports only letters, digits, and hyphens (-). Changing this creates a new instance resource. 

* `password` - (Optional, String, ForceNew) Specifies the password of SASL_SSL user. A password must meet the
  following complexity requirements: Must be 8 to 32 characters long. Must contain at least 2 of the following character
  types: lowercase letters, uppercase letters, digits, and special characters (`~!@#$%^&*()-_=+\\|[{}]:'",<.>/?).
  Changing this creates a new instance resource.

  -> **NOTE:** `access_user` and `password` parameters are mandatory when ssl_enable is set to true. These parameters are invalid when ssl_enable is set to false. \
  -> **NOTE:** If `access_user` and `password` are specified, set `ssl_enable = true`, to enable SASL_SSL for a Kafka instance.

* `description` - (Optional, String) Specifies the description of the DMS Kafka instance. It is a character string
  containing not more than 1,024 characters.

* `maintain_begin` - (Optional, String) Specifies the time at which a maintenance time window starts. Format: HH:mm. The
  start time and end time of a maintenance time window must indicate the time segment of a supported maintenance time
  window. The start time must be set to 22:00, 02:00, 06:00, 10:00, 14:00, or 18:00. Parameters `maintain_begin`
  and `maintain_end` must be set in pairs. If parameter `maintain_begin` is left blank, parameter `maintain_end` is also
  blank. In this case, the system automatically allocates the default start time 02:00.

* `maintain_end` - (Optional, String) Specifies the time at which a maintenance time window ends. Format: HH:mm. The
  start time and end time of a maintenance time window must indicate the time segment of a supported maintenance time
  window. The end time is four hours later than the start time. For example, if the start time is 22:00, the end time is
  02:00. Parameters `maintain_begin`
  and `maintain_end` must be set in pairs. If parameter `maintain_end` is left blank, parameter
  `maintain_begin` is also blank. In this case, the system automatically allocates the default end time 06:00.

* `public_ip_ids` - (Optional, List, ForceNew) Specifies the IDs of the elastic IP address (EIP)
  bound to the DMS Kafka instance. Changing this creates a new instance resource.
  + If the instance is created with `flavor_id`, the total number of public IPs is equal to `broker_num`.
  + If the instance is created with `product_id`, the total number of public IPs must provide as follows:

  | Bandwidth | Total number of public IPs |
    | ---- | ---- |
  | 100MB | 3 |
  | 300MB | 3 |
  | 600MB | 4 |
  | 1,200MB | 8 |

* `retention_policy` - (Optional, String) Specifies the action to be taken when the memory usage reaches the disk
  capacity threshold. The valid values are as follows:
  + **time_base**: Automatically delete the earliest messages.
  + **produce_reject**: Stop producing new messages.

* `dumping` - (Optional, Bool, ForceNew) Specifies whether to enable message dumping.
  Changing this creates a new instance resource.

* `enable_auto_topic` - (Optional, Bool, ForceNew) Specifies whether to enable automatic topic creation. If automatic
  topic creation is enabled, a topic will be automatically created with 3 partitions and 3 replicas when a message is
  produced to or consumed from a topic that does not exist.
  The default value is false.
  Changing this creates a new instance resource.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the Kafka instance.

* `tags` - (Optional, Map) The key/value pairs to associate with the DMS Kafka instance.

* `cross_vpc_accesses` - (Optional, List) Specifies the cross-VPC access information.
  The [object](#dms_cross_vpc_accesses) structure is documented below.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the instance. Valid values are *prePaid*
  and *postPaid*, defaults to *postPaid*. Changing this creates a new resource.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the instance.
  Valid values are *month* and *year*. This parameter is mandatory if `charging_mode` is set to *prePaid*.
  Changing this creates a new resource.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the instance. If `period_unit` is set to *month*
  , the value ranges from 1 to 9. If `period_unit` is set to *year*, the value ranges from 1 to 3. This parameter is
  mandatory if `charging_mode` is set to *prePaid*. Changing this creates a new resource.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled. Valid values are "true" and "false".

<a name="dms_cross_vpc_accesses"></a>
The `cross_vpc_accesses` block supports:

* `advertised_ip` - (Optional, String) The advertised IP Address or domain name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.
* `engine` - Indicates the message engine.
* `partition_num` - Indicates the number of partitions in Kafka instance.
* `used_storage_space` - Indicates the used message storage space. Unit: GB
* `port` - Indicates the port number of the DMS Kafka instance.
* `status` - Indicates the status of the DMS Kafka instance.
* `ssl_enable` - Indicates whether the Kafka SASL_SSL is enabled.
* `enable_public_ip` - Indicates whether public access to the DMS Kafka instance is enabled.
* `resource_spec_code` - Indicates a resource specifications identifier.
* `type` - Indicates the DMS Kafka instance type.
* `user_id` - Indicates the ID of the user who created the DMS Kafka instance
* `user_name` - Indicates the name of the user who created the DMS Kafka instance
* `connect_address` - Indicates the IP address of the DMS Kafka instance.
* `management_connect_address` - Indicates the Kafka Manager connection address of a Kafka instance.
* `cross_vpc_accesses` - Indicates the Access information of cross-VPC. The structure is documented below.
* `charging_mode` - Indicates the charging mode of the instance.

The `cross_vpc_accesses` block supports:

* `listener_ip` - The listener IP address.
* `port` - The port number.
* `port_id` - The port ID associated with the address.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 50 minute.
* `update` - Default is 50 minute.
* `delete` - Default is 15 minute.

## Import

DMS Kafka instance can be imported using the instance id, e.g.

```
 $ terraform import sbercloud_dms_kafka_instance.instance_1 8d3c7938-dc47-4937-a30f-c80de381c5e3
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include:
`password`, `manager_password` and `public_ip_ids`. It is generally recommended running `terraform plan` after importing
a DMS Kafka instance. You can then decide if changes should be applied to the instance, or the resource definition
should be updated to align with the instance. Also you can ignore changes as below.

```
resource "sbercloud_dms_kafka_instance" "instance_1" {
    ...

  lifecycle {
    ignore_changes = [
      password, manager_password,
    ]
  }
}
```
