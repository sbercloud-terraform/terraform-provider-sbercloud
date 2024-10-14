package cbh

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourceCbhInstances_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSourceName := "data.sbercloud_cbh_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceCbhInstances_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.private_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.security_group_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.flavor_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.version"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),

					resource.TestCheckOutput("vpc_id_filter_is_useful", "true"),

					resource.TestCheckOutput("subnet_id_filter_is_useful", "true"),

					resource.TestCheckOutput("security_group_id_filter_is_useful", "true"),

					resource.TestCheckOutput("flavor_id_filter_is_useful", "true"),

					resource.TestCheckOutput("version_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceCbhInstances_basic(name string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_cbh_instances" "test" {
  depends_on = [sbercloud_cbh_instance.test]
}

locals {
  name = data.sbercloud_cbh_instances.test.instances[0].name
}

data "sbercloud_cbh_instances" "filter_by_name" {
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.sbercloud_cbh_instances.filter_by_name.instances[*].name : 
    v == local.name
  ]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

locals {
  vpc_id = data.sbercloud_cbh_instances.test.instances[0].vpc_id
}

data "sbercloud_cbh_instances" "filter_by_vpc_id" {
  vpc_id = local.vpc_id
}

locals {
  vpc_id_filter_result = [
    for v in data.sbercloud_cbh_instances.filter_by_vpc_id.instances[*].vpc_id : 
    v == local.vpc_id
  ]
}

output "vpc_id_filter_is_useful" {
  value = alltrue(local.vpc_id_filter_result) && length(local.vpc_id_filter_result) > 0
}

locals {
  subnet_id = data.sbercloud_cbh_instances.test.instances[0].subnet_id
}

data "sbercloud_cbh_instances" "filter_by_subnet_id" {
  subnet_id = local.subnet_id
}

locals {
  subnet_id_filter_result = [
    for v in data.sbercloud_cbh_instances.filter_by_subnet_id.instances[*].subnet_id : 
    v == local.subnet_id
  ]
}

output "subnet_id_filter_is_useful" {
  value = alltrue(local.subnet_id_filter_result) && length(local.subnet_id_filter_result) > 0
}

locals {
  security_group_id = data.sbercloud_cbh_instances.test.instances[0].security_group_id
}

data "sbercloud_cbh_instances" "filter_by_security_group_id" {
  security_group_id = local.security_group_id
}

locals {
  security_group_id_filter_result = [
    for v in data.sbercloud_cbh_instances.filter_by_security_group_id.instances[*].security_group_id : 
    v == local.security_group_id
  ]
}

output "security_group_id_filter_is_useful" {
  value = alltrue(local.security_group_id_filter_result) && length(local.security_group_id_filter_result) > 0
}

locals {
  flavor_id = data.sbercloud_cbh_instances.test.instances[0].flavor_id
}

data "sbercloud_cbh_instances" "filter_by_flavor_id" {
  flavor_id = local.flavor_id
}

locals {
  flavor_id_filter_result = [
    for v in data.sbercloud_cbh_instances.filter_by_flavor_id.instances[*].flavor_id : 
    v == local.flavor_id
  ]
}

output "flavor_id_filter_is_useful" {
  value = alltrue(local.flavor_id_filter_result) && length(local.flavor_id_filter_result) > 0
}

locals {
  version = data.sbercloud_cbh_instances.test.instances[0].version
}

data "sbercloud_cbh_instances" "filter_by_version" {
  version = local.version
}

locals {
  version_filter_result = [
    for v in data.sbercloud_cbh_instances.filter_by_version.instances[*].version : 
    v == local.version
  ]
}

output "version_filter_is_useful" {
  value = alltrue(local.version_filter_result) && length(local.version_filter_result) > 0
}
`, testCBHInstance_basic(name))
}
