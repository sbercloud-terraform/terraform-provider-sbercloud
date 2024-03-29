---
subcategory: "Elastic IP (EIP)"
---

# sbercloud_vpc_eip

Manages an EIP resource within SberCloud.

## Example Usage

### EIP with Dedicated Bandwidth

```hcl
resource "sbercloud_vpc_eip" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    share_type  = "PER"
    name        = "test"
    size        = 10
    charge_mode = "traffic"
  }
}
```

### EIP with Shared Bandwidth

```hcl
resource "sbercloud_vpc_bandwidth" "bandwidth_1" {
  name = "bandwidth_1"
  size = 5
}

resource "sbercloud_vpc_eip" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    share_type = "WHOLE"
    id         = sbercloud_vpc_bandwidth.bandwidth_1.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the eip resource. If omitted, the provider-level
  region will be used. Changing this creates a new eip resource.

* `publicip` - (Required, List) The elastic IP address object.

* `bandwidth` - (Required, List) The bandwidth object.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the elastic IP.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project id of the elastic IP. Changing this
  creates a new eip.

* `auto_renew` - (Optional, String, ForceNew) Specifies whether auto renew is enabled. Valid values are "true" and "
  false". Changing this creates a new resource.

The `publicip` block supports:

* `type` - (Required, String, ForceNew) The type of the eip. Changing this creates a new eip.

* `ip_address` - (Optional, String, ForceNew) The value must be a valid IP address in the available IP address segment.
  Changing this creates a new eip.

* `port_id` - (Optional, String) The port id which this eip will associate with. If the value is "" or this not
  specified, the eip will be in unbind state.

The `bandwidth` block supports:

* `share_type` - (Required, String, ForceNew) Whether the bandwidth is dedicated or shared. Changing this creates a new
  eip. Possible values are as follows:
  + *PER*: Dedicated bandwidth
  + *WHOLE*: Shared bandwidth

* `name` - (Optional, String) The bandwidth name, which is a string of 1 to 64 characters that contain letters, digits,
  underscores (_), and hyphens (-). This parameter is mandatory when `share_type` is set to *PER*.

* `size` - (Optional, Int) The bandwidth size. The value ranges from 1 to 300 Mbit/s. This parameter is mandatory
  when `share_type` is set to *PER*.

* `id` - (Optional, String, ForceNew) The shared bandwidth id. This parameter is mandatory when
  `share_type` is set to *WHOLE*. Changing this creates a new eip.

* `charge_mode` - (Optional, String, ForceNew) Specifies whether the bandwidth is billed by traffic or by bandwidth
  size. The value can be *traffic* or *bandwidth*. Changing this creates a new eip.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
* `address` - The IP address of the eip.
* `ipv6_address` - The IPv6 address of the EIP.
* `private_ip` - The private IP address bound to the EIP.
* `port_id` - The port ID which the EIP associated with.
* `status` - The status of eip.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minute.

## Import

EIPs can be imported using the `id`, e.g.

```
$ terraform import sbercloud_vpc_eip.eip_1 2c7f39f3-702b-48d1-940c-b50384177ee1
```
