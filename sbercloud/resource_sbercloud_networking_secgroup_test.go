package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/security/groups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccNetworkingV2SecGroup_basic(t *testing.T) {
	var security_group groups.SecGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroup_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(
						"sbercloud_networking_secgroup.secgroup_1", &security_group),
					testAccCheckNetworkingV2SecGroupRuleCount(&security_group, 2),
				),
			},
			{
				ResourceName:      "sbercloud_networking_secgroup.secgroup_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkingV2SecGroup_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(
						"sbercloud_networking_secgroup.secgroup_1", "id", &security_group.ID),
					resource.TestCheckResourceAttr(
						"sbercloud_networking_secgroup.secgroup_1", "name", "security_group_2"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroup_noDefaultRules(t *testing.T) {
	var security_group groups.SecGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroup_noDefaultRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(
						"sbercloud_networking_secgroup.secgroup_1", &security_group),
					testAccCheckNetworkingV2SecGroupRuleCount(&security_group, 0),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroup_timeout(t *testing.T) {
	var security_group groups.SecGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroup_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(
						"sbercloud_networking_secgroup.secgroup_1", &security_group),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2SecGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV2Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_networking_secgroup" {
			continue
		}

		_, err := groups.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Security group still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2SecGroupExists(n string, security_group *groups.SecGroup) resource.TestCheckFunc {
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

		found, err := groups.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Security group not found")
		}

		*security_group = *found

		return nil
	}
}

func testAccCheckNetworkingV2SecGroupRuleCount(
	sg *groups.SecGroup, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(sg.Rules) == count {
			return nil
		}

		return fmt.Errorf("Unexpected number of rules in group %s. Expected %d, got %d",
			sg.ID, count, len(sg.Rules))
	}
}

const testAccNetworkingV2SecGroup_basic = `
resource "sbercloud_networking_secgroup" "secgroup_1" {
  name = "security_group"
  description = "terraform security group acceptance test"
}
`

const testAccNetworkingV2SecGroup_update = `
resource "sbercloud_networking_secgroup" "secgroup_1" {
  name = "security_group_2"
  description = "terraform security group acceptance test"
}
`

const testAccNetworkingV2SecGroup_noDefaultRules = `
resource "sbercloud_networking_secgroup" "secgroup_1" {
	name = "security_group_1"
	description = "terraform security group acceptance test"
	delete_default_rules = true
}
`

const testAccNetworkingV2SecGroup_timeout = `
resource "sbercloud_networking_secgroup" "secgroup_1" {
  name = "security_group"
  description = "terraform security group acceptance test"

  timeouts {
    delete = "5m"
  }
}
`
