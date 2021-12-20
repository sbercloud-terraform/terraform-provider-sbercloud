package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v1/eips"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccNetworkingV2EIPAssociate_basic(t *testing.T) {
	var eip eips.PublicIp

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_networking_eip_associate.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2EIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2EIPAssociate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1EIPExists("sbercloud_vpc_eip.test", &eip),
					resource.TestCheckResourceAttrPtr(
						resourceName, "public_ip", &eip.PublicAddress),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2EIPAssociateDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV1Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating EIP Client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_vpc_eip" {
			continue
		}

		_, err := eips.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("EIP still exists")
		}
	}

	return nil
}

func testAccNetworkingV2EIPAssociate_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%s"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "sbercloud_networking_eip_associate" "test" {
  public_ip   = sbercloud_vpc_eip.test.address
  port_id     = sbercloud_compute_instance.test.network[0].port
}
`, testAccComputeV2Instance_basic(rName), rName)
}
