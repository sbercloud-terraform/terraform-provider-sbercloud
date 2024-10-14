package cbh

import (
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourceCbhFlavors_basic(t *testing.T) {
	rName := "data.sbercloud_cbh_flavors.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceCbhFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.id"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.ecs_system_data_size"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.vcpus"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.memory"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.asset"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.max_connection"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.type"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.data_disk_size"),

					resource.TestCheckOutput("action_is_useful", "true"),
					resource.TestCheckOutput("flavor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("asset_filter_is_useful", "true"),
					resource.TestCheckOutput("memory_filter_is_useful", "true"),
					resource.TestCheckOutput("vcpus_filter_is_useful", "true"),
					resource.TestCheckOutput("max_connection_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceCbhFlavors_basic() string {
	return `
data "sbercloud_cbh_flavors" "test" {}

locals {
  flavor_id = data.sbercloud_cbh_flavors.test.flavors[0].id
}

data "sbercloud_cbh_flavors" "action_filter" {
  action    = "update"
  spec_code = local.flavor_id
}

output "action_is_useful" {
  value = length(data.sbercloud_cbh_flavors.action_filter.flavors) > 0
}

data "sbercloud_cbh_flavors" "flavor_id_filter" {
  flavor_id = local.flavor_id
}
output "flavor_id_filter_is_useful" {
  value = length(data.sbercloud_cbh_flavors.flavor_id_filter.flavors) > 0 && alltrue(
    [for v in data.sbercloud_cbh_flavors.flavor_id_filter.flavors[*].id : v == local.flavor_id]
  )  
}

locals {
  type = data.sbercloud_cbh_flavors.test.flavors[0].type
}
data "sbercloud_cbh_flavors" "type_filter" {
  type = local.type
}
output "type_filter_is_useful" {
  value = length(data.sbercloud_cbh_flavors.type_filter.flavors) > 0 && alltrue(
    [for v in data.sbercloud_cbh_flavors.type_filter.flavors[*].type : v == local.type]
  )  
}

locals {
  asset = data.sbercloud_cbh_flavors.test.flavors[0].asset
}
data "sbercloud_cbh_flavors" "asset_filter" {
  asset = local.asset
}
output "asset_filter_is_useful" {
  value = length(data.sbercloud_cbh_flavors.asset_filter.flavors) > 0 && alltrue(
    [for v in data.sbercloud_cbh_flavors.asset_filter.flavors[*].asset : v == local.asset]
  )  
}

locals {
  memory = data.sbercloud_cbh_flavors.test.flavors[0].memory
}
data "sbercloud_cbh_flavors" "memory_filter" {
  memory = local.memory
}
output "memory_filter_is_useful" {
  value = length(data.sbercloud_cbh_flavors.memory_filter.flavors) > 0 && alltrue(
    [for v in data.sbercloud_cbh_flavors.memory_filter.flavors[*].memory : v == local.memory]
  )  
}

locals {
  vcpus = data.sbercloud_cbh_flavors.test.flavors[0].vcpus
}
data "sbercloud_cbh_flavors" "vcpus_filter" {
  vcpus = local.vcpus
}
output "vcpus_filter_is_useful" {
  value = length(data.sbercloud_cbh_flavors.vcpus_filter.flavors) > 0 && alltrue(
    [for v in data.sbercloud_cbh_flavors.vcpus_filter.flavors[*].vcpus : v == local.vcpus]
  )  
}

locals {
  max_connection = data.sbercloud_cbh_flavors.test.flavors[0].max_connection
}
data "sbercloud_cbh_flavors" "max_connection_filter" {
  max_connection = local.max_connection
}
output "max_connection_filter_is_useful" {
  value = length(data.sbercloud_cbh_flavors.max_connection_filter.flavors) > 0 && alltrue(
    [for v in data.sbercloud_cbh_flavors.max_connection_filter.flavors[*].max_connection : v == local.max_connection]
  )  
}
`
}
