package vpn

import (
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourceVpnGatewayAZs_basic(t *testing.T) {
	rName := "data.sbercloud_vpn_gateway_availability_zones.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceVpnGatewayAZs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "names.#"),
				),
			},
		},
	})
}

func testAccDatasourceVpnGatewayAZs_basic() string {
	return `
data "sbercloud_vpn_gateway_availability_zones" "test" {
  flavor = "Basic"
}`
}
