## Example: Basic ECS instance with password authentication

### Requirements

- subnet is created in SberCloud.Advanced

### Description

This example provisions an ECS instance with the following attributes:
- availability zone: AZ1 ("ru-moscow-1a")
- flavor: s6.large.2
- OS: Ubuntu 20.04
- 1 security group: "sg-ssh"
- one system EVS disk of "High I/O" type (SAS) and 16 GB size
- authentication mode is login/password (by default, the "root" user is used for Linux ECS)

### Notes 

If the password authentication is used, one can't do user data injection into ECS.
