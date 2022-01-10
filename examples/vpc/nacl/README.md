## Example: NACL (Network Access Control List) with several rules

### Requirements

- subnet exists in SberCloud.Advanced

### Description

This example provisions a network access control list (NACL) with 2 inbound rules and 1 outbound rule.  
Rule details are defined in a local variables: *inbound_rules* and *outbound_rules*, respectively. Each has a *map* type.

### Notes

Please note the usage of the **for_each** meta-argument in the *sbercloud_networking_secgroup_rule* resource. It allows to iterate over a set or map object, and in this example it helps add several rules in just one resource.  
This way you can automate the creation of quite complex NACL with many rules. The rules content can be passed in external variables.  
Please also note how one iterates over rules created by **for_each**: the *for .. in ..* statement is very helpful.
