package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVpcRouteV2DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRouteV2Config(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteV2DataSourceID("data.sbercloud_vpc_route.by_id"),
					resource.TestCheckResourceAttr(
						"data.sbercloud_vpc_route.by_id", "type", "peering"),
					resource.TestCheckResourceAttr(
						"data.sbercloud_vpc_route.by_id", "destination", "192.168.0.0/16"),
					testAccCheckRouteV2DataSourceID("data.sbercloud_vpc_route.by_vpc_id"),
					resource.TestCheckResourceAttr(
						"data.sbercloud_vpc_route.by_vpc_id", "type", "peering"),
					resource.TestCheckResourceAttr(
						"data.sbercloud_vpc_route.by_vpc_id", "destination", "192.168.0.0/16"),
				),
			},
		},
	})
}

func testAccCheckRouteV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find vpc route connection data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("vpc route connection data source ID not set")
		}

		return nil
	}
}

func testAccDataSourceRouteV2Config(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc" "test2" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_peering_connection" "test" {
  name        = "%s"
  vpc_id      = sbercloud_vpc.test.id
  peer_vpc_id = sbercloud_vpc.test2.id
}

resource "sbercloud_vpc_route" "test" {
  type        = "peering"
  nexthop     = sbercloud_vpc_peering_connection.test.id
  destination = "192.168.0.0/16"
  vpc_id      = sbercloud_vpc.test.id
}

data "sbercloud_vpc_route" "by_id" {
  id = sbercloud_vpc_route.test.id
}

data "sbercloud_vpc_route" "by_vpc_id" {
  vpc_id = sbercloud_vpc_route.test.vpc_id
}
`, rName, rName+"2", rName)
}
