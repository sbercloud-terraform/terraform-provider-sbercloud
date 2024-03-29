---
subcategory: "NAT Gateway (NAT)"
---

# sbercloud\_nat\_dnat\_rule

Manages a Dnat rule resource within SberCloud Nat.

## Example Usage

### Dnat

```hcl
resource "sbercloud_nat_dnat_rule" "dnat_1" {
  floating_ip_id        = "2bd659ab-bbf7-43d7-928b-9ee6a10de3ef"
  nat_gateway_id        = "bf99c679-9f41-4dac-8513-9c9228e713e1"
  private_ip            = "10.0.0.12"
  protocol              = "tcp"
  internal_service_port = 993
  external_service_port = 242
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the dnat rule resource. If omitted, the provider-level region will be used. Changing this creates a new Dnat rule resource.

* `floating_ip_id` - (Required, String, ForceNew) Specifies the ID of the floating IP address.
  Changing this creates a new resource.

* `internal_service_port` - (Required, Int, ForceNew) Specifies port used by ECSs or BMSs
  to provide services for external systems. Changing this creates a new resource.

* `nat_gateway_id` - (Required, String, ForceNew) ID of the nat gateway this dnat rule belongs to.
   Changing this creates a new dnat rule.

* `port_id` - (Optional, String, ForceNew) Specifies the port ID of an ECS or a BMS.
  This parameter and private_ip are alternative. Changing this creates a
  new dnat rule.

* `private_ip` - (Optional, String, ForceNew) Specifies the private IP address of a
  user, for example, the IP address of a VPC for dedicated connection.
  This parameter and port_id are alternative.
  Changing this creates a new dnat rule.

* `protocol` - (Required, String, ForceNew) Specifies the protocol type. Currently,
  TCP, UDP, and ANY are supported.
  Changing this creates a new dnat rule.

* `external_service_port` - (Required, Int, ForceNew) Specifies port used by ECSs or
  BMSs to provide services for external systems.
  Changing this creates a new dnat rule.

* `internal_service_port_range` - (Optional, String) Specifies port range used by Floating IP provide services
  for external systems.  
  This parameter and `external_service_port_range` are mapped **1:1** in sequence(, ranges must have the same length).
  The valid value for range is **1~65535** and the port ranges can only be concatenated with the `-` character.

* `external_service_port_range` - (Optional, String) Specifies port range used by ECSs or BMSs to provide
  services for external systems.  
  This parameter and `internal_service_port_range` are mapped **1:1** in sequence(, ranges must have the same length).
  The valid value for range is **1~65535** and the port ranges can only be concatenated with the `-` character.  
  Required if `internal_service_port_range` is set.

* `description` - (Optional, String) Specifies the description of the DNAT rule.  
  The value is a string of no more than `255` characters, and angle brackets (<>) are not allowed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `created_at` - Dnat rule creation time.

* `status` - Dnat rule status.

* `floating_ip_address` - The actual floating IP address.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

Dnat can be imported using the following format:

```
$ terraform import sbercloud_nat_dnat_rule.dnat_1 f4f783a7-b908-4215-b018-724960e5df4a
```
