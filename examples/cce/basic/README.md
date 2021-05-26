## Example: Basic CCE Cluster

### Requirements

- VPC exists in SberCloud.Advanced
- subnet exists in SberCloud.Advanced, and the primary DNS of this subnet is **100.125.13.59**: this is very important! If the primary DNS is set to another value, the CCE creation can fail.
- key pair exists in SberCloud.Advanced

### Description

This example provisions a basic CCE cluster with the following attributes:
- Kubernetes version: latest available
- number of master nodes: 3
- master nodes are spread over several AZs
- cluster size: up to 50 worker nodes
- network model: tunnel

Also this example provisions one worker node. The process of creating workers is very similar to ECS.

### Notes

Please note the "flavor_id" argument in the "sbercloud_cce_cluster" resource. It defines the number and size of master nodes and can have the following values:

- cce.s1.small - small-scale cluster (up to 50 nodes) with 1 master node.
- cce.s1.medium - medium-scale cluster (up to 200 nodes) with 1 master node.
- cce.s1.large - large-scale cluster (up to 1000 nodes) with 1 master node.
- cce.s2.small - small-scale HA cluster (up to 50 nodes) with 3 master nodes.
- cce.s2.medium - medium-scale HA cluster (up to 200 nodes) with 3 master nodes.
- cce.s2.large - large-scale HA cluster (up to 1000 nodes) with 3 master nodes.

Please note the "root_volume" and "data_volumes" arguments in the "sbercloud_cce_node" resource. The minimum size of the CCE node root volume is 50 GB. CCE node must also have at least one data disk, and its minimum size is 100 GB.
