## Example: Basic ECS Instance with tag

### Requirements

- key pair exists in SberCloud.Advanced
- subnet exists in SberCloud.Advanced

### Description

This example provisions a basic ECS instance with the following attributes, including tags:
- availability zone: AZ1 ("ru-moscow-1a")
- flavor: s6.large.2
- OS: Ubuntu 20.04
- security group: "default"
- one EVS disk of "High I/O" type (SAS) and 16 GB size
- tag is assigned to ECS: "created_by" = "terraform"

### Notes

None.
