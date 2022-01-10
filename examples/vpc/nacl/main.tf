# Define rules content for NACL rules in local variables. You can pass it via terraform.tfvars, too.
locals {
  inbound_rules = {
    rule_http = {
      name                   = "rule-http",
      description            = "Allow HTTP from anywhere",
      protocol               = "tcp",
      action                 = "allow",
      source_ip_address      = "0.0.0.0/0",
      destination_ip_address = "172.16.10.0/24",
      destination_port       = "80"
    },
    rule-ssh = {
      name                   = "rule-ssh",
      description            = "Allow SSH from 172.16.20.23 to 172.16.10.16",
      protocol               = "tcp",
      action                 = "allow",
      source_ip_address      = "172.16.20.23/32",
      destination_ip_address = "172.16.10.16/32",
      destination_port       = "22"
    }
  }
  outbound_rules = {
    rule_all = {
      name                   = "rule-all",
      description            = "Allow all from 172.16.10.100",
      protocol               = "any",
      action                 = "allow",
      source_ip_address      = "172.16.10.100/32",
      destination_ip_address = "0.0.0.0/0"
    }
  }
}

# Get the subnet which NACL will be associated with
data "sbercloud_vpc_subnet" "subnet_01" {
  name = "put_here_name_of_your_existing_subnet"
}

# Create inbound NACL rules
resource "sbercloud_network_acl_rule" "inbound_rules" {
  for_each = local.inbound_rules

  name                   = each.value.name
  description            = each.value.description
  protocol               = each.value.protocol
  action                 = each.value.action
  source_ip_address      = each.value.source_ip_address
  destination_ip_address = each.value.destination_ip_address
  destination_port       = each.value.destination_port
}

# Create outbound NACL rules
resource "sbercloud_network_acl_rule" "outbound_rules" {
  for_each = local.outbound_rules

  name                   = each.value.name
  description            = each.value.description
  protocol               = each.value.protocol
  action                 = each.value.action
  source_ip_address      = each.value.source_ip_address
  destination_ip_address = each.value.destination_ip_address
}

# Create NACL and associate it with subnet
resource "sbercloud_network_acl" "nacl_01" {
  name           = "nacl-tf"
  subnets        = [data.sbercloud_vpc_subnet.subnet_01.id]
  inbound_rules  = [for rule in sbercloud_network_acl_rule.inbound_rules: rule.id]
  outbound_rules = [for rule in sbercloud_network_acl_rule.outbound_rules: rule.id]
}
