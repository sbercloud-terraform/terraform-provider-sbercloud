## Example: Several ECS Instances, placed in different AZs

### Requirements

- key pair exists in SberCloud.Advanced
- subnet exists in SberCloud.Advanced

### Description

This example provisions ECS instances in different AZs:
- number of ECS: 4
- flavor: s6.large.2
- OS: Ubuntu 20.04
- 1 security group: "default"
- one system EVS disk of "High I/O" type (SAS) and 16 GB size. 

Availability Zone name is defined automatically for each ECS at the creation time.

### Notes

Please note how the list of availability zones is obtained.

Then the length of the list of availability zones is declared as a local variable.

Please note the "count = 4" line inside the "sbercloud_compute_instance" resource. This is the number of ECS to be created. Terraform [count meta-argument](https://www.terraform.io/docs/language/meta-arguments/count.html) allows creating several similar resources from one just block.  

Later the code refers to the current iteration with **count.index**. Current iteration is used to set a unique name for each ECS.  

Finally, note how each ECS is placed into a different AZ. The code doesn't depend on the names of AZs at all.
