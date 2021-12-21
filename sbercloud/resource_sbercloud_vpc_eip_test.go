package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccVpcV1EIP_basic(t *testing.T) {
	var eip eips.PublicIp

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_vpc_eip.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcV1EIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1EIP_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1EIPExists(resourceName, &eip),
				),
			},
		},
	})
}

func TestAccVpcV1EIP_share(t *testing.T) {
	var eip eips.PublicIp

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_vpc_eip.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcV1EIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1EIP_share(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1EIPExists(resourceName, &eip),
				),
			},
		},
	})
}

func testAccCheckVpcV1EIPDestroy(s *terraform.State) error {
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

func testAccCheckVpcV1EIPExists(n string, eip *eips.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		networkingClient, err := config.NetworkingV1Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating networking client: %s", err)
		}

		found, err := eips.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("EIP not found")
		}

		*eip = found

		return nil
	}
}

func testAccVpcV1EIP_basic(rName string) string {
	return fmt.Sprintf(`
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
`, rName)
}

func testAccVpcV1EIP_share(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc_bandwidth" "test" {
	name = "%s"
	size = 5
}

resource "sbercloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    id         = sbercloud_vpc_bandwidth.test.id
    share_type = "WHOLE"
  }
}
`, rName)
}
