---
subcategory: "Deprecated"
---

# sbercloud_dcs_product

Use this data source to get the ID of an available DCS product.

!> **WARNING:** It has been deprecated. This data source is used for the `product_id` of the
`sbercloud_dcs_instance` resource. Now `product_id` has been deprecated and this data source is no longer used.

## Example Usage

```hcl
data "sbercloud_dcs_product" "product1" {
  spec_code = "dcs.single_node"
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to obtain the dcs products.
  If omitted, the provider-level region will be used.

* `spec_code` - (Optional, String) Specifies the DCS instance specification code. For details, see
    [Querying Service Specifications](https://support.hc.sbercloud.ru/api/dcs/dcs-api-0312040.html).
  + Log in to the DCS console, click *Buy DCS Instance*, and find the corresponding instance specification.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A data source ID in UUID format.
* `engine` - The cache engine. The value is *redis* or *memcached*.
* `engine_version` - The supported versions of a cache engine.
* `cache_mode` - The mode of a cache engine. The value is one of *single*, *ha*, *cluster*,
  *proxy* and *ha_rw_split*.
