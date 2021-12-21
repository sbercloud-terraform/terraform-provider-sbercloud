package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3/groups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccIdentityV3Group_basic(t *testing.T) {
	var group groups.Group
	var groupName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3GroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3Group_basic(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3GroupExists("sbercloud_identity_group.group_1", &group),
					resource.TestCheckResourceAttrPtr(
						"sbercloud_identity_group.group_1", "name", &group.Name),
					resource.TestCheckResourceAttrPtr(
						"sbercloud_identity_group.group_1", "description", &group.Description),
				),
			},
			{
				Config: testAccIdentityV3Group_update(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3GroupExists("sbercloud_identity_group.group_1", &group),
					resource.TestCheckResourceAttrPtr(
						"sbercloud_identity_group.group_1", "name", &group.Name),
					resource.TestCheckResourceAttrPtr(
						"sbercloud_identity_group.group_1", "description", &group.Description),
				),
			},
		},
	})
}

func testAccCheckIdentityV3GroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	identityClient, err := config.IdentityV3Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_identity_group" {
			continue
		}

		_, err := groups.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Group still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3GroupExists(n string, group *groups.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		identityClient, err := config.IdentityV3Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud identity client: %s", err)
		}

		found, err := groups.Get(identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Group not found")
		}

		*group = *found

		return nil
	}
}

func testAccIdentityV3Group_basic(groupName string) string {
	return fmt.Sprintf(`
    resource "sbercloud_identity_group" "group_1" {
      name = "%s"
      description = "A ACC test group"
    }
  `, groupName)
}

func testAccIdentityV3Group_update(groupName string) string {
	return fmt.Sprintf(`
    resource "sbercloud_identity_group" "group_1" {
      name = "%s"
      description = "Some Group"
    }
  `, groupName)
}
