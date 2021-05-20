## Example: Creating ECS Instance with 2 NICs. Then attaching one more NIC

### Requirements

- key pair exists in SberCloud.Advanced
- 3 subnets exist in SberCloud.Advanced

### Description

This example provisions an ECS instance with the following attributes:
- availability zone: AZ2 ("ru-moscow-1b")
- flavor: c6.xlarge.2
- OS: Ubuntu 20.04
- 1 security group: "default"
- one system EVS disk of "High I/O" type (SAS) and 16 GB size
- two NICs: attaching ECS to your first and second subnets during ECS creation

After creation, the ECS is attached to the third subnet.

### Notes

Please note how the list of availability zones is obtained.  

Please note the **sbercloud_compute_flavors** data source.  
It gets the right flavor name based on the performance type, vCPU and RAM amount required.  
The following performance types are available:
- "normal": General-purpose
- "computingv3": Dedicated general-purpose
- "highmem": Memory-optimized
- "diskintensive": Disk-intensive
- "gpu": GPU-accelerated

