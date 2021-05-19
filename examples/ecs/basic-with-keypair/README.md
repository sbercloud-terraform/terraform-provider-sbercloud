## Example: Key Pair creation and then basic ECS instance using this key pair

This example provisions an ECS instance with the following attributes:
- availability zone: AZ1 ("ru-moscow-1a")
- flavor: s6.large.2
- OS: Ubuntu 20.04
- 1 security group: "sg-ssh"
- one system EVS disk of "High I/O" type (SAS) and 16 GB size

One needs to create a public key first (in Linux, you can do it with "ssh-keygen").

The example expects that one already has a subnet.  

**Please note how the list of availability zones is obtained.**

Please note the **sbercloud_compute_flavors** data source.  
It gets the right flavor name based on the performance type, vCPU and RAM amount required.  
The following performance types are available:
- "normal": General-purpose
- "computingv3": Dedicated general-purpose
- "highmem": Memory-optimized
- "diskintensive": Disk-intensive
- "gpu": GPU-accelerated

