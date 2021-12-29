package iam

import (
	"fmt"
	"testing"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIdentityGroupDataSource_basic(t *testing.T) {
	resourceName := "data.sbercloud_identity_group.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityGroupDataSource_by_name(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func testAccIdentityGroupDataSource_by_name(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_identity_group" "group_1" {
  name        = "%s"
  description = "A ACC test group"
}

data "sbercloud_identity_group" "test" {
  name = sbercloud_identity_group.group_1.name
  
  depends_on = [
    sbercloud_identity_group.group_1
  ]
}
`, rName)
}
