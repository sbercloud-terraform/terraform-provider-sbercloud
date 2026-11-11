---
subcategory: "Elastic Cloud Server (ECS)"
layout: "sbercloud"
page_title: "Sbercloud: sbercloud_compute_password_delete"
description: |-
  Manages an ECS password records delete generated during initial installation of a Windows ECS resource within Sbercloud.
---

# sbercloud_compute_password_delete

Manages an ECS password records delete generated during initial installation of a Windows ECS resource within Sbercloud.

## Example Usage

```hcl
variable "server_id" {}

resource "sbercloud_compute_password_delete" "test" {
  server_id = var.server_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `server_id` - (Required, String, NonUpdatable) Specifies the ECS ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the instance ID.