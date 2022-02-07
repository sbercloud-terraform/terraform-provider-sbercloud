## Example: Node Pool for CCE Cluster

### Requirements

- CCE cluster (master(s)) exists in SberCloud.Advanced
- key pair exists in SberCloud.Advanced

### Description

This example provisions a node pool for CCE cluster with the following attributes:

- Node flavor: s6.xlarge.4
- Minimal number of nodes in the pool: 2
- Initial number of nodes in the pool: 2
- Maximum number of nodes in the pool: 10
- Availability zone: ru-moscow-1a

As a result, this example provisions two worker nodes. 

### Notes

The **os** paramter is described as Optional, but it's better to set it explicitly to CentOS 7.6  
It may simplify upgrades to next CCE releases.
