---
subcategory: "Elastic Cloud Server (ECS)"
layout: "sbercloud"
page_title: "sbercloud: sbercloud_compute_flavor_sales_policies"
description: |-
  Use this data source to get the flavor sales policies of spot pricing ECSs and IES instances.
---

# sbercloud_compute_flavor_sales_policies

Use this data source to get the flavor sales policies of spot pricing ECSs and IES instances.

## Example Usage
а
```hcl
data "sbercloud_compute_flavor_sales_policies" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `flavor_id` - (Optional, String) Specifies the flavor ID of the ECS.

* `sell_status` - (Optional, String) Specifies the sales status of the ECS system flavor.
  Value options:
  + **available**: indicates that the flavor is available.
  + **sellout**: indicates that the flavor has been sold out.

* `sell_mode` - (Optional, String) Specifies the billing mode.
  Value options:
  + **postPaid**: indicates the pay-per-use billing mode, which is not supported currently.
  + **prePaid**: indicates the yearly/monthly billing mode, which is not supported currently.

* `availability_zone_id` - (Optional, String) Specifies the AZ.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `sell_policies` - Indicates the list of ECS flavor sales policies

  The [sell_policies](#sell_policies_struct) structure is documented below.

<a name="sell_policies_struct"></a>
The `sell_policies` block supports:

* `id` - Indicates the index of the ECS flavor.

* `flavor_id` - Indicates the ECS flavor ID.

* `sell_status` - Indicates the sales status of the ECS flavor.
  The value can be:
  + **sellout**: indicates that the flavor has been sold out.
  + **available**: indicates that the flavor is available.

* `availability_zone_id` - Indicates the AZ of the ECS flavor.

* `sell_mode` - Indicates the billing mode of the ECS flavor.
  The value can be:
  + **postPaid**: indicates the pay-per-use billing mode, which is not supported currently.
  + **prePaid**: indicates the yearly/monthly billing mode, which is not supported currently.