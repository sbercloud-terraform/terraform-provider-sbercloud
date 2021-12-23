---
subcategory: "Virtual Private Cloud (VPC)"
---

# sbercloud_vpc_peering_connection_accepter

Provides a resource to manage the accepter's side of a VPC Peering Connection. This is an alternative
to `sbercloud_vpc_peering_connection_accepter_v2`

-> **NOTE:** When a cross-tenant (requester's tenant differs from the accepter's tenant) VPC Peering Connection
  is created, a VPC Peering Connection resource is automatically created in the accepter's account.
  The requester can use the `sbercloud_vpc_peering_connection` resource to manage its side of the connection and
  the accepter can use the `sbercloud_vpc_peering_connection_accepter` resource to accept its side of the connection
  into management.

## Example Usage

```hcl
provider "sbercloud" {
  alias = "main"
}

provider "sbercloud" {
  alias = "peer"
}

resource "sbercloud_vpc" "vpc_main" {
  provider = "sbercloud.main"
  name     = var.vpc_name
  cidr     = var.vpc_cidr
}

resource "sbercloud_vpc" "vpc_peer" {
  provider = "sbercloud.peer"
  name     = var.peer_vpc_name
  cidr     = var.peer_vpc_cidr
}

# Requester's side of the connection.
resource "sbercloud_vpc_peering_connection" "peering" {
  provider       = "sbercloud.main"
  name           = var.peer_name
  vpc_id         = sbercloud_vpc.vpc_main.id
  peer_vpc_id    = sbercloud_vpc.vpc_peer.id
  peer_tenant_id = var.tenant_id
}

# Accepter's side of the connection.
resource "sbercloud_vpc_peering_connection_accepter" "peer" {
  provider = "sbercloud.peer"
  accept   = true

  vpc_peering_connection_id = sbercloud_vpc_peering_connection.peering.id
}
 ```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the vpc peering connection accepter. If omitted,
  the provider-level region will be used. Changing this creates a new VPC peering connection accepter resource.

* `vpc_peering_connection_id` (Required, String, ForceNew) - The VPC Peering Connection ID to manage. Changing this
  creates a new VPC peering connection accepter.

* `accept` (Optional, Bool)- Whether or not to accept the peering request. Defaults to `false`.

## Removing sbercloud_vpc_peering_connection_accepter from your configuration

SberCloud allows a cross-tenant VPC Peering Connection to be deleted from either the requester's or accepter's side.
However, Terraform only allows the VPC Peering Connection to be deleted from the requester's side by removing the
corresponding `sbercloud_vpc_peering_connection` resource from your configuration.
Removing a `sbercloud_vpc_peering_connection_accepter` resource from your configuration will remove it from your
state file and management, but will not destroy the VPC Peering Connection.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The VPC peering connection name.

* `id` - The VPC peering connection ID.

* `status` - The VPC peering connection status.

* `vpc_id` - The ID of requester VPC involved in a VPC peering connection.

* `peer_vpc_id` - The VPC ID of the accepter tenant.

* `peer_tenant_id` - The Tenant Id of the accepter tenant.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `delete` - Default is 10 minute.
