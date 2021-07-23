---
subcategory: "Cloud Data Migration (CDM)"
---

# sbercloud\_cdm\_flavors

Use this data source to get available SberCloud cdm flavors.

## Example Usage

```hcl
data "sbercloud_cdm_flavors" "flavor" {
}
```

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `region` - The region in which to obtain the CDM flavors. If omitted, the provider-level region will be used.

* `version` -
  The version of the flavor.

* `flavors` -
  Indicates the flavors information. Structure is documented below.

The `flavors` block contains:

* `name` - The name of the cdm flavor.
* `id` - The id of the cdm flavor.
