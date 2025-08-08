---
subcategory: "Cloud Firewall (CFW)"
layout: "sbercloud"
page_title: "SberCloud: sbercloud_cfw_ips_rule_details"
description: |-
  Use this data source to get the list of CFW IPS rule details.
---

# sbercloud_cfw_ips_rule_details

Use this data source to get the list of CFW IPS rule details.

## Example Usage

```hcl
variable "fw_instance_id" {}

data "sbercloud_cfw_ips_rule_details" "test" {
  fw_instance_id = var.fw_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The IPS information.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `ips_type` - The IPS type.

* `ips_version` - The IPS version.

* `update_time` - The update time.
