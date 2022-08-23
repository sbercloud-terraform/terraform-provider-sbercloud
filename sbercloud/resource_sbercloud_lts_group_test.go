package sbercloud

import (
	"testing"

	"github.com/chnsz/golangsdk/openstack/lts/huawei/loggroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccLogTankGroupV2_basic(t *testing.T) {
	var group loggroups.LogGroup

	resourceName := "sbercloud_lts_group.testacc_group"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogTankGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLogTankGroupV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogTankGroupV2Exists(
						resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "group_name", "testacc_group"),
					resource.TestCheckResourceAttr(
						resourceName, "ttl_in_days", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLogTankGroupV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	ltsclient, err := config.LtsV2Client(SBC_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating SberCloud LTS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_lts_group" {
			continue
		}

		groups, err := loggroups.List(ltsclient).Extract()
		if err != nil {
			return fmtp.Errorf("Log group get list err: %s", err.Error())
		}
		for _, group := range groups.LogGroups {
			if group.ID == rs.Primary.ID {
				return fmtp.Errorf("Log group (%s) still exists.", rs.Primary.ID)
			}
		}

	}
	return nil
}

func testAccCheckLogTankGroupV2Exists(n string, group *loggroups.LogGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		ltsclient, err := config.LtsV2Client(SBC_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating SberCloud LTS client: %s", err)
		}

		var founds *loggroups.LogGroups
		founds, err = loggroups.List(ltsclient).Extract()
		if err != nil {
			return err
		}
		for _, loggroup := range founds.LogGroups {
			if rs.Primary.ID == loggroup.ID {
				*group = loggroup
				return nil
			}
		}

		return fmtp.Errorf("Error SberCloud log group %s: No Found", rs.Primary.ID)
	}
}

const testAccLogTankGroupV2_basic = `
resource "sbercloud_lts_group" "testacc_group" {
	group_name  = "testacc_group"
	ttl_in_days = 1
}
`
