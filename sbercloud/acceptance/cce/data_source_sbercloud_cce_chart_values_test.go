package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccDataSourceShowChartValues_basic(t *testing.T) {
	dataSource := "data.sbercloud_cce_chart_values.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceChartPath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceShowChartValues_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "values.%"),
				),
			},
		},
	})
}

func testAccDataSourceShowChartValues_basic() string {
	return fmt.Sprintf(`
%s

data "sbercloud_cce_chart_values" "test" {
  chart_id = sbercloud_cce_chart.test.id
}
`, testAccChart_basic())
}
