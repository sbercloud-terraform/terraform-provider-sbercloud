## Example: DMS Cluster (Kafka)

### Requirements

- VPC exists in SberCloud.Advanced
- subnet exists in SberCloud.Advanced
- security group exists in SberCloud.Advanced

### Description

This example provisions a DMS cluster of Kafka with the following attributes:

- engine: Kafka 
- version: 2.3.0
- availability zones: ru-moscow-1a, ru-moscow-1b and ru-moscow-1c. That is, our Kafka cluster is geo-redundant and distributed between 3 AZs.
- bandwidth: 300 MB/s
- storage space: 1200 GB
- storage type: SAS

Also, it creates a topic called *topic_01*.

### Notes

No notes so far.
