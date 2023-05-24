package lb

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/security/groups"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getLoadBalancerResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.LoadBalancerClient(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB v2 Client: %s", err)
	}
	resp, err := loadbalancers.Get(c, state.Primary.ID).Extract()
	if resp == nil && err == nil {
		return resp, fmt.Errorf("unable to find the LoadBalancer (%s)", state.Primary.ID)
	}
	return resp, err
}

func TestAccLBV2LoadBalancer_basic(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameUpdate := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_lb_loadbalancer.loadbalancer_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		CheckDestroy: testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2LoadBalancerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists(resourceName, &lb),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestMatchResourceAttr(resourceName, "vip_port_id",
						regexp.MustCompile("^[a-f0-9-]+")),
				),
			},
			{
				Config: testAccLBV2LoadBalancerConfig_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
				),
			},
		},
	})
}

func TestAccLBV2LoadBalancer_secGroup(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	var sg_1, sg_2 groups.SecGroup
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameSecg1 := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameSecg2 := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_lb_loadbalancer.loadbalancer_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		CheckDestroy: testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2LoadBalancer_secGroup(rName, rNameSecg1, rNameSecg2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists(resourceName, &lb),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "1"),
					testAccCheckNetworkingV2SecGroupExists(
						"sbercloud_networking_secgroup.secgroup_1", &sg_1),
					testAccCheckNetworkingV2SecGroupExists(
						"sbercloud_networking_secgroup.secgroup_1", &sg_2),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg_1),
				),
			},
			{
				Config: testAccLBV2LoadBalancer_secGroup_update1(rName, rNameSecg1, rNameSecg2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists(resourceName, &lb),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "2"),
					testAccCheckNetworkingV2SecGroupExists(
						"sbercloud_networking_secgroup.secgroup_2", &sg_1),
					testAccCheckNetworkingV2SecGroupExists(
						"sbercloud_networking_secgroup.secgroup_2", &sg_2),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg_1),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg_2),
				),
			},
			{
				Config: testAccLBV2LoadBalancer_secGroup_update2(rName, rNameSecg1, rNameSecg2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists(resourceName, &lb),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "1"),
					testAccCheckNetworkingV2SecGroupExists(
						"sbercloud_networking_secgroup.secgroup_2", &sg_1),
					testAccCheckNetworkingV2SecGroupExists(
						"sbercloud_networking_secgroup.secgroup_2", &sg_2),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg_2),
				),
			},
		},
	})
}

func testAccCheckLBV2LoadBalancerDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	elbClient, err := config.ElbV2Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud elb client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_lb_loadbalancer" {
			continue
		}

		_, err := loadbalancers.Get(elbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("LoadBalancer still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2LoadBalancerExists(
	n string, lb *loadbalancers.LoadBalancer) resource.TestCheckFunc {
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
			return fmt.Errorf("Error creating SberCloud networking client: %s", err)
		}

		found, err := loadbalancers.Get(elbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*lb = *found

		return nil
	}
}

func testAccCheckLBV2LoadBalancerHasSecGroup(
	lb *loadbalancers.LoadBalancer, sg *groups.SecGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acceptance.TestAccProvider.Meta().(*config.Config)
		networkingClient, err := config.NetworkingV2Client(acceptance.SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud networking client: %s", err)
		}

		port, err := ports.Get(networkingClient, lb.VipPortID).Extract()
		if err != nil {
			return err
		}

		for _, p := range port.SecurityGroups {
			if p == sg.ID {
				return nil
			}
		}

		return fmt.Errorf("LoadBalancer does not have the security group")
	}
}

func testAccCheckNetworkingV2SecGroupExists(n string, security_group *groups.SecGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		networkingClient, err := config.NetworkingV2Client(acceptance.SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}

		found, err := groups.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Security group not found")
		}

		*security_group = *found

		return nil
	}
}

func testAccLBV2LoadBalancerConfig_basic(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "sbercloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "%s"
  vip_subnet_id = data.sbercloud_vpc_subnet.test.subnet_id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, rName)
}

func testAccLBV2LoadBalancerConfig_update(rNameUpdate string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "sbercloud_lb_loadbalancer" "loadbalancer_1" {
  name           = "%s"
  admin_state_up = "true"
  vip_subnet_id  = data.sbercloud_vpc_subnet.test.subnet_id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, rNameUpdate)
}

func testAccLBV2LoadBalancer_secGroup(rName, rNameSecg1, rNameSecg2 string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "sbercloud_networking_secgroup" "secgroup_1" {
  name        = "%s"
  description = "secgroup_1"
}

resource "sbercloud_networking_secgroup" "secgroup_2" {
  name        = "%s"
  description = "secgroup_2"
}

resource "sbercloud_lb_loadbalancer" "loadbalancer_1" {
  name               = "%s"
  vip_subnet_id      = data.sbercloud_vpc_subnet.test.subnet_id
  security_group_ids = [
    sbercloud_networking_secgroup.secgroup_1.id
  ]
}
`, rNameSecg1, rNameSecg2, rName)
}

func testAccLBV2LoadBalancer_secGroup_update1(rName, rNameSecg1, rNameSecg2 string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "sbercloud_networking_secgroup" "secgroup_1" {
  name        = "%s"
  description = "secgroup_1"
}

resource "sbercloud_networking_secgroup" "secgroup_2" {
  name        = "%s"
  description = "secgroup_2"
}

resource "sbercloud_lb_loadbalancer" "loadbalancer_1" {
  name               = "%s"
  vip_subnet_id      = data.sbercloud_vpc_subnet.test.subnet_id
  security_group_ids = [
    sbercloud_networking_secgroup.secgroup_1.id,
    sbercloud_networking_secgroup.secgroup_2.id
  ]
}
`, rNameSecg1, rNameSecg2, rName)
}

func testAccLBV2LoadBalancer_secGroup_update2(rName, rNameSecg1, rNameSecg2 string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "sbercloud_networking_secgroup" "secgroup_1" {
  name        = "%s"
  description = "secgroup_1"
}

resource "sbercloud_networking_secgroup" "secgroup_2" {
  name        = "%s"
  description = "secgroup_2"
}

resource "sbercloud_lb_loadbalancer" "loadbalancer_1" {
  name               = "%s"
  vip_subnet_id      = data.sbercloud_vpc_subnet.test.subnet_id
  security_group_ids = [
    sbercloud_networking_secgroup.secgroup_2.id
  ]
}
`, rNameSecg1, rNameSecg2, rName)
}

