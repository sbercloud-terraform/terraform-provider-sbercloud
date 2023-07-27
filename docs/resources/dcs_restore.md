---
subcategory: "Distributed Cache Service (DCS)"
---

# sbercloud_dcs_restore

Manages a DCS instance within SberCloud.


## Example Usage

### Create a single mode Redis instance

```hcl
variable project_id {}
variable instance_id {}
variable backup_id {}


resource "sbercloud_dcs_restore" "test" {
  project_id  = var.project_id
  instance_id = var.instance_id
  backup_id   = var.backup_id
  remark      = "restore instance"
}
```



## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String, ForceNew) The enterprise project id of the dcs instance. Changing this creates a new instance.

* `instance_id` - (Required, String, ForceNew) A dcs_instance ID in UUID format.

* `backup_id` - (Required, String, ForceNew) ID of the backup record.

* `remark` - (Optional, String, ForceNew) Description of DCS instance restoration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.

* `restore_records` - Array of the restoration records.

The `restore_records` block supports:

* `status` - Restoration status:
  * `waiting` - DCS instance restoration is waiting to begin.
  * `restoring` - DCS instance restoration is in progress.
  * `succeed` - DCS instance restoration succeeded.
  * `failed` - DCS instance restoration failed.

* `progress` - Restoration progress.
* `restore_id` - ID of the restoration record.
* `backup_id` - ID of the backup record.
* `restore_remark` - Description of DCS instance restoration.
* `backup_remark` - Description of DCS instance backup.
* `created_at` - Time at which the restoration task is created.
* `updated_at` - Time at which DCS instance restoration completed.
* `restore_name` - Name of the restoration record.
* `backup_name` - Name of the backup record.
* `error_code` - Error code returned if DCS instance restoration fails.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 120 minutes.
* `update` - Default is 120 minutes.
* `delete` - Default is 15 minutes.

## Import

DCS instance can be imported using the `id`, e.g.

```bash
terraform import sbercloud_dcs_instance.instance_1 80e373f9-872e-4046-aae9-ccd9ddc55511
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `password`, `auto_renew`, `period`, `period_unit`, `rename_commands`,
`internal_version`, `save_days`, `backup_type`, `begin_at`, `period_type`, `backup_at`.
It is generally recommended running `terraform plan` after importing an instance.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```
resource "sbercloud_dcs_instance" "instance_1" {
    ...

  lifecycle {
    ignore_changes = [
      password, rename_commands,
    ]
  }
}
```
