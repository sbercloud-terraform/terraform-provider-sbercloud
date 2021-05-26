## Example: Security group with many rules

### Requirements

None.

### Description

This example provisions a security group with many rules.  
Rule details are defined in a local variable "rules", which has a "map" type.

### Notes

Please note the usage of the **for_each** meta-argument in the "sbercloud_networking_secgroup_rule" resource. It allows to iterate over a set or map object, and in this example it helps add several rules in just one resource.
