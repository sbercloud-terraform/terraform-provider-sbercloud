## Example: ECS Instance of c6 flavor, with 2 security groups and 2 disks

### Requirements

- key pair exists in SberCloud.Advanced
- subnet exists in SberCloud.Advanced

### Description

This example provisions an ECS instance with the following attributes:
- availability zone: AZ2 ("ru-moscow-1b")
- flavor: c6.large.2
- OS: Ubuntu 20.04
- 2 security groups: "default", "sg-ssh"
- one system EVS disk of "High I/O" type (SAS) and 16 GB size
- one data EVS disk of "High I/O" type (SAS) and 32 GB size

### Notes

Please note the **sbercloud_compute_flavors** data source.  
It gets the right flavor name based on the performance type, vCPU and RAM amount required.  
The following performance types are available:
- "normal": General-purpose
- "computingv3": Dedicated general-purpose
- "highmem": Memory-optimized
- "diskintensive": Disk-intensive
- "gpu": GPU-accelerated

