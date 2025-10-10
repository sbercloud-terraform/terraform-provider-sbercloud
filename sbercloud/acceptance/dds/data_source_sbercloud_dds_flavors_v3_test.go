package dds

import (
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDDSFlavorV3DataSource_basic(t *testing.T) {
	dataSourceName := "data.sbercloud_dds_flavors.flavor"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSFlavorV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.engine_name", "DDS-Community"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.engine_versions.0", "4.0"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.vcpus", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.memory", "4"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.type", "mongos"),

					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.az_status.%"),
				),
			},
		},
	})
}

var testAccDDSFlavorV3DataSource_basic = `
data "sbercloud_dds_flavors" "flavor" {
  engine_name    = "DDS-Community"
  engine_version = "4.0"
  vcpus          = 2
  memory         = 4
  type           = "mongos"
}
`
