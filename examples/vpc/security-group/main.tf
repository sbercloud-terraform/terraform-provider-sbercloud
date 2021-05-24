# Create security group
resource "sbercloud_networking_secgroup" "sg_01" {
  name        = "sg-demo"
  description = "Security group for HTTP"
}

# Create a security group rule, which allows HTTP traffic
resource "sbercloud_networking_secgroup_rule" "sg_rule_http" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 80
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = sbercloud_networking_secgroup.sg_01.id
}
