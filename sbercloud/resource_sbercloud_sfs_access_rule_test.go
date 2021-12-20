package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/sfs/v2/shares"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccSFSAccessRuleV2_basic(t *testing.T) {
	var rule shares.AccessRight
	shareName := fmt.Sprintf("sfs-acc-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSAccessRuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: configAccSFSAccessRuleV2_basic(shareName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSAccessRuleV2Exists("sbercloud_sfs_access_rule.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"sbercloud_sfs_access_rule.rule_1", "access_level", "rw"),
					resource.TestCheckResourceAttr(
						"sbercloud_sfs_access_rule.rule_1", "status", "active"),
				),
			},
			{
				Config: configAccSFSAccessRuleV2_ipAuth(shareName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSAccessRuleV2Exists("sbercloud_sfs_access_rule.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"sbercloud_sfs_access_rule.rule_1", "status", "active"),
				),
			},
		},
	})
}

func testAccCheckSFSAccessRuleV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	sfsClient, err := config.SfsV2Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud sfs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_sfs_access_rule" {
			continue
		}

		sfsID := rs.Primary.Attributes["sfs_id"]
		if sfsID == "" {
			return fmt.Errorf("No SFSID is set in sbercloud_sfs_access_rule")
		}
		rules, err := shares.ListAccessRights(sfsClient, sfsID).ExtractAccessRights()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}

			return err
		}

		for _, v := range rules {
			if v.ID == rs.Primary.ID {
				return fmt.Errorf("resource sbercloud_sfs_access_rule still exists")
			}
		}
	}

	return nil
}

func testAccCheckSFSAccessRuleV2Exists(n string, rule *shares.AccessRight) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set in %s", n)
		}

		config := testAccProvider.Meta().(*config.Config)
		sfsClient, err := config.SfsV2Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud sfs client: %s", err)
		}

		sfsID := rs.Primary.Attributes["sfs_id"]
		if sfsID == "" {
			return fmt.Errorf("No SFSID is set in %s", n)
		}

		rules, err := shares.ListAccessRights(sfsClient, sfsID).ExtractAccessRights()
		if err != nil {
			return err
		}

		for _, v := range rules {
			if v.ID == rs.Primary.ID {
				*rule = v
				return nil
			}
		}

		return fmt.Errorf("sfs access rule %s was not found", n)
	}
}

func configAccSFSAccessRuleV2_basic(sfsName string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc" "vpc_default" {
  name = "vpc-default"
  enterprise_project_id = "0"
}

resource "sbercloud_sfs_file_system" "sfs_1" {
  share_proto = "NFS"
  size        = 10
  name        = "%s"
  description = "sfs file system created by terraform testacc"
}

resource "sbercloud_sfs_access_rule" "rule_1" {
  sfs_id = sbercloud_sfs_file_system.sfs_1.id
  access_to = data.sbercloud_vpc.vpc_default.id
}`, sfsName)
}

func configAccSFSAccessRuleV2_ipAuth(sfsName string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc" "vpc_default" {
  name = "vpc-default"
  enterprise_project_id = "0"
}

resource "sbercloud_sfs_file_system" "sfs_1" {
  share_proto = "NFS"
  size        = 10
  name        = "%s"
  description = "sfs file system created by terraform testacc"
}

resource "sbercloud_sfs_access_rule" "rule_1" {
  sfs_id = sbercloud_sfs_file_system.sfs_1.id
  access_to = join("#", [data.sbercloud_vpc.vpc_default.id, "192.168.10.0/24", "0", "no_all_squash,no_root_squash"])
}`, sfsName)
}
