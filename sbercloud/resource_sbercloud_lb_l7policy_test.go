package sbercloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/lbaas_v2/l7policies"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccLBV2L7Policy_basic(t *testing.T) {
	var l7Policy l7policies.L7Policy
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_lb_l7policy.l7policy_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2L7PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLBV2L7PolicyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7PolicyExists(resourceName, &l7Policy),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "action", "REDIRECT_TO_POOL"),
					resource.TestCheckResourceAttr(resourceName, "position", "1"),
					resource.TestMatchResourceAttr(resourceName, "listener_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
				),
			},
		},
	})
}

func testAccCheckLBV2L7PolicyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	lbClient, err := config.ElbV2Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Sbercloud load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_lb_l7policy" {
			continue
		}

		_, err := l7policies.Get(lbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("L7 Policy still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2L7PolicyExists(n string, l7Policy *l7policies.L7Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		lbClient, err := config.ElbV2Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating Sbercloud load balancing client: %s", err)
		}

		found, err := l7policies.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Policy not found")
		}

		*l7Policy = *found

		return nil
	}
}

func testAccCheckLBV2L7PolicyConfig_basic(rName string) string {
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
  name            = "%s"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = sbercloud_lb_loadbalancer.loadbalancer_1.id
}

resource "sbercloud_lb_l7policy" "l7policy_1" {
  name         = "%s"
  action       = "REDIRECT_TO_POOL"
  description  = "test description"
  position     = 1
  listener_id  = sbercloud_lb_listener.listener_1.id
  redirect_pool_id = sbercloud_lb_pool.pool_1.id
}
`, rName, rName, rName, rName)
}
