---
subcategory: "Direct Connect (DC)"
layout: "sbercloud"
page_title: "SberCloud: sbercloud_dc_tags"
description: |-
  Use this data source to query tags under specified resource type of DC service.
---

# sbercloud_dc_tags

Use this data source to query tags under specified resource type of DC service.

## Example Usage

```hcl
data "sbercloud_dc_tags" "test" {
  resource_type = "dc-vif"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Resource type of DC service. The valid values are:
  + **dc-directconnect**: Direct connect connection.
  + **dc-vgw**: Virtual gateway.
  + **dc-vif**: Virtual interface.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of DC tags.
  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `values` - The values of the tag.
