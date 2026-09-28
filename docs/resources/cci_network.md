---
subcategory: "Cloud Container Instance (CCI)"
layout: "sbercloud"
page_title: "Sbercloud: sbercloud_cci_network"
description: |-
  Manages a CCI v2 network resource within Sbercloud.
---

# sbercloud_cci_network

Manages a CCI v2 network resource within Sbercloud.

## Example Usage

```hcl
variable "namespace_name" {}
variable "network_name" {}
variable "subnet_id" {}
variable "security_group_id" {}

resource "sbercloud_cci_network" "test" {
  namespace = var.namespace_name
  name      = var.network_name

  subnets {
    subnet_id = var.subnet_id
  }

  security_group_ids = [var.security_group_id]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CCI network.
  If omitted, the provider-level region will be used. Changing this will create a new CCI network resource.

* `namespace` - (Required, String, NonUpdatable) Specifies the namespace of the CCI network.

* `name` - (Required, String, NonUpdatable) Specifies the name of the CCI network.

* `annotations` - (Optional, Map) Specifies the annotations of the CCI network.

* `ip_families` - (Optional, List, NonUpdatable) Specifies the IP families of the CCI network.

* `security_group_ids` - (Optional, List) Specifies the security group IDs of the CCI network.

* `subnets` - (Optional, List, NonUpdatable) Specifies the subnets of the CCI network.
  The [subnets](#cci_network_subnets) structure is documented below.

<a name="cci_network_subnets"></a>
The `subnets` block supports:

* `subnet_id` - (Optional, String) Specifies the IPv4 subnet ID of the VPC subnet.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `<namespace>/<name>`.

* `api_version` - The API version of the CCI network.

* `kind` - The kind of the CCI network.

* `creation_timestamp` - The creation timestamp of the CCI network.

* `finalizers` - The finalizers of the CCI network.

* `resource_version` - The resource version of the CCI network.

* `uid` - The uid of the CCI network.

* `status` - The status of the CCI network.
  The [status](#cci_network_status) structure is documented below.

<a name="cci_network_status"></a>
The `status` block supports:

* `status` - The status of the CCI network.

* `conditions` - The conditions of the CCI network.
  The [conditions](#cci_network_status_conditions) structure is documented below.

* `subnet_attrs` - The subnet attributes of the CCI network.
  The [subnet_attrs](#cci_network_status_subnet_attrs) structure is documented below.

<a name="cci_network_status_conditions"></a>
The `conditions` block supports:

* `type` - The type of the condition.

* `status` - The status of the condition.

* `last_transition_time` - The last transition time of the condition.

* `reason` - The reason of the condition.

* `message` - The message of the condition.

<a name="cci_network_status_subnet_attrs"></a>
The `subnet_attrs` block supports:

* `network_id` - The network ID of the VPC subnet.

* `subnet_v4_id` - The IPv4 subnet ID.

* `subnet_v6_id` - The IPv6 subnet ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

CCI networks can be imported using their `namespace` and `name`, separated by a slash, e.g.

```bash
$ terraform import sbercloud_cci_network.test <namespace>/<name>
```
