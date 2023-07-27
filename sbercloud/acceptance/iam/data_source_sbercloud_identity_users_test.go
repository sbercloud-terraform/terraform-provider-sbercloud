package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccIdentityUsersDataSource_basic(t *testing.T) {
	dataSourceName := "data.sbercloud_identity_users.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUsersDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.name"),
				),
			},
		},
	})
}

const testAccUsersDataSourceBasic string = `
data "sbercloud_identity_users" "test" {
}
`
