---
subcategory: "Enterprise Router (ER)"
layout: "sbercloud"
page_title: "SberCloud: sbercloud_er_propagation"
description: ""
---

# sbercloud_er_propagation

Manages a propagation resource under the route table for ER service within SberCloud.

Before using enterprise router, define custom endpoint as shown below:
```terraform
provider "sbercloud" {
  auth_url = "https://iam.ru-moscow-1.hc.sbercloud.ru/v3"
  region   = "ru-moscow-1"
  access_key = var.access_key
  secret_key = var.secret_key

  endpoints = {
    er  = "https://er.ru-moscow-1.hc.cloud.ru"
  }
}
```

## Example Usage

```hcl
variable "instance_id" {}
variable "route_table_id" {}
variable "attachment_id" {}

resource "sbercloud_er_propagation" "test" {
  instance_id    = var.instance_id
  route_table_id = var.route_table_id
  attachment_id  = var.attachment_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the ER instance and route table are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the ER instance to which the route table and the
  attachment belongs.  
  Changing this parameter will create a new resource.

* `route_table_id` - (Required, String, ForceNew) Specifies the ID of the route table to which the propagation
  belongs.  
  Changing this parameter will create a new resource.

* `attachment_id` - (Required, String, ForceNew) Specifies the ID of the attachment corresponding to the propagation.  
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `attachment_type` - The type of the attachment corresponding to the propagation.

* `status` - The current status of the propagation.

* `created_at` - The creation time.

* `updated_at` - The latest update time.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 2 minutes.

## Import

Propagations can be imported using their `id` and the related `instance_id` and `route_table_id`, separated by
slashes (/), e.g.

```bash
$ terraform import sbercloud_er_propagation.test <instance_id>/<route_table_id>/<id>
```
