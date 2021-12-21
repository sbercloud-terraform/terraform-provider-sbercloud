package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingV2PortDataSource_basic(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.sbercloud_networking_port.port_3", "all_fixed_ips.#", "1"),
				),
			},
		},
	})
}

func testAccNetworkingV2PortDataSource_basic() string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "mynet" {
  name = "subnet-default"
}

data "sbercloud_networking_port" "port_3" {
  network_id = data.sbercloud_vpc_subnet.mynet.id
  fixed_ip = data.sbercloud_vpc_subnet.mynet.gateway_ip
}
`)
}
