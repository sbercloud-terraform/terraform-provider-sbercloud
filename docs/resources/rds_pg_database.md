---
subcategory: "Relational Database Service (RDS)"
layout: "sbercloud"
page_title: "SberCloud: sbercloud_rds_pg_database"
description: ""
---

# sbercloud_rds_pg_database

Manages RDS PostgreSQL database resource within SberCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "sbercloud_rds_pg_database" "test" {
  instance_id = var.instance_id
  name        = "test"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the RDS PostgreSQL instance.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the database name. The value contains 1 to 63 characters, including
  letters, digits, and underscores (_). It cannot start with pg or a digit, and must be different from RDS for
  PostgreSQL template library names. RDS for PostgreSQL template libraries include **postgres**, **template0**, and
  **template1**.

  Changing this parameter will create a new resource.

* `owner` - (Optional, String) Specifies the database user. The value must be an existing username and must be different
  from system usernames. Defaults to **root**.

* `template` - (Optional, String, ForceNew) Specifies the name of the database template. Value options: **template0**,
  **template1**. Defaults to **template1**.

  Changing this parameter will create a new resource.

* `character_set` - (Optional, String, ForceNew) Specifies the database character set.
  For details, see [documentation](https://www.postgresql.org/docs/16/infoschema-character-sets.html).
  Defaults to **UTF8**.

  Changing this parameter will create a new resource.

* `lc_collate` - (Optional, String, ForceNew) Specifies the database collocation.
  Defaults to **en_US.UTF-8**.

  -> **NOTE:** For different collation rules, the execution result of a statement may be different.
  <br/> For example, the execution result of select 'a'>'A'; is false when this parameter is set to
  **en_US.utf8** and is true when this parameter is set to 'C'. If a database is migrated from "O" to
  PostgreSQL, this parameter needs to be set to 'C' to meet your expectations. You can query the supported
  collation rules from the pg_collation table.

  Changing this parameter will create a new resource.

* `lc_ctype` - (Optional, String, ForceNew) Specifies the database classification.
  Defaults to: **en_US.UTF-8**.

  Changing this parameter will create a new resource.

* `is_revoke_public_privilege` - (Optional, Bool, ForceNew) Specifies whether to revoke the PUBLIC CREATE permission of
  the public schema.
  + **true**: indicates that the permission will be revoked.
  + **false**: indicates that the permission will not be revoked.

  Defaults to **false**.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the database description. The value contains 0 to 512 characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of database which is formatted `<instance_id>/<name>`.

* `size` - Indicates the database size, in bytes.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The RDS postgresql database can be imported using the `instance_id` and `name` separated by a slash, e.g.

```bash
$ terraform import sbercloud_rds_pg_database.test <instance_id>/<name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `template`, `is_revoke_public_privilege`
`lc_ctype`. It is generally recommended running `terraform plan` after importing the RDS PostgreSQL database. You can
then decide if changes should be applied to the RDS PostgreSQL database, or the resource definition should be updated
to align with the RDS PostgreSQL database. Also you can ignore changes as below.

```hcl
resource "sbercloud_rds_pg_database" "account_1" {
    ...

  lifecycle {
    ignore_changes = [
      template, is_revoke_public_privilege, lc_ctype,
    ]
  }
}
```
