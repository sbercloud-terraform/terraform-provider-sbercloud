package rds

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRdsPgPluginParameterValueRange_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.sbercloud_rds_pg_plugin_parameter_value_range.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePgPluginParameterValueRange_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "restart_required"),
					resource.TestCheckResourceAttrSet(rName, "values.#"),
					resource.TestCheckResourceAttrSet(rName, "default_values.#"),
				),
			},
		},
	})
}

func testAccDatasourcePgPluginParameterValueRange_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.sbercloud_availability_zones.test.names[0]]
  security_group_id = sbercloud_networking_secgroup.test.id
  subnet_id         = data.sbercloud_vpc_subnet.test.id
  vpc_id            = data.sbercloud_vpc.test.id

  db {
    type    = "PostgreSQL"
    version = "14"
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
}
`, testAccRdsInstance_base(name), name)
}

func testAccDatasourcePgPluginParameterValueRange_basic(name string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_rds_pg_plugin_parameter_value_range" "test" {
  instance_id = sbercloud_rds_instance.test.id
  name        = "shared_preload_libraries"
}
`, testAccDatasourcePgPluginParameterValueRange_base(name))
}
