package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccComputeServerGroupsDataSource_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	dataSourceName := "data.sbercloud_compute_servergroups.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroupsDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "servergroups.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "servergroups.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "servergroups.0.name"),
				),
			},
		},
	})
}

func testAccComputeServerGroupsDataSource_basic(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_compute_servergroup" "test" {
  name     = "%s"
  policies = ["anti-affinity"]
}

data "sbercloud_compute_servergroups" "test" {
  name = sbercloud_compute_servergroup.test.name
}
`, rName)
}
