package vpn

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVPNCustomerGatewaysDataSource_Basic(t *testing.T) {
	resourceName := "data.sbercloud_vpn_customer_gateways.services"
	dc := acceptance.InitDataSourceCheck(resourceName)
	rName := acceptance.RandomAccResourceName()
	ipAddress := "172.16.1.2"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNCustomerGatewaysDataSourceBasic(rName, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "customer_gateways.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "customer_gateways.0.asn", "65000"),
					resource.TestCheckResourceAttr(resourceName, "customer_gateways.0.ip", ipAddress),
					resource.TestCheckResourceAttr(resourceName, "customer_gateways.0.route_mode", "bgp"),
					resource.TestCheckResourceAttr(resourceName, "customer_gateways.0.id_type", "ip"),
					resource.TestCheckResourceAttr(resourceName, "customer_gateways.0.id_value", ipAddress),
					resource.TestCheckResourceAttrSet(resourceName, "customer_gateways.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_gateways.0.created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_gateways.0.updated_at"),
					resource.TestCheckResourceAttr(resourceName, "customer_gateways.0.ca_certificate.#", "1"),
				),
			},
		},
	})
}

func testAccVPNCustomerGatewaysDataSourceBasic(rName string, ipAddress string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpn_customer_gateway" "test" {
  name = "%s"
  ip   = "%s"
}

data "sbercloud_vpn_customer_gateways" "services" {
  asn                 = 65000
  route_mode          = "bgp"
  customer_gateway_id = sbercloud_vpn_customer_gateway.test.id
  name                = sbercloud_vpn_customer_gateway.test.name
  ip                  = sbercloud_vpn_customer_gateway.test.ip
}`, rName, ipAddress)
}
