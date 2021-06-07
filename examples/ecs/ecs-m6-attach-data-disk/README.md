## Example: ECS Instance of m6 flavor. Attaching data disk

### Requirements

- key pair exists in SberCloud.Advanced
- subnet exists in SberCloud.Advanced

### Description

This example provisions an ECS instance with the following attributes:
- availability zone: AZ1 ("ru-moscow-1a")
- flavor: m6.xlarge.8
- OS: Ubuntu 20.04
- 1 security group: "default"
- one system EVS disk of "High I/O" type (SAS) and 16 GB size

Then an EVS disk is created and attached to the ECS:
- disk type: Ultra-high I/O (SSD)
- size: 64 GB

### Notes

Please note the **sbercloud_compute_flavors** data source.  
It gets the right flavor name based on the performance type, vCPU and RAM amount required.  
The following performance types are available:
- "normal": General-purpose
- "computingv3": Dedicated general-purpose
- "highmem": Memory-optimized
- "diskintensive": Disk-intensive
- "gpu": GPU-accelerated

