# Set local variable of the "map" type.
# It contains details of the rules to be created
locals {
  rules = {
    http-rule = {
      description = "Allow HTTP from anywhere",
      protocol = "tcp",
      port = 80,
      source = "0.0.0.0/0"
    },
    ssh-rule = {
      description = "Allow SSH from only one source",
      protocol = "tcp",
      port = 22,
      source = "10.20.30.40/32"
    }
  }
}

# Create security group
resource "sbercloud_networking_secgroup" "sg_01" {
  name        = "sg-demo"
  description = "Security group with many rules"
}

# Create all security group rules in one go
resource "sbercloud_networking_secgroup_rule" "sg_rule_01" {
  for_each = local.rules

  direction         = "ingress"
  ethertype         = "IPv4"
  description       = each.value.description
  protocol          = each.value.protocol
  port_range_min    = each.value.port
  port_range_max    = each.value.port
  remote_ip_prefix  = each.value.source

  security_group_id = sbercloud_networking_secgroup.sg_01.id
}
