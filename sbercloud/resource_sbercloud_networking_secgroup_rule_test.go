package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/security/groups"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/security/rules"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccNetworkingV2SecGroupRule_basic(t *testing.T) {
	var secgroup_1 groups.SecGroup
	var secgroup_2 groups.SecGroup
	var secgroup_rule_1 rules.SecGroupRule
	var secgroup_rule_2 rules.SecGroupRule

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(
						"sbercloud_networking_secgroup.secgroup_1", &secgroup_1),
					testAccCheckNetworkingV2SecGroupExists(
						"sbercloud_networking_secgroup.secgroup_2", &secgroup_2),
					testAccCheckNetworkingV2SecGroupRuleExists(
						"sbercloud_networking_secgroup_rule.secgroup_rule_1", &secgroup_rule_1),
					testAccCheckNetworkingV2SecGroupRuleExists(
						"sbercloud_networking_secgroup_rule.secgroup_rule_2", &secgroup_rule_2),
				),
			},
			{
				ResourceName:      "sbercloud_networking_secgroup_rule.secgroup_rule_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkingV2SecGroupRule_lowerCaseCIDR(t *testing.T) {
	var secgroup_1 groups.SecGroup
	var secgroup_rule_1 rules.SecGroupRule

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRule_lowerCaseCIDR,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(
						"sbercloud_networking_secgroup.secgroup_1", &secgroup_1),
					testAccCheckNetworkingV2SecGroupRuleExists(
						"sbercloud_networking_secgroup_rule.secgroup_rule_1", &secgroup_rule_1),
					resource.TestCheckResourceAttr(
						"sbercloud_networking_secgroup_rule.secgroup_rule_1", "remote_ip_prefix", "2001:558:fc00::/39"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroupRule_timeout(t *testing.T) {
	var secgroup_1 groups.SecGroup
	var secgroup_2 groups.SecGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRule_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(
						"sbercloud_networking_secgroup.secgroup_1", &secgroup_1),
					testAccCheckNetworkingV2SecGroupExists(
						"sbercloud_networking_secgroup.secgroup_2", &secgroup_2),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroupRule_numericProtocol(t *testing.T) {
	var secgroup_1 groups.SecGroup
	var secgroup_rule_1 rules.SecGroupRule

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRule_numericProtocol,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(
						"sbercloud_networking_secgroup.secgroup_1", &secgroup_1),
					testAccCheckNetworkingV2SecGroupRuleExists(
						"sbercloud_networking_secgroup_rule.secgroup_rule_1", &secgroup_rule_1),
					resource.TestCheckResourceAttr(
						"sbercloud_networking_secgroup_rule.secgroup_rule_1", "protocol", "115"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2SecGroupRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV2Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_networking_secgroup_rule" {
			continue
		}

		_, err := rules.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Security group rule still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2SecGroupRuleExists(n string, security_group_rule *rules.SecGroupRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		networkingClient, err := config.NetworkingV2Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}

		found, err := rules.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Security group rule not found")
		}

		*security_group_rule = *found

		return nil
	}
}

const testAccNetworkingV2SecGroupRule_basic = `
resource "sbercloud_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "sbercloud_networking_secgroup" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group rule acceptance test"
}

resource "sbercloud_networking_secgroup_rule" "secgroup_rule_1" {
  direction = "ingress"
  ethertype = "IPv4"
  port_range_max = 22
  port_range_min = 22
  protocol = "tcp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = "${sbercloud_networking_secgroup.secgroup_1.id}"
}

resource "sbercloud_networking_secgroup_rule" "secgroup_rule_2" {
  direction = "ingress"
  ethertype = "IPv4"
  port_range_max = 80
  port_range_min = 80
  protocol = "tcp"
  remote_group_id = "${sbercloud_networking_secgroup.secgroup_1.id}"
  security_group_id = "${sbercloud_networking_secgroup.secgroup_2.id}"
}
`

const testAccNetworkingV2SecGroupRule_lowerCaseCIDR = `
resource "sbercloud_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "sbercloud_networking_secgroup_rule" "secgroup_rule_1" {
  direction = "ingress"
  ethertype = "IPv6"
  port_range_max = 22
  port_range_min = 22
  protocol = "tcp"
  remote_ip_prefix = "2001:558:FC00::/39"
  security_group_id = "${sbercloud_networking_secgroup.secgroup_1.id}"
}
`

const testAccNetworkingV2SecGroupRule_timeout = `
resource "sbercloud_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "sbercloud_networking_secgroup" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group rule acceptance test"
}

resource "sbercloud_networking_secgroup_rule" "secgroup_rule_1" {
  direction = "ingress"
  ethertype = "IPv4"
  port_range_max = 22
  port_range_min = 22
  protocol = "tcp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = "${sbercloud_networking_secgroup.secgroup_1.id}"

  timeouts {
    delete = "5m"
  }
}

resource "sbercloud_networking_secgroup_rule" "secgroup_rule_2" {
  direction = "ingress"
  ethertype = "IPv4"
  port_range_max = 80
  port_range_min = 80
  protocol = "tcp"
  remote_group_id = "${sbercloud_networking_secgroup.secgroup_1.id}"
  security_group_id = "${sbercloud_networking_secgroup.secgroup_2.id}"

  timeouts {
    delete = "5m"
  }
}
`

const testAccNetworkingV2SecGroupRule_numericProtocol = `
resource "sbercloud_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "sbercloud_networking_secgroup_rule" "secgroup_rule_1" {
  direction = "ingress"
  ethertype = "IPv4"
  port_range_max = 22
  port_range_min = 22
  protocol = "115"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = "${sbercloud_networking_secgroup.secgroup_1.id}"
}
`
