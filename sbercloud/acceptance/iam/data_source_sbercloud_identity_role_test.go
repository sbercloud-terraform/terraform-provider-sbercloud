package iam

import (
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIdentityRoleDataSource_basic(t *testing.T) {
	dataSourceName := "data.sbercloud_identity_role.role_1"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

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
					resource.TestCheckResourceAttr(dataSourceName, "name", "system_all_64"),
					resource.TestCheckResourceAttr(dataSourceName, "display_name", "EPS ReadOnlyAccess"),
				),
			},
			{
				Config: testAccIdentityRoleDataSource_by_displayname,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", "kms_adm"),
					resource.TestCheckResourceAttr(dataSourceName, "display_name", "KMS Administrator"),
				),
			},
		},
	})
}

const testAccIdentityRoleDataSource_by_name = `
data "sbercloud_identity_role" "role_1" {
  # OBS ReadOnlyAccess
  name = "system_all_64"
}
`

const testAccIdentityRoleDataSource_by_displayname = `
data "sbercloud_identity_role" "role_1" {
  display_name = "KMS Administrator"
}
`
