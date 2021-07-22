package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/vpnaas/services"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccVpnServiceV2_basic(t *testing.T) {
	var service services.Service
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpnServiceV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnServiceV2_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnServiceV2Exists(
						"sbercloud_vpnaas_service.service_1", &service),
					resource.TestCheckResourceAttrPtr("sbercloud_vpnaas_service.service_1", "router_id", &service.RouterID),
					resource.TestCheckResourceAttr("sbercloud_vpnaas_service.service_1", "admin_state_up", "true"),
				),
			},
		},
	})
}

func testAccCheckVpnServiceV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV2Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_vpnaas_service" {
			continue
		}
		_, err = services.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Service (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckVpnServiceV2Exists(n string, serv *services.Service) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		networkingClient, err := config.NetworkingV2Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud networking client: %s", err)
		}

		var found *services.Service

		found, err = services.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*serv = *found

		return nil
	}
}

func testAccVpnServiceV2_basic(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpnaas_service" "service_1" {
  name = "%s"
  router_id = sbercloud_vpc.test.id
}
`, name, name)
}
