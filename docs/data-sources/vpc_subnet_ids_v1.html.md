---
layout: "sbercloud"
page_title: "SberCloud: sbercloud_subnet_ids_v1"
sidebar_current: "docs-sbercloud-datasource-subnet-ids-v1"
description: |-
  Provides a list of subnet Ids for a VPC
---

# Data Source: sbercloud_vpc_subnet_ids

`sbercloud_vpc_subnet_ids` provides a list of subnet ids for a vpc_id

This resource can be useful for getting back a list of subnet ids for a vpc.

## Example Usage

The following example shows outputing all cidr blocks for every subnet id in a vpc.

 ```hcl
data "sbercloud_vpc_subnet_ids" "subnet_ids" {
  vpc_id = var.vpc_id
}

data "sbercloud_vpc_subnet" "subnet" {
  count = length(data.sbercloud_vpc_subnet_ids.subnet_ids.ids)
  id    = tolist(data.sbercloud_vpc_subnet_ids.subnet_ids.ids)[count.index]
 }

output "subnet_cidr_blocks" {
  value = [for s in data.sbercloud_vpc_subnet.subnet: "${s.name}: ${s.id}: ${s.cidr}"]
}
 ```

## Argument Reference

The following arguments are supported:

* `vpc_id` (Required) - Specifies the VPC ID used as the query filter.

## Attributes Reference

The following attributes are exported:

* `ids` - A set of all the subnet ids found. This data source will fail if none are found.
