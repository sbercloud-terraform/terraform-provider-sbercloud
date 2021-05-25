## Example: Elastic Load Balancer (ELB) with 2 ECS in the backend servers group

### Requirements

- subnet exists in SberCloud.Advanced
- key pair exists in SberCloud.Advanced

### Description

This example creates two ECS, which will be behind the load balancer, in the backend servers group.  
For more details on the ECS creation see corresponding examples.  

Next, it creates EIP with the following attributes:
- bandwidth size: 3 Mbit/s
- billed by: bandwidth

Then it provisions Elastic Load Balancer (ELB) in the subnet.  

Next, it attached the EIP to the ELB, so that ELB becomes **public**.

Then it created a "listener" object, which will listen on TCP:80 (HTTP).

Next, it created an ECS backend group with the following attributes:
- backend protocol: HTTP
- load balancing method: Round Robin

Then it created a health check policy for servers in the backend with the following attributes:
- probe type: HTTP
- HTTP probing URL path: /
- expected HTTP code: 200-202 (range of 3)
- delay between probes: 10 seconds
- number of seconds to wait for a probe reply: 5 seconds
- number of allowed ping failures before changing the backend host status to "unhealth": 3

Finally, it adds both ECS into the backend servers group with the following attributes:
- backend protocol port: 80

### Notes

Please note the **count** meta-argument both in the "sbercloud_compute_instance" and "sbercloud_lb_member" resources. It helps handle several similar resources.  

Please note the **depends_on** meta-argument in the "sbercloud_lb_member" resource. It makes Terraform wait until the health check policy creation is completed, and only then adds the server into the backend group.
