package rds

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRdsPgPluginParameterValues_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.sbercloud_rds_pg_plugin_parameter_values.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePgPluginParameterValues_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "values.0",
						"data.sbercloud_rds_pg_plugin_parameter_value_range.test", "values.0"),
				),
			},
		},
	})
}

func testAccDatasourcePgPluginParameterValues_basic(name string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_rds_pg_plugin_parameter_values" "test" {
  depends_on = [sbercloud_rds_pg_plugin_parameter.test]

  instance_id = sbercloud_rds_instance.test.id
  name        = "shared_preload_libraries"
}
`, testPgPluginParameter_basic(name))
}
