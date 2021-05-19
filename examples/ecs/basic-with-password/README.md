## Example: Basic ECS instance with password authentication

### Requirements

- root password has been salted as described [here](https://support.hc.sbercloud.ru/en-us/api/cce/cce_02_0242.html) 
- subnet is created in SberCloud.Advanced

### Description

This example provisions an ECS instance with the following attributes:
- availability zone: AZ1 ("ru-moscow-1a")
- flavor: s6.large.2
- OS: Ubuntu 20.04
- 1 security group: "sg-ssh"
- one system EVS disk of "High I/O" type (SAS) and 16 GB size

### Notes 

None.
