package lb

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/lbaas_v2/pools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccLBV2Pool_basic(t *testing.T) {
	var pool pools.Pool
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameUpdate := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_lb_pool.pool_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		CheckDestroy: testAccCheckLBV2PoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2PoolConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2PoolExists(resourceName, &pool),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "ROUND_ROBIN"),
				),
			},
			{
				Config: testAccLBV2PoolConfig_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "LEAST_CONNECTIONS"),
				),
			},
		},
	})
}

func testAccCheckLBV2PoolDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	elbClient, err := config.ElbV2Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud elb client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_lb_pool" {
			continue
		}

		_, err := pools.Get(elbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Pool still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2PoolExists(n string, pool *pools.Pool) resource.TestCheckFunc {
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

		found, err := pools.Get(elbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*pool = *found

		return nil
	}
}

func testAccLBV2PoolConfig_basic(rName string) string {
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

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, rName, rName, rName)
}

func testAccLBV2PoolConfig_update(rName, rNameUpdate string) string {
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
  name           = "%s"
  protocol       = "HTTP"
  lb_method      = "LEAST_CONNECTIONS"
  admin_state_up = "true"
  listener_id    = sbercloud_lb_listener.listener_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, rName, rName, rNameUpdate)
}
