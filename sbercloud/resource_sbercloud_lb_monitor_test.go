package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/lbaas_v2/monitors"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccLBV2Monitor_basic(t *testing.T) {
	var monitor monitors.Monitor
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameUpdate := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_lb_monitor.monitor_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2MonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2MonitorConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MonitorExists(resourceName, &monitor),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "delay", "20"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
				),
			},
			{
				Config: testAccLBV2MonitorConfig_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "delay", "30"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "15"),
				),
			},
		},
	})
}

func testAccCheckLBV2MonitorDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	elbClient, err := config.ElbV2Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud elb client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_lb_monitor" {
			continue
		}

		_, err := monitors.Get(elbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Monitor still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2MonitorExists(n string, monitor *monitors.Monitor) resource.TestCheckFunc {
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

		found, err := monitors.Get(elbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Monitor not found")
		}

		*monitor = *found

		return nil
	}
}

func testAccLBV2MonitorConfig_basic(rName string) string {
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

resource "sbercloud_lb_pool" "pool_1" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = sbercloud_lb_listener.listener_1.id
}

resource "sbercloud_lb_monitor" "monitor_1" {
  name        = "%s"
  type        = "PING"
  delay       = 20
  timeout     = 10
  max_retries = 5
  pool_id     = sbercloud_lb_pool.pool_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, rName, rName, rName, rName)
}

func testAccLBV2MonitorConfig_update(rName, rNameUpdate string) string {
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

resource "sbercloud_lb_pool" "pool_1" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = sbercloud_lb_listener.listener_1.id
}

resource "sbercloud_lb_monitor" "monitor_1" {
  name           = "%s"
  type           = "PING"
  delay          = 30
  timeout        = 15
  max_retries    = 10
  admin_state_up = "true"
  pool_id        = sbercloud_lb_pool.pool_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, rName, rName, rName, rNameUpdate)
}
