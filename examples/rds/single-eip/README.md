## Example: Single RDS Instance with EIP (PostgreSQL)

### Requirements

- VPC exists in SberCloud.Advanced
- subnet exists in SberCloud.Advanced
- security group exists in SberCloud.Advanced

### Description

This example provisions a single RDS instance of PostgreSQL with an EIP attached to it, with the following other attributes:

- flavor: I don't know and I don't care :) All I know is I need 2 vCPUs and 4 GB RAM. See Notes below.
- availability zone: ru-moscow-1b
- PG version: 13
- Disk storage size: 100 GB
- Disk storage type: SSD
- EIP charge mode: by bandwidth
- EIP bandwidth size: 4 Mbit/s

### Notes

Please note the **sbercloud_rds_flavors** data source.
It gets the right flavor name based on the number of vCPU and RAM amount. It helps avoid setting flavor names explicitly.  
