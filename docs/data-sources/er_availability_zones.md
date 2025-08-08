---
subcategory: "Enterprise Router (ER)"
layout: "sbercloud"
page_title: "SberCloud: sbercloud_er_availability_zones"
description: ""
---

# sbercloud_er_availability_zones

Use this data source to query availability zones where ER instances can be created within SberCloud.

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
data "sbercloud_er_availability_zones" "all" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are supported:

* `id` - The data source ID.

* `names` - The names of availability zone.
