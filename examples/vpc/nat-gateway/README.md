## Example: NAT Gateway with SNAT rule

### Requirements

- VPC exists in SberCloud.Advanced
- subnet exists in SberCloud.Advanced

### Description

This example creates EIP with the following attributes:
- bandwidth size: 4 Mbit/s
- billed by: "bandwidth"

Then it provisions NAT Gateway with the following attributes:
- type: "small"

Finally, it creates an SNAT rule, which attached the subnet to NAT Gateway.

### Notes

None.
