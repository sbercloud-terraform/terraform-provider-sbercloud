---
subcategory: "Scalable File Service (SFS)"
---

# sbercloud\_sfs\_file\_system

Provides information about a Shared File System (SFS).

## Example Usage

```hcl
variable "share_name" {}
variable "share_id" {}

data "sbercloud_sfs_file_system" "shares" {
  name = var.share_name
  id   = var.share_id
}
```

## Argument Reference
The following arguments are supported:

* `name` - (Optional, String) The name of the shared file system.

* `id` - (Optional, String) The UUID of the shared file system.

* `status` - (Optional, String) The status of the shared file system.

* `region` - (Optional, String) Specifies the region in which to obtain the shared file system.
  If omitted, the provider-level region will be used.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `availability_zone` - The availability zone name.
* 
* `description` - The description of the shared file system.

* `state` - The state of the shared file system.

* `size` - The size (GB) of the shared file system.

* `status` - The status of the shared file system.

* `is_public` - The level of visibility for the shared file system.

* `share_proto` - The protocol for sharing file systems.

* `metadata` - Metadata key and value pairs as a dictionary of strings.

* `export_location` - The path for accessing the shared file system.

* `access_level` - The level of the access rule.

* `access_rules_status` - The status of the share access rule.

* `access_type` - The type of the share access rule.

* `access_to` - The access that the back end grants or denies.

* `share_access_id` - The UUID of the share access rule.

* `mount_id` - The UUID of the mount location of the shared file system.

* `share_instance_id` - The access that the back end grants or denies.

* `preferred` - Identifies which mount locations are most efficient and are used preferentially when multiple mount locations exist.
