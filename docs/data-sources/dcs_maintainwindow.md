---
subcategory: "Distributed Cache Service"
---

# sbercloud\_dcs\_maintainwindow

Use this data source to get the ID of an available SberCloud DCS maintainwindow.

## Example Usage

```hcl

data "sbercloud_dcs_maintainwindow" "maintainwindow1" {
  seq = 1
}

```

## Argument Reference

For details, See [Querying Maintenance Time Window](https://support.hc.sbercloud.ru/api/dcs/dcs-api-0312041.html).

* `region` - (Optional, String) The region in which to obtain the dcs maintainwindows. If omitted, the provider-level region will be used.

* `seq` - (Required, Int) Indicates the sequential number of a maintenance time window.

* `begin` - (Optional, String) Indicates the time at which a maintenance time window starts.

* `end` - (Required, String) Indicates the time at which a maintenance time window ends.

* `default` - (Required, Bool) Indicates whether a maintenance time window is set to the default time segment.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a data source ID in UUID format.
