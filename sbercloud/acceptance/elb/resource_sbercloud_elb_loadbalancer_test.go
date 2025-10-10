package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/loadbalancers"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func getELBResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.ElbV3Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB client: %s", err)
	}

	eipID := state.Primary.Attributes["ipv4_eip_id"]
	eipType := state.Primary.Attributes["iptype"]
	if eipType != "" && eipID != "" {
		eipClient, err := c.NetworkingV1Client(acceptance.SBC_REGION_NAME)
		if err != nil {
			return nil, fmt.Errorf("error creating VPC v1 client: %s", err)
		}

		if _, err := eips.Get(eipClient, eipID).Extract(); err != nil {
			return nil, err
		}
	}

	return loadbalancers.Get(client, state.Primary.ID).Extract()
}

func TestAccElbV3LoadBalancer_basic(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_backend", "false"),
					resource.TestCheckResourceAttr(resourceName, "backend_subnets.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "backend_subnets.0",
						"sbercloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "nonProtection"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "charge_mode", "lcu"),
					resource.TestCheckResourceAttr(resourceName, "guaranteed", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_failure_action", "discard"),
					resource.TestCheckResourceAttrSet(resourceName, "elb_virsubnet_type"),
					resource.TestCheckResourceAttrSet(resourceName, "operating_status"),
					resource.TestCheckResourceAttrSet(resourceName, "public_border_group"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccElbV3LoadBalancerConfig_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_backend", "true"),
					resource.TestCheckResourceAttr(resourceName, "backend_subnets.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "consoleProtection"),
					resource.TestCheckResourceAttr(resourceName, "protection_reason", "test protection reason"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
					resource.TestCheckResourceAttr(resourceName, "guaranteed", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "waf_failure_action", "forward"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection_enable"},
			},
		},
	})
}

func TestAccElbV3LoadBalancer_with_deletion_protection(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_with_deletion_protection(rName, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				Config: testAccElbV3LoadBalancerConfig_with_deletion_protection(rName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func TestAccElbV3LoadBalancer_withEpsId(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_withEpsId(rName, "0"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				Config: testAccElbV3LoadBalancerConfig_withEpsId(rName, acceptance.SBC_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.SBC_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccElbV3LoadBalancer_withEIP(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_withEIP(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "iptype", "5_bgp"),
					resource.TestCheckResourceAttrSet(resourceName, "ipv4_eip_id"),
				),
			},
		},
	})
}

func TestAccElbV3LoadBalancer_withEIP_Bandwidth_Id(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_withEIP_Bandwidth_Id(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "iptype", "5_bgp"),
					resource.TestCheckResourceAttrSet(resourceName, "ipv4_eip_id"),
				),
			},
		},
	})
}

func TestAccElbV3LoadBalancer_availabilityZone(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_availabilityZone(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_backend", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone.0",
						"data.sbercloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccElbV3LoadBalancerConfig_availabilityZoneUpdate(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_backend", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone.0",
						"data.sbercloud_availability_zones.test", "names.1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
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

// API does not support updating availability_zone_list (maybe not for all tenants)
func TestAccElbV3LoadBalancer_withIpv6(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_withIpv6(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "ipv6_network_id",
						"data.sbercloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "ipv6_address", "2407:c080:1200:5f0:9bb8:4438:299b:9083"),
				),
			},
			{
				Config: testAccElbV3LoadBalancerConfig_withIpv6_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "ipv6_network_id",
						"data.sbercloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "ipv6_address", "2407:c080:1200:5f0:9bb8:4438:299b:9084"),
				),
			},
		},
	})
}

func TestAccElbV3LoadBalancer_gateway(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckElbGatewayType(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_gateway(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "loadbalancer_type", "gateway"),
					resource.TestCheckResourceAttrPair(resourceName, "ipv4_subnet_id",
						"sbercloud_vpc_subnet.test", "ipv4_subnet_id"),
					resource.TestCheckResourceAttrPair(resourceName, "ipv6_network_id",
						"sbercloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "description", "test gateway description"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "gw_flavor_id"),
				),
			},
			{
				Config: testAccElbV3LoadBalancerConfig_gateway_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "loadbalancer_type", "gateway"),
					resource.TestCheckResourceAttrPair(resourceName, "ipv4_subnet_id",
						"sbercloud_vpc_subnet.test", "ipv4_subnet_id"),
					resource.TestCheckResourceAttr(resourceName, "description", "test gateway description update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
					resource.TestCheckResourceAttrSet(resourceName, "gw_flavor_id"),
				),
			},
		},
	})
}

func testAccElbV3LoadBalancerConfig_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "sbercloud_availability_zones" "test" {}

resource "sbercloud_vpc_subnet" "test_1" {
  name       = "%[2]s_1"
  vpc_id     = sbercloud_vpc.test.id
  cidr       = "192.168.1.0/24"
  gateway_ip = "192.168.1.1"
}

resource "sbercloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "sbercloud_elb_loadbalancer" "test" {
  name                = "%[2]s"
  vpc_id              = sbercloud_vpc.test.id
  ipv4_subnet_id      = sbercloud_vpc_subnet.test.ipv4_subnet_id
  ipv4_eip_id         = sbercloud_vpc_eip.test.id
  waf_failure_action  = "discard"
  autoscaling_enabled = true

  availability_zone = [
    data.sbercloud_availability_zones.test.names[0]
  ]

  backend_subnets = [
    sbercloud_vpc_subnet.test.id
  ]

  protection_status = "nonProtection"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, acceptance.TestVpc(rName), rName)
}

func testAccElbV3LoadBalancerConfig_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%[1]s

data "sbercloud_availability_zones" "test" {}

resource "sbercloud_vpc_subnet" "test_1" {
  name       = "%[2]s_1"
  vpc_id     = sbercloud_vpc.test.id
  cidr       = "192.168.1.0/24"
  gateway_ip = "192.168.1.1"
}

resource "sbercloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "sbercloud_elb_loadbalancer" "test" {
  name                = "%[3]s"
  cross_vpc_backend   = true
  vpc_id              = sbercloud_vpc.test.id
  ipv4_subnet_id      = sbercloud_vpc_subnet.test.ipv4_subnet_id
  ipv4_eip_id         = sbercloud_vpc_eip.test.id
  waf_failure_action  = "forward"
  autoscaling_enabled = true

  availability_zone = [
    data.sbercloud_availability_zones.test.names[0]
  ]

  backend_subnets = [
    sbercloud_vpc_subnet.test.id,
    sbercloud_vpc_subnet.test_1.id,
  ]

  protection_status = "consoleProtection"
  protection_reason = "test protection reason"

  tags = {
    key1  = "value1"
    owner = "terraform_update"
  }
}
`, acceptance.TestVpc(rName), rName, rNameUpdate)
}

func testAccElbV3LoadBalancerConfig_with_deletion_protection(rName string, deletionProtection bool) string {
	return fmt.Sprintf(`
%[1]s

data "sbercloud_availability_zones" "test" {}

resource "sbercloud_vpc_subnet" "test_1" {
  name       = "%[2]s_1"
  vpc_id     = sbercloud_vpc.test.id
  cidr       = "192.168.1.0/24"
  gateway_ip = "192.168.1.1"
}

resource "sbercloud_elb_loadbalancer" "test" {
  name                       = "%[2]s"
  vpc_id                     = sbercloud_vpc.test.id
  ipv4_subnet_id             = sbercloud_vpc_subnet.test.ipv4_subnet_id
  deletion_protection_enable = %[3]v

  availability_zone = [
    data.sbercloud_availability_zones.test.names[0]
  ]
}
`, acceptance.TestVpc(rName), rName, deletionProtection)
}

func testAccElbV3LoadBalancerConfig_withEpsId(rName, enterpriseProjectId string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "sbercloud_availability_zones" "test" {}

resource "sbercloud_elb_loadbalancer" "test" {
  name                  = "%[1]s"
  ipv4_subnet_id        = data.sbercloud_vpc_subnet.test.ipv4_subnet_id
  enterprise_project_id = "%[2]s"

  availability_zone = [
    data.sbercloud_availability_zones.test.names[0]
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, rName, enterpriseProjectId)
}

func testAccElbV3LoadBalancerConfig_withEIP(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "sbercloud_availability_zones" "test" {}

resource "sbercloud_elb_loadbalancer" "test" {
  name           = "%s"
  ipv4_subnet_id = data.sbercloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.sbercloud_availability_zones.test.names[0]
  ]

  iptype                = "5_bgp"
  bandwidth_charge_mode = "traffic"
  sharetype             = "PER"
  bandwidth_size        = 5
}
`, rName)
}

func testAccElbV3LoadBalancerConfig_withEIP_Bandwidth_Id(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "sbercloud_availability_zones" "test" {}

resource "sbercloud_vpc_bandwidth" "test" {
  name = "%[1]s"
  size = 5
}

resource "sbercloud_elb_loadbalancer" "test" {
  name           = "%[1]s"
  ipv4_subnet_id = data.sbercloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.sbercloud_availability_zones.test.names[0]
  ]

  iptype       = "5_bgp"
  sharetype    = "WHOLE"
  bandwidth_id = sbercloud_vpc_bandwidth.test.id
}
`, rName)
}

func testAccElbV3LoadBalancerConfig_postPaid(rName string) string {
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

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, rName)
}

func testAccElbV3LoadBalancerConfig_prePaid(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "sbercloud_availability_zones" "test" {}

resource "sbercloud_elb_loadbalancer" "test" {
  name            = "%s"
  ipv4_subnet_id  = data.sbercloud_vpc_subnet.test.ipv4_subnet_id

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "false"

  availability_zone = [
    data.sbercloud_availability_zones.test.names[0]
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, rName)
}

func testAccElbV3LoadBalancerConfig_prePaidUpdate(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "sbercloud_availability_zones" "test" {}

data "sbercloud_elb_flavors" "l4flavors" {
  type            = "L4"
  max_connections = 1000000
  cps             = 20000
  bandwidth       = 100
}

data "sbercloud_elb_flavors" "l7flavors" {
  type            = "L7"
  max_connections = 400000
  cps             = 4000
  bandwidth       = 100
}

resource "sbercloud_elb_loadbalancer" "test" {
  name            = "%s"
  ipv4_subnet_id  = data.sbercloud_vpc_subnet.test.ipv4_subnet_id
  description     = "update flavors"
  l4_flavor_id    = data.sbercloud_elb_flavors.l4flavors.ids[0]
  l7_flavor_id    = data.sbercloud_elb_flavors.l7flavors.ids[0]

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"

  availability_zone = [
    data.sbercloud_availability_zones.test.names[1]
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, rName)
}

func testAccElbV3LoadBalancerConfig_availabilityZone(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "sbercloud_availability_zones" "test" {}

resource "sbercloud_elb_loadbalancer" "test" {
  name           = "%[2]s"
  vpc_id         = sbercloud_vpc.test.id
  ipv4_subnet_id = sbercloud_vpc_subnet.test.ipv4_subnet_id
	
  availability_zone = [
    data.sbercloud_availability_zones.test.names[0]
  ]

  backend_subnets = [
    sbercloud_vpc_subnet.test.id
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, acceptance.TestVpc(rName), rName)
}

func testAccElbV3LoadBalancerConfig_availabilityZoneUpdate(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%[1]s

data "sbercloud_availability_zones" "test" {}

resource "sbercloud_elb_loadbalancer" "test" {
  name              = "%[2]s"
  cross_vpc_backend = true
  vpc_id            = sbercloud_vpc.test.id
  ipv4_subnet_id    = sbercloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.sbercloud_availability_zones.test.names[1]
  ]

  tags = {
    key1  = "value1"
    owner = "terraform_update"
  }
}
`, acceptance.TestVpc(rName), rNameUpdate)
}

func testAccElbV3LoadBalancerConfig_withIpv6(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "sbercloud_availability_zones" "test" {}

resource "sbercloud_elb_loadbalancer" "test" {
  name           = "%s"
  ipv4_subnet_id = data.sbercloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.sbercloud_availability_zones.test.names[0]
  ]

  iptype                = "5_bgp"
  bandwidth_charge_mode = "traffic"
  sharetype             = "PER"
  bandwidth_size        = 5

  ipv6_network_id = data.sbercloud_vpc_subnet.test.id
  ipv6_address    = "2407:c080:1200:5f0:9bb8:4438:299b:9083"
}
`, rName)
}

func testAccElbV3LoadBalancerConfig_withIpv6_update(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "sbercloud_availability_zones" "test" {}

resource "sbercloud_elb_loadbalancer" "test" {
  name           = "%s"
  ipv4_subnet_id = data.sbercloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.sbercloud_availability_zones.test.names[0]
  ]

  iptype                = "5_bgp"
  bandwidth_charge_mode = "traffic"
  sharetype             = "PER"
  bandwidth_size        = 5

  ipv6_network_id = data.sbercloud_vpc_subnet.test.id
  ipv6_address    = "2407:c080:1200:5f0:9bb8:4438:299b:9084"
}
`, rName)
}

func testAccElbV3LoadBalancerConfig_gateway_base(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_availability_zones" "test" {}

resource "sbercloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_subnet" "test" {
  name        = "%[1]s"
  vpc_id      = sbercloud_vpc.test.id
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  ipv6_enable = true
}
`, rName)
}

func testAccElbV3LoadBalancerConfig_gateway(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_elb_loadbalancer" "test" {
  name              = "%[2]s"
  vpc_id            = sbercloud_vpc.test.id
  ipv4_subnet_id    = sbercloud_vpc_subnet.test.ipv4_subnet_id
  ipv6_network_id   = sbercloud_vpc_subnet.test.id
  loadbalancer_type = "gateway"
  description       = "test gateway description"
  availability_zone = [
    data.sbercloud_availability_zones.test.names[0]
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccElbV3LoadBalancerConfig_gateway_base(rName), rName)
}

func testAccElbV3LoadBalancerConfig_gateway_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_elb_loadbalancer" "test" {
  name              = "%[2]s"
  vpc_id            = sbercloud_vpc.test.id
  ipv4_subnet_id    = sbercloud_vpc_subnet.test.ipv4_subnet_id
  loadbalancer_type = "gateway"
  description       = "test gateway description update"
  description       = "test gateway description update"
  availability_zone = [
    data.sbercloud_availability_zones.test.names[0]
  ]

  tags = {
    key1  = "value_update"
    owner = "terraform_update"
  }
}
`, testAccElbV3LoadBalancerConfig_gateway_base(rName), rName, rNameUpdate)
}
