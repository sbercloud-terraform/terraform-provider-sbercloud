---
subcategory: "Distributed Message Service (DMS)"
---

# sbercloud\_dms\_az

Use this data source to get the ID of an available SberCloud DMS AZ.

## Example Usage

```hcl

data "sbercloud_dms_az" "az1" {
  code = "ru-moscow-1a"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the dms az. If omitted, the provider-level region will be used.

* `name` - (Optional, String) Indicates the name of an AZ.

* `code` - (Optional, String) Indicates the code of an AZ.

* `port` - (Optional, String) Indicates the port number of an AZ.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a data source ID in UUID format.

