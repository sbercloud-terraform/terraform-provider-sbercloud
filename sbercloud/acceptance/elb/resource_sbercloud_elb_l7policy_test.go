package elb

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/l7policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccElbV3L7Policy_basic(t *testing.T) {
	var l7Policy l7policies.L7Policy
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_elb_l7policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckElbV3L7PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckElbV3L7PolicyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElbV3L7PolicyExists(resourceName, &l7Policy),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestMatchResourceAttr(resourceName, "listener_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckElbV3L7PolicyDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	elbClient, err := cfg.ElbV3Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating ELB client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_elb_l7policy" {
			continue
		}

		_, err := l7policies.Get(elbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("the L7 Policy still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckElbV3L7PolicyExists(n string, l7Policy *l7policies.L7Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		lbClient, err := cfg.ElbV3Client(acceptance.SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating ELB client: %s", err)
		}

		found, err := l7policies.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("policy not found")
		}

		*l7Policy = *found

		return nil
	}
}

func testAccCheckElbV3L7PolicyConfig_basic(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "sbercloud_availability_zones" "test" {}

resource "sbercloud_elb_loadbalancer" "test" {
  name            = "%s"
  ipv4_subnet_id  = data.sbercloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.sbercloud_availability_zones.test.names[0]
  ]
}

resource "sbercloud_elb_listener" "test" {
  name             = "%s"
  description      = "test description"
  protocol         = "HTTP"
  protocol_port    = 8080
  loadbalancer_id  = sbercloud_elb_loadbalancer.test.id
  forward_eip      = true
  idle_timeout     = 60
  request_timeout  = 60
  response_timeout = 60
}

resource "sbercloud_elb_pool" "test" {
  name            = "%s"
  protocol        = "HTTP"
  lb_method       = "LEAST_CONNECTIONS"
  loadbalancer_id = sbercloud_elb_loadbalancer.test.id
}

resource "sbercloud_elb_l7policy" "test" {
  name             = "%s"
  description      = "test description"
  listener_id      = sbercloud_elb_listener.test.id
  redirect_pool_id = sbercloud_elb_pool.test.id
}
`, rName, rName, rName, rName)
}
