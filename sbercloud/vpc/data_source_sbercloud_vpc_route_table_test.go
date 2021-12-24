package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccVpcRouteTableDataSource_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dataSourceName := "data.sbercloud_vpc_route_table.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRouteTable_base(rName),
			},
			{
				Config: testAccDataSourceRouteTable_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "default", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceRouteTable_base(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc" "test" {
  name = "%s"
  cidr = "172.16.0.0/16"
}

resource "sbercloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "172.16.10.0/24"
  gateway_ip = "172.16.10.1"
  vpc_id     = sbercloud_vpc.test.id
}
`, rName, rName)
}

func testAccDataSourceRouteTable_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc_route_table" "test" {
  vpc_id = sbercloud_vpc.test.id
}
`, testAccDataSourceRouteTable_base(rName))
}
