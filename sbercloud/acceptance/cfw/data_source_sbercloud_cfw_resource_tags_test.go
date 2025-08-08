package cfw

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCfwResourceTags_basic(t *testing.T) {
	dataSource := "data.sbercloud_cfw_resource_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCfwResourceTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.value"),
				),
			},
		},
	})
}

func testDataSourceCfwResourceTags_basic() string {
	return fmt.Sprintf(`
data "sbercloud_cfw_resource_tags" "test" {
  fw_instance_id = "%[1]s"
}
`, acceptance.SBC_CFW_INSTANCE_ID)
}
