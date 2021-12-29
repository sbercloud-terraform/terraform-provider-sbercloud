package iam

import (
	"testing"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIdentityRoleDataSource_basic(t *testing.T) {
	resourceName := "data.sbercloud_identity_role.role_1"
	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityRoleDataSource_by_name,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", "system_all_64"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "EPS Viewer"),
				),
			},
		},
	})
}

const testAccIdentityRoleDataSource_by_name = `
data "sbercloud_identity_role" "role_1" {
  # EPS Viewer
  name = "system_all_64"
}
`
