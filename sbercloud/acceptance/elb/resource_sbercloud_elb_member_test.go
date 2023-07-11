package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/pools"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccElbV3Member_basic(t *testing.T) {
	var member_1 pools.Member
	var member_2 pools.Member
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckElbV3MemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3MemberConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElbV3MemberExists("sbercloud_elb_member.member_1", &member_1),
					testAccCheckElbV3MemberExists("sbercloud_elb_member.member_2", &member_2),
				),
			},
			{
				Config: testAccElbV3MemberConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sbercloud_elb_member.member_1", "weight", "10"),
					resource.TestCheckResourceAttr("sbercloud_elb_member.member_2", "weight", "15"),
				),
			},
			{
				ResourceName:      "sbercloud_elb_member.member_1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccELBMemberImportStateIdFunc(),
			},
		},
	})
}

func TestAccElbV3Member_crossVpcBackend(t *testing.T) {
	var member_1 pools.Member
	var member_2 pools.Member
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckElbV3MemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3MemberConfig_crossVpcBackend_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElbV3MemberExists("sbercloud_elb_member.member_1", &member_1),
					testAccCheckElbV3MemberExists("sbercloud_elb_member.member_2", &member_2),
				),
			},
			{
				Config: testAccElbV3MemberConfig_crossVpcBackend_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sbercloud_elb_member.member_1", "weight", "10"),
					resource.TestCheckResourceAttr("sbercloud_elb_member.member_2", "weight", "15"),
				),
			},
			{
				ResourceName:      "sbercloud_elb_member.member_1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccELBMemberImportStateIdFunc(),
			},
		},
	})
}

func testAccCheckElbV3MemberDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	elbClient, err := cfg.ElbV3Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating ELB client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_elb_member" {
			continue
		}

		poolId := rs.Primary.Attributes["pool_id"]
		_, err := pools.GetMember(elbClient, poolId, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("member still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckElbV3MemberExists(n string, member *pools.Member) resource.TestCheckFunc {
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

		poolId := rs.Primary.Attributes["pool_id"]
		found, err := pools.GetMember(elbClient, poolId, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("member not found")
		}

		*member = *found

		return nil
	}
}

func testAccELBMemberImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		pool, ok := s.RootModule().Resources["sbercloud_elb_pool.test"]
		if !ok {
			return "", fmt.Errorf("pool not found: %s", pool)
		}
		member, ok := s.RootModule().Resources["sbercloud_elb_member.member_1"]
		if !ok {
			return "", fmt.Errorf("member not found: %s", member)
		}
		if pool.Primary.ID == "" || member.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", pool.Primary.ID, member.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", pool.Primary.ID, member.Primary.ID), nil
	}
}

func testAccElbV3MemberConfig_basic(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet_elb_test"
}

data "sbercloud_availability_zones" "test" {}

resource "sbercloud_elb_loadbalancer" "test" {
  name            = "%s"
  ipv4_subnet_id  = data.sbercloud_vpc_subnet.test.ipv4_subnet_id
  ipv6_network_id = data.sbercloud_vpc_subnet.test.id

  availability_zone = [
    data.sbercloud_availability_zones.test.names[0]
  ]
}

resource "sbercloud_elb_listener" "test" {
  name            = "%s"
  description     = "test description"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = sbercloud_elb_loadbalancer.test.id

  forward_eip      = true
  idle_timeout     = 60
  request_timeout  = 60
  response_timeout = 60
}

resource "sbercloud_elb_pool" "test" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = sbercloud_elb_listener.test.id
}

resource "sbercloud_elb_member" "member_1" {
  address       = "192.168.10.10"
  protocol_port = 8080
  pool_id       = sbercloud_elb_pool.test.id
  subnet_id     = data.sbercloud_vpc_subnet.test.ipv4_subnet_id
}

resource "sbercloud_elb_member" "member_2" {
  address       = "192.168.10.11"
  protocol_port = 8080
  pool_id       = sbercloud_elb_pool.test.id
  subnet_id     = data.sbercloud_vpc_subnet.test.ipv4_subnet_id
}
`, rName, rName, rName)
}

func testAccElbV3MemberConfig_update(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet_elb_test"
}

data "sbercloud_availability_zones" "test" {}

resource "sbercloud_elb_loadbalancer" "test" {
  name            = "%s"
  ipv4_subnet_id  = data.sbercloud_vpc_subnet.test.ipv4_subnet_id
  ipv6_network_id = data.sbercloud_vpc_subnet.test.id

  availability_zone = [
    data.sbercloud_availability_zones.test.names[0]
  ]
}

resource "sbercloud_elb_listener" "test" {
  name            = "%s"
  description     = "test description"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = sbercloud_elb_loadbalancer.test.id

  forward_eip      = true
  idle_timeout     = 60
  request_timeout  = 60
  response_timeout = 60
}

resource "sbercloud_elb_pool" "test" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = sbercloud_elb_listener.test.id
}

resource "sbercloud_elb_member" "member_1" {
  address        = "192.168.10.10"
  protocol_port  = 8080
  weight         = 10
  pool_id        = sbercloud_elb_pool.test.id
  subnet_id      = data.sbercloud_vpc_subnet.test.ipv4_subnet_id
}

resource "sbercloud_elb_member" "member_2" {
  address        = "192.168.10.11"
  protocol_port  = 8080
  weight         = 15
  pool_id        = sbercloud_elb_pool.test.id
  subnet_id      = data.sbercloud_vpc_subnet.test.ipv4_subnet_id
}
`, rName, rName, rName)
}

func testAccElbV3MemberConfig_crossVpcBackend_basic(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet_elb_test"
}

data "sbercloud_availability_zones" "test" {}

resource "sbercloud_elb_loadbalancer" "test" {
  name              = "%s"
  cross_vpc_backend = true
  ipv4_subnet_id    = data.sbercloud_vpc_subnet.test.ipv4_subnet_id
  ipv6_network_id   = data.sbercloud_vpc_subnet.test.id

  availability_zone = [
    data.sbercloud_availability_zones.test.names[0]
  ]
}

resource "sbercloud_elb_listener" "test" {
  name            = "%s"
  description     = "test description"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = sbercloud_elb_loadbalancer.test.id

  forward_eip      = true
  idle_timeout     = 60
  request_timeout  = 60
  response_timeout = 60
}

resource "sbercloud_elb_pool" "test" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = sbercloud_elb_listener.test.id
}

resource "sbercloud_elb_member" "member_1" {
  address       = "121.121.0.120"
  protocol_port = 8080
  pool_id       = sbercloud_elb_pool.test.id
}

resource "sbercloud_elb_member" "member_2" {
  address       = "121.121.0.121"
  protocol_port = 8080
  pool_id       = sbercloud_elb_pool.test.id
}
`, rName, rName, rName)
}

func testAccElbV3MemberConfig_crossVpcBackend_update(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet_elb_test"
}

data "sbercloud_availability_zones" "test" {}

resource "sbercloud_elb_loadbalancer" "test" {
  name              = "%s"
  cross_vpc_backend = true
  ipv4_subnet_id    = data.sbercloud_vpc_subnet.test.ipv4_subnet_id
  ipv6_network_id   = data.sbercloud_vpc_subnet.test.id

  availability_zone = [
    data.sbercloud_availability_zones.test.names[0]
  ]
}

resource "sbercloud_elb_listener" "test" {
  name            = "%s"
  description     = "test description"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = sbercloud_elb_loadbalancer.test.id

  forward_eip      = true
  idle_timeout     = 60
  request_timeout  = 60
  response_timeout = 60
}

resource "sbercloud_elb_pool" "test" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = sbercloud_elb_listener.test.id
}

resource "sbercloud_elb_member" "member_1" {
  address        = "121.121.0.120"
  protocol_port  = 8080
  weight         = 10
  pool_id        = sbercloud_elb_pool.test.id
}

resource "sbercloud_elb_member" "member_2" {
  address        = "121.121.0.121"
  protocol_port  = 8080
  weight         = 15
  pool_id        = sbercloud_elb_pool.test.id
}
`, rName, rName, rName)
}
