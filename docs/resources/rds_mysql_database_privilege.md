---
subcategory: "Relational Database Service (RDS)"
layout: "sbercloud"
page_title: "SberCloud: sbercloud_rds_mysql_database_privilege"
description: ""
---

# sbercloud_rds_mysql_database_privilege

Manages RDS Mysql database privilege resource within SberCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "db_name" {}
variable "user_name_1" {}
variable "user_name_2" {}

resource "sbercloud_rds_mysql_database_privilege" "test" {
  instance_id = var.instance_id
  db_name     = var.db_name

  users {
    name     = var.user_name_1
    readonly = true
  }

  users {
    name     = var.user_name_2
    readonly = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the RDS database privilege resource. If omitted,
  the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the RDS instance ID. Changing this will create a new resource.

* `db_name` - (Required, String, ForceNew) Specifies the database name. Changing this creates a new resource.

* `users` - (Required, List) Specifies the account that associated with the database. Structure is documented below.

The `users` block supports:

* `name` - (Required, String) Specifies the username of the database account.

* `readonly` - (Optional, Bool) Specifies the read-only permission. The value can be:
  + **true**: indicates the read-only permission.
  + **false**: indicates the read and write permission.

  The default value is **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of database privilege which is formatted `<instance_id>/<db_name>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

RDS database privilege can be imported using the `instance id` and `db_name`, e.g.

```bash
$ terraform import sbercloud_rds_mysql_database_privilege.test <instance_id>/<db_name>
```
