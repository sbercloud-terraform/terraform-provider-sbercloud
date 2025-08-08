package cfw

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCfwRegions_basic(t *testing.T) {
	dataSource := "data.sbercloud_cfw_regions.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCfwRegions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data"),
				),
			},
		},
	})
}

func testDataSourceCfwRegions_basic() string {
	return fmt.Sprintf(`
data "sbercloud_cfw_regions" "test" {
  fw_instance_id = "%[1]s"
}
`, acceptance.SBC_CFW_INSTANCE_ID)
}
