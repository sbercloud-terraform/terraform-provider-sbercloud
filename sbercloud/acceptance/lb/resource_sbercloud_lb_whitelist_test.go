package lb

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/lbaas_v2/whitelists"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccLBV2Whitelist_basic(t *testing.T) {
	var whitelist whitelists.Whitelist
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_lb_whitelist.whitelist_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV2WhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2WhitelistConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2WhitelistExists(resourceName, &whitelist),
				),
			},
			{
				Config: testAccLBV2WhitelistConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enable_whitelist", "true"),
				),
			},
		},
	})
}

func testAccCheckLBV2WhitelistDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	elbClient, err := config.ElbV2Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud elb client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_lb_whitelist" {
			continue
		}

		_, err := whitelists.Get(elbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Whitelist still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2WhitelistExists(n string, whitelist *whitelists.Whitelist) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		elbClient, err := config.ElbV2Client(acceptance.SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud elb client: %s", err)
		}

		found, err := whitelists.Get(elbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Whitelist not found")
		}

		*whitelist = *found

		return nil
	}
}

func testAccLBV2WhitelistConfig_basic(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "sbercloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "%s"
  vip_subnet_id = data.sbercloud_vpc_subnet.test.subnet_id
}

resource "sbercloud_lb_listener" "listener_1" {
  name            = "%s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = sbercloud_lb_loadbalancer.loadbalancer_1.id
}

resource "sbercloud_lb_whitelist" "whitelist_1" {
  enable_whitelist = true
  whitelist        = "192.168.11.1,192.168.0.1/24"
  listener_id      = sbercloud_lb_listener.listener_1.id
}
`, rName, rName)
}

func testAccLBV2WhitelistConfig_update(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "sbercloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "%s"
  vip_subnet_id = data.sbercloud_vpc_subnet.test.subnet_id
}

resource "sbercloud_lb_listener" "listener_1" {
  name            = "%s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = sbercloud_lb_loadbalancer.loadbalancer_1.id
}

resource "sbercloud_lb_whitelist" "whitelist_1" {
  enable_whitelist = true
  whitelist        = "192.168.11.1,192.168.0.1/24,192.168.201.18/8"
  listener_id      = sbercloud_lb_listener.listener_1.id
}
`, rName, rName)
}
