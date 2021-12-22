package sbercloud

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/lbaas_v2/listeners"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccLBV2Listener_basic(t *testing.T) {
	var listener listeners.Listener
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameUpdate := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_lb_listener.listener_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2ListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2ListenerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2ListenerExists(resourceName, &listener),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "connection_limit", "-1"),
				),
			},
			{
				Config: testAccLBV2ListenerConfig_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
				),
			},
		},
	})
}

func testAccCheckLBV2ListenerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	elbClient, err := config.ElbV2Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud elb client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_lb_listener" {
			continue
		}

		_, err := listeners.Get(elbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Listener still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2ListenerExists(n string, listener *listeners.Listener) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		elbClient, err := config.ElbV2Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud elb client: %s", err)
		}

		found, err := listeners.Get(elbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*listener = *found

		return nil
	}
}

func testAccLBV2ListenerConfig_basic(rName string) string {
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
`, rName, rName)
}

func testAccLBV2ListenerConfig_update(rName, rNameUpdate string) string {
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
  admin_state_up  = "true"
  loadbalancer_id = sbercloud_lb_loadbalancer.loadbalancer_1.id
}
`, rName, rNameUpdate)
}
