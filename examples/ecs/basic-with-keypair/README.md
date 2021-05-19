## Example: Key Pair creation and then basic ECS instance using this key pair

### Requirements

- public key is created (in Linux, you can do it with "ssh-keygen")
- subnet is created in SberCloud.Advanced

### Description

This example creates Key Pair and then provisions an ECS instance with the following attributes:
- availability zone: AZ1 ("ru-moscow-1a")
- flavor: s6.large.2
- OS: Ubuntu 20.04
- 1 security group: "sg-ssh"
- one system EVS disk of "High I/O" type (SAS) and 16 GB size
- key pair created from the public key

### Notes 

None.
