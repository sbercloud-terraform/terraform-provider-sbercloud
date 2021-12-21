package sbercloud

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/fwaas_v2/firewall_groups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccNetworkACL_basic(t *testing.T) {
	rName := fmt.Sprintf("acc-fw-%s", acctest.RandString(5))
	resourceKey := "sbercloud_network_acl.fw_1"
	var fwGroup huaweicloud.FirewallGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACL_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists(resourceKey, &fwGroup),
					resource.TestCheckResourceAttr(resourceKey, "name", rName),
					resource.TestCheckResourceAttr(resourceKey, "description", "created by terraform test acc"),
					resource.TestCheckResourceAttr(resourceKey, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(resourceKey, "inbound_policy_id"),
					testAccCheckFWFirewallPortCount(&fwGroup, 1),
				),
			},
			{
				Config: testAccNetworkACL_basic_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists("sbercloud_network_acl.fw_1", &fwGroup),
					resource.TestCheckResourceAttr(resourceKey, "name", rName+"_update"),
					resource.TestCheckResourceAttr(resourceKey, "description", "updated by terraform test acc"),
					resource.TestCheckResourceAttr(resourceKey, "status", "ACTIVE"),
					testAccCheckFWFirewallPortCount(&fwGroup, 2),
				),
			},
		},
	})
}

func TestAccNetworkACL_no_subnets(t *testing.T) {
	rName := fmt.Sprintf("acc-fw-%s", acctest.RandString(5))
	resourceKey := "sbercloud_network_acl.fw_1"
	var fwGroup huaweicloud.FirewallGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACL_no_subnets(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists(resourceKey, &fwGroup),
					resource.TestCheckResourceAttr(resourceKey, "name", rName),
					resource.TestCheckResourceAttr(resourceKey, "description", "network acl without subents"),
					resource.TestCheckResourceAttr(resourceKey, "status", "INACTIVE"),
					testAccCheckFWFirewallPortCount(&fwGroup, 0),
				),
			},
		},
	})
}

func TestAccNetworkACL_remove(t *testing.T) {
	rName := fmt.Sprintf("acc-fw-%s", acctest.RandString(5))
	resourceKey := "sbercloud_network_acl.fw_1"
	var fwGroup huaweicloud.FirewallGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACL_basic_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists("sbercloud_network_acl.fw_1", &fwGroup),
					resource.TestCheckResourceAttr(resourceKey, "status", "ACTIVE"),
					testAccCheckFWFirewallPortCount(&fwGroup, 2),
				),
			},
			{
				Config: testAccNetworkACL_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists("sbercloud_network_acl.fw_1", &fwGroup),
					resource.TestCheckResourceAttr(resourceKey, "status", "ACTIVE"),
					testAccCheckFWFirewallPortCount(&fwGroup, 1),
				),
			},
			{
				Config: testAccNetworkACL_no_subnets(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists(resourceKey, &fwGroup),
					resource.TestCheckResourceAttr(resourceKey, "status", "INACTIVE"),
					testAccCheckFWFirewallPortCount(&fwGroup, 0),
				),
			},
		},
	})
}

func testAccCheckNetworkACLDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	fwClient, err := config.FwV2Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud fw client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_network_acl" {
			continue
		}

		_, err = firewall_groups.Get(fwClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Firewall group (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckNetworkACLExists(n string, fwGroup *huaweicloud.FirewallGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set in %s", n)
		}

		config := testAccProvider.Meta().(*config.Config)
		fwClient, err := config.FwV2Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud fw client: %s", err)
		}

		var found huaweicloud.FirewallGroup
		err = firewall_groups.Get(fwClient, rs.Primary.ID).ExtractInto(&found)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Firewall group not found")
		}

		*fwGroup = found

		return nil
	}
}

func testAccCheckFWFirewallPortCount(firewall_group *huaweicloud.FirewallGroup, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(firewall_group.PortIDs) != expected {
			return fmt.Errorf("Expected %d Ports, got %d", expected, len(firewall_group.PortIDs))
		}

		return nil
	}
}

func testAccNetworkACLRules(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc" "vpc_1" {
  name = "%s_vpc"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_subnet" "subnet_1" {
  name = "%s_subnet_1"
  cidr = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id = sbercloud_vpc.vpc_1.id
}

resource "sbercloud_vpc_subnet" "subnet_2" {
	name = "%s_subnet_2"
	cidr = "192.168.10.0/24"
	gateway_ip = "192.168.10.1"
	vpc_id = sbercloud_vpc.vpc_1.id
  }

resource "sbercloud_network_acl_rule" "rule_1" {
  name             = "%s-rule-1"
  description      = "drop TELNET traffic"
  action           = "deny"
  protocol         = "tcp"
  destination_port = "23"
}

resource "sbercloud_network_acl_rule" "rule_2" {
  name             = "%s-rule-2"
  description      = "drop NTP traffic"
  action           = "deny"
  protocol         = "udp"
  destination_port = "123"
}
`, name, name, name, name, name)
}

func testAccNetworkACL_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_network_acl" "fw_1" {
  name        = "%s"
  description = "created by terraform test acc"

  inbound_rules = [sbercloud_network_acl_rule.rule_1.id]
  subnets = [sbercloud_vpc_subnet.subnet_1.id]
}
`, testAccNetworkACLRules(name), name)
}

func testAccNetworkACL_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_network_acl" "fw_1" {
  name        = "%s_update"
  description = "updated by terraform test acc"

  inbound_rules = [sbercloud_network_acl_rule.rule_1.id,
      sbercloud_network_acl_rule.rule_2.id]
  subnets = [sbercloud_vpc_subnet.subnet_1.id,
      sbercloud_vpc_subnet.subnet_2.id]
}
`, testAccNetworkACLRules(name), name)
}

func testAccNetworkACL_no_subnets(name string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_network_acl" "fw_1" {
  name        = "%s"
  description = "network acl without subents"
}
`, testAccNetworkACLRules(name), name)
}
