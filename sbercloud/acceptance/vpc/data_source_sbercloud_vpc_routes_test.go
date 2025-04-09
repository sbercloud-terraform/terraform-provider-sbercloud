package vpc

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVpcRoutes_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dataSource1 := "data.sbercloud_vpc_routes.basic"
	dataSource2 := "data.sbercloud_vpc_routes.filter_by_type"
	dataSource3 := "data.sbercloud_vpc_routes.filter_by_vpc_id"
	dataSource4 := "data.sbercloud_vpc_routes.filter_by_destination"
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)
	dc3 := acceptance.InitDataSourceCheck(dataSource3)
	dc4 := acceptance.InitDataSourceCheck(dataSource4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceVpcRoutes_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					dc3.CheckResourceExists(),
					dc4.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_vpc_id_filter_useful", "true"),
					resource.TestCheckOutput("is_destination_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceVpcRoutes_basic(name string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc_routes" "basic" {
  depends_on = [sbercloud_vpc_route.test]
}

data "sbercloud_vpc_routes" "filter_by_type" {
  type = "peering"

  depends_on = [sbercloud_vpc_route.test]
}

data "sbercloud_vpc_routes" "filter_by_vpc_id" {
  vpc_id = sbercloud_vpc.test1.id

  depends_on = [sbercloud_vpc_route.test]
}

data "sbercloud_vpc_routes" "filter_by_destination" {
  destination = sbercloud_vpc.test2.cidr

  depends_on = [sbercloud_vpc_route.test]
}

locals {
  type_filter_result = [for v in data.sbercloud_vpc_routes.filter_by_type.routes[*].type : v == "peering"]
  vpc_id_filter_result = [
	for v in data.sbercloud_vpc_routes.filter_by_vpc_id.routes[*].vpc_id : v == sbercloud_vpc.test1.id
  ]
  destination_filter_result = [
	for v in data.sbercloud_vpc_routes.filter_by_destination.routes[*].destination : v == sbercloud_vpc.test2.cidr
  ]
}

output "is_results_not_empty" {
  value = length(data.sbercloud_vpc_routes.basic.routes) > 0
}

output "is_type_filter_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}

output "is_vpc_id_filter_useful" {
  value = alltrue(local.vpc_id_filter_result) && length(local.vpc_id_filter_result) > 0
}

output "is_destination_filter_useful" {
  value = alltrue(local.destination_filter_result) && length(local.destination_filter_result) > 0
}
`, testAccRouteDataSource_base(name))
}
