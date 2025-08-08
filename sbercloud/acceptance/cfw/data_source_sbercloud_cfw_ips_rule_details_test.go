package cfw

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCfwIpsRuleDetails_basic(t *testing.T) {
	dataSource := "data.sbercloud_cfw_ips_rule_details.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCfwIpsRuleDetails_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.ips_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.ips_version"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.update_time"),
				),
			},
		},
	})
}

func testDataSourceCfwIpsRuleDetails_basic() string {
	return fmt.Sprintf(`
data "sbercloud_cfw_ips_rule_details" "test" {
  fw_instance_id = "%[1]s"
}
`, acceptance.SBC_CFW_INSTANCE_ID)
}
