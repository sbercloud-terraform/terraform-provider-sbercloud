---
subcategory: "Virtual Private Cloud (VPC)"
---

# sbercloud_networking_secgroup_rule

Manages a Security Group Rule resource within SberCloud.

## Example Usage

### Create an ingress rule that opens TCP port 8080 with port range parameters

```hcl
variable "security_group_id" {}

resource "sbercloud_networking_secgroup_rule" "test" {
  security_group_id = var.security_group_id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 8080
  port_range_max    = 8080
  remote_ip_prefix  = "0.0.0.0/0"
}
```

### Create an ingress rule that enable the remote address group and open some TCP ports

```hcl
variable "group_name" {}
variable "security_group_id" {}

resource "sbercloud_vpc_address_group" "test" {
  name = var.group_name

  addresses = [
    "192.168.10.12",
    "192.168.11.0-192.168.11.240",
  ]
}

resource "sbercloud_networking_secgroup_rule" "test" {
  security_group_id       = var.security_group_id
  direction               = "ingress"
  action                  = "allow"
  ethertype               = "IPv4"
  ports                   = "80,500,600-800"
  protocol                = "tcp"
  priority                = 5
  remote_address_group_id = sbercloud_vpc_address_group.test.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to obtain the V2 networking client.
    A networking client is needed to create a port. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    security group rule.

* `direction` - (Required, String, ForceNew) The direction of the rule, valid values are __ingress__
    or __egress__. Changing this creates a new security group rule.

* `ethertype` - (Required, String, ForceNew) The layer 3 protocol type, valid values are __IPv4__
    or __IPv6__. Changing this creates a new security group rule.

* `description` - (Optional, String, ForceNew) Specifies the supplementary information about the networking security
  group rule. This parameter can contain a maximum of 255 characters and cannot contain angle brackets (< or >).
  Changing this creates a new security group rule.

* `protocol` - (Optional, String, ForceNew) The layer 4 protocol type, valid values are following. Changing this creates a new security group rule. This is required if you want to specify a port range.
  * __tcp__
  * __udp__
  * __icmp__
  * __ah__
  * __dccp__
  * __egp__
  * __esp__
  * __gre__
  * __igmp__
  * __ipv6-encap__
  * __ipv6-frag__
  * __ipv6-icmp__
  * __ipv6-nonxt__
  * __ipv6-opts__
  * __ipv6-route__
  * __ospf__
  * __pgm__
  * __rsvp__
  * __sctp__
  * __udplite__
  * __vrrp__

* `port_range_min` - (Optional, String, ForceNew) The lower part of the allowed port range, valid
    integer value needs to be between 1 and 65535. This parameter and `ports` are alternative. Changing this creates a new
    security group rule.

* `port_range_max` - (Optional, Int, ForceNew) The higher part of the allowed port range, valid
    integer value needs to be between 1 and 65535. This parameter and `ports` are alternative. Changing this creates a new
    security group rule.

* `ports` - (Optional, String, ForceNew) Specifies the allowed port value range, which supports single port (80),
  continuous port (1-30) and discontinous port (22, 3389, 80) The valid port values is range form `1` to `65,535`.
  Changing this creates a new security group rule.

* `remote_ip_prefix` - (Optional, String, ForceNew) The remote CIDR, the value needs to be a valid
    CIDR (i.e. 192.168.0.0/16). Changing this creates a new security group rule.

* `remote_group_id` - (Optional, String, ForceNew) The remote group id, the value needs to be an
    Openstack ID of a security group in the same tenant. Changing this creates
    a new security group rule.

* `remote_address_group_id` - (Optional, String, ForceNew) Specifies the remote address group ID.
  This parameter is not used with `port_range_min` and `port_range_max`.
  Changing this creates a new security group rule.

* `security_group_id` - (Required, String, ForceNew) The security group id the rule should belong
    to, the value needs to be an Openstack ID of a security group in the same
    tenant. Changing this creates a new security group rule.

* `action` - (Optional, String, ForceNew) Specifies the effective policy. The valid values are **allow** and **deny**.
  This parameter is not used with `port_range_min` and `port_range_max`.
  Changing this creates a new security group rule.

* `priority` - (Optional, Int, ForceNew) Specifies the priority number.
  The valid value is range from **1** to **100**. The default value is **1**.
  This parameter is not used with `port_range_min` and `port_range_max`.
  Changing this creates a new security group rule.



* `security_group_id` - (Required, String, ForceNew) Specifies the security group ID the rule should belong to. Changing
  this creates a new security group rule.

* `direction` - (Required, String, ForceNew) Specifies the direction of the rule, valid values are **ingress** or
  **egress**. Changing this creates a new security group rule.

* `ethertype` - (Required, String, ForceNew) Specifies the layer 3 protocol type, valid values are **IPv4** or **IPv6**.
  Changing this creates a new security group rule.

* `description` - (Optional, String, ForceNew) Specifies the supplementary information about the networking security
  group rule. This parameter can contain a maximum of 255 characters and cannot contain angle brackets (< or >).
  Changing this creates a new security group rule.

* `protocol` - (Optional, String, ForceNew) Specifies the layer 4 protocol type, valid values are **tcp**, **udp**,
  **icmp** and **icmpv6**. If omitted, the protocol means that all protocols are supported.
  This is required if you want to specify a port range. Changing this creates a new security group rule.

* `port_range_min` - (Optional, Int, ForceNew) Specifies the lower part of the allowed port range, valid integer value
  needs to be between `1` and `65,535`. Changing this creates a new security group rule.
  This parameter and `ports` are alternative.

* `port_range_max` - (Optional, Int, ForceNew) Specifies the higher part of the allowed port range, valid integer value
  needs to be between `1` and `65,535`. Changing this creates a new security group rule.
  This parameter and `ports` are alternative.

* `ports` - (Optional, String, ForceNew) Specifies the allowed port value range, which supports single port (80),
  continuous port (1-30) and discontinous port (22, 3389, 80) The valid port values is range form `1` to `65,535`.
  Changing this creates a new security group rule.

* `remote_ip_prefix` - (Optional, String, ForceNew) Specifies the remote CIDR, the value needs to be a valid CIDR (i.e.
  192.168.0.0/16). Changing this creates a new security group rule.

* `remote_group_id` - (Optional, String, ForceNew) Specifies the remote group ID. Changing this creates a new security
  group rule.

* `remote_address_group_id` - (Optional, String, ForceNew) Specifies the remote address group ID.
  This parameter is not used with `port_range_min` and `port_range_max`.
  Changing this creates a new security group rule.

* `action` - (Optional, String, ForceNew) Specifies the effective policy. The valid values are **allow** and **deny**.
  This parameter is not used with `port_range_min` and `port_range_max`.
  Changing this creates a new security group rule.

* `priority` - (Optional, Int, ForceNew) Specifies the priority number.
  The valid value is range from **1** to **100**. The default value is **1**.
  This parameter is not used with `port_range_min` and `port_range_max`.
  Changing this creates a new security group rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Timeouts

This resource provides the following timeouts configuration options:

* `delete` - Default is 10 minute.

## Import

Security Group Rules can be imported using the `id`, e.g.

```
$ terraform import sbercloud_networking_secgroup_rule.secgroup_rule_1 aeb68ee3-6e9d-4256-955c-9584a6212745
```
