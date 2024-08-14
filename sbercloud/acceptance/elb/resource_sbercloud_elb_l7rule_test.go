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

func TestAccElbV3L7Rule_basic(t *testing.T) {
	var l7rule l7policies.Rule
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_elb_l7rule.l7rule_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckElbV3L7RuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckElbV3L7RuleConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElbV3L7RuleExists(resourceName, &l7rule),
					resource.TestCheckResourceAttr(resourceName, "type", "PATH"),
					resource.TestCheckResourceAttr(resourceName, "compare_type", "EQUAL_TO"),
					resource.TestCheckResourceAttr(resourceName, "value", "/api"),
					resource.TestMatchResourceAttr(resourceName, "l7policy_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
				),
			},
			{
				Config: testAccCheckElbV3L7RuleConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElbV3L7RuleExists(resourceName, &l7rule),
					resource.TestCheckResourceAttr(resourceName, "type", "PATH"),
					resource.TestCheckResourceAttr(resourceName, "compare_type", "STARTS_WITH"),
					resource.TestCheckResourceAttr(resourceName, "value", "/images"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccELBL7RuleImportStateIdFunc(),
			},
		},
	})
}

func testAccCheckElbV3L7RuleDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	elbClient, err := cfg.ElbV3Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating ELB client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_elb_l7rule" {
			continue
		}

		l7policyID := ""
		for k, v := range rs.Primary.Attributes {
			if k == "l7policy_id" {
				l7policyID = v
				break
			}
		}

		if l7policyID == "" {
			return fmt.Errorf("unable to find l7policy_id")
		}

		_, err := l7policies.GetRule(elbClient, l7policyID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("the L7 Rule still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckElbV3L7RuleExists(n string, l7rule *l7policies.Rule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		elbClient, err := cfg.ElbV3Client(acceptance.SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating ELB client: %s", err)
		}

		l7policyID := ""
		for k, v := range rs.Primary.Attributes {
			if k == "l7policy_id" {
				l7policyID = v
				break
			}
		}

		if l7policyID == "" {
			return fmt.Errorf("unable to find l7policy_id")
		}

		found, err := l7policies.GetRule(elbClient, l7policyID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("policy not found")
		}

		*l7rule = *found

		return nil
	}
}

func testAccELBL7RuleImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		policy, ok := s.RootModule().Resources["sbercloud_elb_l7policy.test"]
		if !ok {
			return "", fmt.Errorf("policy not found: %s", policy)
		}
		rule, ok := s.RootModule().Resources["sbercloud_elb_l7rule.l7rule_1"]
		if !ok {
			return "", fmt.Errorf("rule not found: %s", rule)
		}
		if policy.Primary.ID == "" || rule.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", policy.Primary.ID, rule.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", policy.Primary.ID, rule.Primary.ID), nil
	}
}

func testAccCheckElbV3L7RuleConfig(rName string) string {
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

func testAccCheckElbV3L7RuleConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_elb_l7rule" "l7rule_1" {
  l7policy_id  = sbercloud_elb_l7policy.test.id
  type         = "PATH"
  compare_type = "EQUAL_TO"
  value        = "/api"
}
`, testAccCheckElbV3L7RuleConfig(rName))
}

func testAccCheckElbV3L7RuleConfig_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_elb_l7rule" "l7rule_1" {
  l7policy_id  = sbercloud_elb_l7policy.test.id
  type         = "PATH"
  compare_type = "STARTS_WITH"
  value        = "/images"
}
`, testAccCheckElbV3L7RuleConfig(rName))
}
