## Example: DCS Cluster (Redis)

### Requirements

- VPC exists in SberCloud.Advanced
- subnet exists in SberCloud.Advanced

### Description

This example provisions a DCS cluster of Redis with the following attributes:

- flavor: redis.cluster.xu1.large.r2.8
- availability zones: ru-moscow-1a, ru-moscow-1b. That is, our Redis cluster is geo-redundant.
- Engine: Redis 5.0
- Cache size: 8 GB
- Password protected

### Notes

You can get the list of flavor names by execuring [this API call](https://support.hc.sbercloud.ru/api/dcs/dcs-api-0312040.html) or by looking at the DCS instance creation process in the console.  
  
Note the tag attached to the cluster.  
Note the backup schedule configured for the cluster. Backups will be performed each Tuesday, Thursday and Saturday, at **05:00 MSK** (in main.tf time is GMT), and stored for 5 days. 
