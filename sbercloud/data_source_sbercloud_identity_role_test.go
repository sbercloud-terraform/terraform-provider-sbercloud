package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccIdentityV3RoleDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3RoleDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3DataSourceID("data.sbercloud_identity_role.role_1"),
					resource.TestCheckResourceAttr(
						"data.sbercloud_identity_role.role_1", "name", "secu_admin"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find role data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Role data source ID not set")
		}

		return nil
	}
}

const testAccIdentityV3RoleDataSource_basic = `
data "sbercloud_identity_role" "role_1" {
    name = "secu_admin"
}
`
