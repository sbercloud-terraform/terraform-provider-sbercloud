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

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a data source ID in UUID format.

* `flavors` -
  Indicates the flavors information. Structure is documented below.

The `flavors` block contains:

* `name` - The name of the rds flavor.
* `vcpus` - Indicates the CPU size.
* `memory` - Indicates the memory size in GB.
* `mode` - See 'instance_mode' above.
