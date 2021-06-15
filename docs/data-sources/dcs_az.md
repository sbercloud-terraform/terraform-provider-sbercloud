---
subcategory: "Distributed Cache Service"
---

# sbercloud\_dcs\_az

Use this data source to get the ID of an available SberCloud DCS AZ.

## Example Usage

```hcl

data "sbercloud_dcs_az" "az1" {
  port = "443"
  code = "ru-moscow-1a"
}
```

## Argument Reference

For details, See [Querying AZ Information](https://support.hc.sbercloud.ru/api/dcs/dcs-api-0312039.html).

* `region` - (Optional, String) The region in which to obtain the dcs az. If omitted, the provider-level region will be used.

* `name` - (Optional, String) Indicates the name of an AZ.

* `code` - (Optional, String) Indicates the code of an AZ.

* `port` - (Optional, String) Indicates the port number of an AZ.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a data source ID in UUID format.
