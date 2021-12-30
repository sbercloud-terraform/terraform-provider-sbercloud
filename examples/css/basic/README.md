## Example: CSS Instance

### Requirements

- VPC exists in SberCloud.Advanced
- subnet exists in SberCloud.Advanced
- security group exists in SberCloud.Advanced

### Description

This example provisions a CSS instance with the following attributes:

- flavor: ess.spec-4u32g (which corresponds to m6.xlarge.8)
- number of nodes: 1
- availability zone: ru-moscow-1a
- Elasticsearch version: 7.9.3
- Disk storage size: 80 GB
- Disk storage type: SSD

### Notes

Please note that there is a backup policy configured and a tag attached to the instance.  
Backups (snapshots) will be performed daily at 01:00 MSK, put into the "p-test-02" bucket, into the "css_backups/css-terraform" folder, and stored for 4 days.
The "css_obs_agency" will be used to store backups in OBS.
