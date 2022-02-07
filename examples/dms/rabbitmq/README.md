## Example: DMS Cluster (RabbitMQ)

### Requirements

- VPC exists in SberCloud.Advanced
- subnet exists in SberCloud.Advanced
- security group exists in SberCloud.Advanced

### Description

This example provisions a DMS cluster of RabbitMQ with the following attributes:

- engine: RabbitMQ 
- version: 3.7.17
- availability zones: ru-moscow-1a, ru-moscow-1b and ru-moscow-1c. That is, our RabbitMQ cluster is geo-redundant and distributed between 3 AZs.
- storage space: 1000 GB
- storage type: SSD
- user name: admin

### Notes

It will be a 5 nodes cluster.
