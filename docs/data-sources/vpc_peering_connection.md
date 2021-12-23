---
subcategory: "Virtual Private Cloud (VPC)"
---

# sbercloud_vpc_peering_connection

The VPC Peering Connection data source provides details about a specific VPC peering connection.

## Example Usage

```hcl
data "sbercloud_vpc" "vpc" {
  name = "vpc"
}

data "sbercloud_vpc" "peer_vpc" {
  name = "peer_vpc"
}

data "sbercloud_vpc_peering_connection" "peering" {
  vpc_id      = data.sbercloud_vpc.vpc.id
  peer_vpc_id = data.sbercloud_vpc.peer_vpc.id
}

resource "sbercloud_vpc_route" "vpc_route" {
  type        = "peering"
  nexthop     = data.sbercloud_vpc_peering_connection.peering.id
  destination = "192.168.0.0/16"
  vpc_id      = data.sbercloud_vpc.vpc.id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available VPC peering connection. The given filters
must match exactly one VPC peering connection whose data will be exported as attributes.

* `region` - (Optional, String) The region in which to obtain the VPC Peering Connection. If omitted, the provider-level
  region will be used.

* `id` - (Optional, String) The ID of the specific VPC Peering Connection to retrieve.

* `status` - (Optional, String) The status of the specific VPC Peering Connection to retrieve.

* `vpc_id` - (Optional, String) The ID of the requester VPC of the specific VPC Peering Connection to retrieve.

* `peer_vpc_id` - (Optional, String) The ID of the accepter/peer VPC of the specific VPC Peering Connection to retrieve.

* `peer_tenant_id` - (Optional, String) The Tenant ID of the accepter/peer VPC of the specific VPC Peering Connection to
  retrieve.

* `name` - (Optional, String) The name of the specific VPC Peering Connection to retrieve.
