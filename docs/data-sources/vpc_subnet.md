---
subcategory: "Virtual Private Cloud (VPC)"
---

# sbercloud_vpc_subnet

This resource can prove useful when a module accepts a subnet id as an input variable and needs to, for example,
determine the id of the VPC that the subnet belongs to.

## Example Usage

```hcl
data "sbercloud_vpc_subnet" "subnet" {
  id = var.subnet_id
}

output "subnet_vpc_id" {
  value = data.sbercloud_vpc_subnet.subnet.vpc_id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available subnets in the current tenant. The given
filters must match exactly one subnet whose data will be exported as attributes.

* `region` - (Optional, String) Specifies the region in which to obtain the subnet. If omitted, the provider-level
  region will be used.

* `id` - (Optional, String) - Specifies a resource ID in UUID format.

* `name` - (Optional, String) Specifies the name of the specific subnet to retrieve.

* `cidr` - (Optional, String) Specifies the network segment of specific subnet to retrieve. The value must be in CIDR
  format.

* `status` - (Optional, String) Specifies the value can be ACTIVE, DOWN, UNKNOWN, or ERROR.

* `vpc_id` - (Optional, String) Specifies the id of the VPC that the desired subnet belongs to.

* `gateway_ip` - (Optional, String) Specifies the subnet gateway address of specific subnet.

* `primary_dns` - (Optional, String) Specifies the IP address of DNS server 1 on the specific subnet.

* `secondary_dns` - (Optional, String) Specifies the IP address of DNS server 2 on the specific subnet.

* `availability_zone` - (Optional, String) Specifies the availability zone (AZ) to which the subnet should belong.

## **Attributes Reference**

In addition to all arguments above, the following attributes are exported:

* `dns_list` - The IP address list of DNS servers on the subnet.

* `dhcp_enable` - Whether the DHCP is enabled.

* `subnet_id` - The subnet (Native OpenStack API) ID.
* 