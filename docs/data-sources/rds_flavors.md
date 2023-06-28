---
subcategory: "Relational Database Service (RDS)"
---

# sbercloud\_rds\_flavors

Use this data source to get available SberCloud rds flavors.

## Example Usage

```hcl
data "sbercloud_rds_flavors" "flavor" {
  db_type       = "PostgreSQL"
  db_version    = "12"
  instance_mode = "ha"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the RDS flavors. If omitted, the provider-level region will be used.

* `db_type` - (Required, String) Specifies the DB engine. Value: MySQL, PostgreSQL, SQLServer.

* `db_version` - (Required, String) Specifies the database version. Available value:

type | version
---- | ---
MySQL| 5.6 <br>5.7 <br>8.0
PostgreSQL | 9.5 <br> 9.6 <br>10 <br>11 <br>12
SQLServer| 2012_SE <br>2014_SE <br>2016_SE <br>2012_EE <br>2014_EE <br>2016_EE <br>2017_EE

* `instance_mode` - (Required, String) The mode of instance. Value: *ha*(indicates primary/standby instance),
  *single*(indicates single instance) and *replica*(indicates read replicas).

* `vcpus` - (Optional, Int) Specifies the number of vCPUs in the RDS flavor.

* `memory` - (Optional, Int) Specifies the memory size(GB) in the RDS flavor.

* `group_type` - (Optional, String) Specifies the performance specification, the valid values are as follows:
  + **normal**: General enhanced.
  + **normal2**: General enhanced type II.
  + **armFlavors**: KunPeng general enhancement.
  + **dedicatedNormal**: (dedicatedNormalLocalssd): Dedicated for x86.
  + **armLocalssd**: KunPeng general type.
  + **normalLocalssd**: x86 general type.
  + **general**: General type.
  + **dedicated**:  
    For MySQL engine: Dedicated type.  
    For PostgreSQL and SQL Server engines: Dedicated type, only supported by cloud disk SSD.
  + **rapid**:  
    For MySQL engine: Dedicated (discontinued).  
    For PostgreSQL and SQL Server engines: Dedicated, only supported by ultra-fast SSDs.
  + **bigmem**: Large memory type.

* `availability_zone` - (Optional, String) Specifies the availability zone which the RDS flavor belongs to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - Indicates the flavors information. Structure is documented below.

The `flavors` block contains:

* `id` - The ID of the rds flavor.
* `name` - The name of the rds flavor.
* `vcpus` - The CPU size.
* `memory` - The memory size in GB.
* `group_type` - The performance specification.
* `instance_mode` - The mode of instance.
* `availability_zones` - The availability zones which the RDS flavor belongs to.
* `db_versions` - The Available versions of the database.
