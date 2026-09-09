package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccChartsDataSource_basic(t *testing.T) {
	datasourceName := "data.sbercloud_cce_charts.test"
	dc := acceptance.InitDataSourceCheck(datasourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceChartPath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccChartsDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
				),
			},
		},
	})
}

func testAccChartsDataSource_basic() string {
	return fmt.Sprintf(`
%s

data "sbercloud_cce_charts" "test" {
  depends_on = [sbercloud_cce_chart.test]
}

output "is_results_not_empty" {
  value = length(data.sbercloud_cce_charts.test.charts) > 0
}
`, testAccChart_basic())
}
