## Example: Basic ECS Instance

This example provisions a basic ECS instance with the following attributes:
- availability zone: AZ1 ("ru-moscow-1a")
- flavor: s6.large.2
- OS: Ubuntu 20.04
- security group: "default"
- one EVS disk of "High I/O" type (SAS) and 16 GB size

The example expects that one already has a key pair and a subnet.  
