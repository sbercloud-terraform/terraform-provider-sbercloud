## Example: Primary-Standby RDS Instance (PostgreSQL)

### Requirements

- VPC exists in SberCloud.Advanced
- subnet exists in SberCloud.Advanced
- security group exists in SberCloud.Advanced

### Description

This example provisions a primary-standby RDS instance of PostgreSQL with the following attributes:

- flavor: I don't know and I don't care :) All I know is there are 2 vCPUs and 8 GB RAM. See Notes below.
- availability zones: ru-moscow-1a and ru-moscow-1b. That is, my PG is geo-redundant.
- PG version: 12
- Disk storage size: 100 GB
- Disk storage type: SSD

### Notes

Please note the **sbercloud_rds_flavors** data source.
It gets the right flavor name based on the number of vCPU and RAM amount. It helps avoid setting flavor names explicitly.  

The **ha_replication_mode** parameter is described as Optional, but it's better to set it explicitly for HA configurations.  

Please note that there is a backup policy configured and a tag attached to the instance.
