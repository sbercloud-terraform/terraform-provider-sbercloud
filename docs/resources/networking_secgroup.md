---
subcategory: "Virtual Private Cloud (VPC)"
---

# sbercloud\_networking\_secgroup

Manages a V2 neutron security group resource within SberCloud.
Unlike Nova security groups, neutron separates the group from the rules
and also allows an admin to target a specific tenant_id.

## Example Usage

```hcl
resource "sbercloud_networking_secgroup" "secgroup_1" {
  name        = "secgroup_1"
  description = "My neutron security group"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to obtain the V2 networking client.
    A networking client is needed to create a port. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    security group.

* `name` - (Required, String) A unique name for the security group.

* `description` - (Optional, String) Description of the security group.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id of the security group.
  Changing this creates a new security group.

* `delete_default_rules` - (Optional, Bool, ForceNew) Whether or not to delete the default
    egress security rules. This is `false` by default. See the below note
    for more information.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `rules` - The array of security group rules associating with the security group.
  The [rule object](#security_group_rule) is documented below.

* `created_at` - The creation time, in UTC format.

* `updated_at` - The last update time, in UTC format.

<a name="security_group_rule"></a>
The `rules` block supports:

* `id` - The security group rule ID.
* `description` - The supplementary information about the security group rule.
* `direction` - The direction of the rule. The value can be *egress* or *ingress*.
* `ethertype` - The IP protocol version. The value can be *IPv4* or *IPv6*.
* `protocol` - The protocol type.
* `ports` - The port value range.
* `remote_ip_prefix` - The remote IP address. The value can be in the CIDR format or IP addresses.
* `remote_group_id` - The ID of the peer security group.
* `remote_address_group_id` - The ID of the remote address group.
* `action` - The effective policy.
* `priority` - The priority number.

## Default Security Group Rules

In most cases, SberCloud will create some egress security group rules for each
new security group. These security group rules will not be managed by
Terraform, so if you prefer to have *all* aspects of your infrastructure
managed by Terraform, set `delete_default_rules` to `true` and then create
separate security group rules such as the following:

```hcl
resource "sbercloud_networking_secgroup_rule" "secgroup_rule_v4" {
  direction         = "egress"
  ethertype         = "IPv4"
  security_group_id = sbercloud_networking_secgroup.secgroup.id
}

resource "sbercloud_networking_secgroup_rule" "secgroup_rule_v6" {
  direction         = "egress"
  ethertype         = "IPv6"
  security_group_id = sbercloud_networking_secgroup.secgroup.id
}
```

Please note that this behavior may differ depending on the configuration of
the SberCloud cloud. The above illustrates the current default Neutron
behavior. Some SberCloud clouds might provide additional rules and some might
not provide any rules at all (in which case the `delete_default_rules` setting
is moot).


## Timeouts
This resource provides the following timeouts configuration options:
- `delete` - Default is 10 minute.
## Import

Security Groups can be imported using the `id`, e.g.

```
$ terraform import sbercloud_networking_secgroup.secgroup_1 38809219-5e8a-4852-9139-6f461c90e8bc
```
