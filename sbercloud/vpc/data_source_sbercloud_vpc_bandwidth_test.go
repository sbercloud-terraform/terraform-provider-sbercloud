package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccBandWidthDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.sbercloud_vpc_bandwidth.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccBandWidthDataSource_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "size", "10"),
				),
			},
		},
	})
}

func testAccBandWidthDataSource_basic(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc_bandwidth" "test" {
	name = "%s"
	size = 10
}

data "sbercloud_vpc_bandwidth" "test" {
  name = sbercloud_vpc_bandwidth.test.name
}
`, rName)
}
