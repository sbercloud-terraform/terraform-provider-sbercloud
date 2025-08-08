---
subcategory: "Enterprise Router (ER)"
layout: "sbercloud"
page_title: "SberCloud: sbercloud_er_quotas"
description: |-
  Using this data source to query the list of available resource quotas within SberCloud.
---

# sbercloud_er_quotas

Using this data source to query the list of available resource quotas within SberCloud.

~> Using an invalid ID to filter the results will not report an error or return an empty list, but will return a quota
   list with all usage equal to 0.

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
data "sbercloud_er_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) The quota type to be queried.
  The valid values are as follows:
  + **er_instance**: Quotas and usage for enterprise router instances.
  + **dc_attachment**: Quotas and usage for DC attachment.
  + **vpc_attachment**: Quotas and usage for VPC attachment.
  + **vpn_attachment**: Quotas and usage for VPN attachment.
  + **peering_attachment**: Quotas and usage for peering attachment.
  + **can_attachment**: Quotas and usage for can attachment.
  + **route_table**: Quotas and usage for route table.
  + **static_route**: Quotas and usage for static route.
  + **vpc_er**: The number of enterprise routers that each VPC can access and the current usage.
  + **flow_log**: The number of flow logs that can be created per attachment.

* `instance_id` - (Optional, String) The instance ID.

* `route_table_id` - (Optional, String) The route table ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - All quotas that match the filter parameters.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `used` - The number of quota used.

* `unit` - The unit of usage.

* `type` - The quota type.

* `limit` - The number of available quotas, `-1` means unlimited.
