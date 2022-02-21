---
subcategory: "Cloud Backup and Recovery (CBR)"
---

# sbercloud_cbr_vault

Manages a CBR Vault resource within Sbercloud.

## Example Usage

### Create a server type vault

```hcl
variable "vault_name" {}
variable "ecs_instance_id" {}
variable "evs_volume_id" {}

data "sbercloud_compute_instance" "test" {
  ...
}

resource "sbercloud_cbr_vault" "test" {
  name             = var.vault_name
  type             = "server"
  protection_type  = "backup"
  consistent_level = "crash_consistent"
  size             = 100

  resources {
    server_id = var.ecs_instance_id
  }
}
```

### Create a disk type vault

```hcl
variable "vault_name" {}
variable "evs_volume_id" {}

resource "sbercloud_cbr_vault" "test" {
  name             = var.vault_name
  type             = "disk"
  protection_type  = "backup"
  consistent_level = "crash_consistent"
  size             = 50
  auto_expand      = true

  resources {
    includes = [
      var.evs_volume_id
    ]
  }
}
```

### Create an SFS turbo type vault

```hcl
variable "vault_name" {}
variable "sfs_turbo_id" {}

resource "sbercloud_cbr_vault" "test" {
  name             = var.vault_name
  consistent_level = "crash_consistent"
  type             = "turbo"
  protection_type  = "backup"
  size             = 1000

  resources {
    includes = [
      var.sfs_turbo_id
    ]
  }
}
```

## Argument reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CBR vault. If omitted, the
  provider-level region will be used. Changing this will create a new vault.

* `name` - (Required, String) Specifies a unique name of the CBR vault. This parameter can contain a maximum of 64
  characters, which may consist of letters, digits, underscores(_) and hyphens (-).

* `type` - (Required, String, ForceNew) Specifies the object type of the CBR vault.
  Changing this will create a new vault. Vaild values are as follows:
  + **server** (Cloud Servers)
  + **disk** (EVS Disks)
  + **turbo** (SFS Turbo file systems)

* `consistent_level` - (Required, String, ForceNew) Specifies the backup specifications.
  The value is crash_consistent by default (crash consistent backup).

* `protection_type` - (Required, String, ForceNew) Specifies the protection type of the CBR vault.
  The valid value is **backup**.
  Changing this will create a new vault.

* `size` - (Required, Int) Specifies the vault sapacity, in GB. The valid value range is `1` to `10,485,760`.

* `auto_expand` - (Optional, Bool) Specifies to enable auto capacity expansion for the backup protection type vault.
  Default to **false**.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies a unique ID in UUID format of enterprise project.
  Changing this will create a new vault.

* `policy_id` - (Optional, String) Specifies a policy to associate with the CBR vault.

* `resources` - (Optional, List) Specifies an array of one or more resources to attach to the CBR vault.
  The [object](#cbr_vault_resources) structure is documented below.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the CBR vault.

<a name="cbr_vault_resources"></a>
The `resources` block supports:

* `server_id` - (Optional, String) Specifies the ID of the ECS instance to be backed up.

* `includes` - (Optional, List) Specifies the array of disk or SFS file system IDs which will be included in the backup.
  Only **disk** and **turbo** vault support this parameter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.

* `allocated` - The allocated capacity of the vault, in GB.

* `used` - The used capacity, in GB.

* `spec_code` - The specification code.

* `status` - The vault status.

* `storage` - The name of the bucket for the vault.

## Import

Vaults can be imported by their `id`. For example,

```
$ terraform import sbercloud_cbr_vault.test 01c33779-7c83-4182-8b6b-24a671fcedf8
```
