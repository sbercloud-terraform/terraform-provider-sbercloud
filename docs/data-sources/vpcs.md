---
subcategory: "Virtual Private Cloud (VPC)"
---

# sbercloud_vpcs

Use this data source to get a list of VPC.

## Example Usage

An example filter by name and tag

```hcl
variable "vpc_name" {}

data "sbercloud_vpcs" "vpc" {
  name = var.vpc_name
  tags {
    foo = "bar,value"
  }
}

output "vpc_ids" {
  value = data.sbercloud_vpcs.vpc.vpcs[*].id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available VPCs in the current region.
 All VPCs that meet the filter criteria will be exported as attributes.

* `region` - (Optional, String) Specifies the region in which to obtain the VPC. If omitted, the provider-level region
  will be used.

* `id` - (Optional, String) Specifies the id of the desired VPC.

* `name` - (Optional, String) Specifies the name of the desired VPC. The value is a string of no more than 64 characters
  and can contain digits, letters, underscores (_) and hyphens (-).

* `status` - (Optional, String) Specifies the current status of the desired VPC. The value can be CREATING, OK or ERROR.

* `cidr` - (Optional, String) Specifies the cidr block of the desired VPC.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID which the desired VPC belongs to.

* `tags` - (Optional, Map) Specifies the included key/value pairs which associated with the desired VPC.

 -> A maximum of 10 tag keys are allowed for each query operation. Each tag key can have up to 10 tag values.
  The tag key cannot be left blank or set to an empty string. Each tag key must be unique, and each tag value in a
  tag must be unique, use commas(,) to separate the multiple values. An empty for values indicates any value.
  The values are in the OR relationship.

## Attributes Reference

The following attributes are exported:

* `id` - Indicates a data source ID.
* `vpcs` - Indicates a list of all VPCs found. Structure is documented below.

The `vpcs` block supports:

* `id` - Indicates the ID of the VPC.
* `name` - Indicates the name of the VPC.
* `cidr` - Indicates the cidr block of the VPC.
* `status` - Indicates the current status of the VPC.
* `enterprise_project_id` - Indicates the the enterprise project ID of the VPC.
* `description` - Indicates the description of the VPC.
* `tags` - Indicates the key/value pairs which associated with the VPC.
