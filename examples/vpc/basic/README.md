## Example: Basic VPC with subnets

### Requirements

None.

### Description

This example provisions basic VPC instance with two subnets.

### Notes

Please note the "primary_dns" attribute. It's set to 100.125.13.59  
**It's better if you keep this value.**  

This is an internal cloud DNS server, and it plays a very important role as it resolves internal requests to cloud services.  
If you change its value to something else for a subnet, you won't be able to use some cloud services in this subnet, such as CCE, SFS.  

If you really need to use your own DNS servers, please contact SberCloud representatives.
