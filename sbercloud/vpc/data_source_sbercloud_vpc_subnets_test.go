package vpc

import (
	"fmt"
	"testing"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVpcSubnetsDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.sbercloud_vpc_subnets.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetsDataSource_Base(randName, randCidr, randGatewayIp),
			},
			{
				Config: testAccVpcSubnetsDataSource_Basic(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.dhcp_enable", "true"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "subnets.0.vpc_id",
						"${sbercloud_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccVpcSubnetsDataSource_Basic(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc_subnets" "test" {
  id = sbercloud_vpc_subnet.test.id
}
`, testAccVpcSubnetsDataSource_Base(rName, cidr, gatewayIp))
}

func TestAccVpcSubnetsDataSource_ipv4ByCidr(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.sbercloud_vpc_subnets.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetsDataSource_Base(randName, randCidr, randGatewayIp),
			},
			{
				Config: testAccVpcSubnetsDataSource_ipv4ByCidr(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.dhcp_enable", "true"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "subnets.0.vpc_id",
						"${sbercloud_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccVpcSubnetsDataSource_ipv4ByCidr(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc_subnets" "test" {
  cidr = sbercloud_vpc_subnet.test.cidr
}
`, testAccVpcSubnetsDataSource_Base(rName, cidr, gatewayIp))
}

func TestAccVpcSubnetsDataSource_ipv4ByName(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.sbercloud_vpc_subnets.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetsDataSource_Base(randName, randCidr, randGatewayIp),
			},
			{
				Config: testAccVpcSubnetsDataSource_ipv4ByName(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.dhcp_enable", "true"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "subnets.0.vpc_id",
						"${sbercloud_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccVpcSubnetsDataSource_ipv4ByName(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc_subnets" "test" {
  name = sbercloud_vpc_subnet.test.name
}
`, testAccVpcSubnetsDataSource_Base(rName, cidr, gatewayIp))
}

func TestAccVpcSubnetsDataSource_ipv4ByVpcId(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.sbercloud_vpc_subnets.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetsDataSource_Base(randName, randCidr, randGatewayIp),
			},
			{
				Config: testAccVpcSubnetsDataSource_ipv4ByVpcId(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.dhcp_enable", "true"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "subnets.0.vpc_id",
						"${sbercloud_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccVpcSubnetsDataSource_ipv4ByVpcId(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc_subnets" "test" {
  vpc_id = sbercloud_vpc_subnet.test.vpc_id
}
`, testAccVpcSubnetsDataSource_Base(rName, cidr, gatewayIp))
}

func testAccVpcSubnetsDataSource_Base(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc" "test" {
  name = "%s"
  cidr = "%s"
}

resource "sbercloud_vpc_subnet" "test" {
  name       = "%s"
  vpc_id     = sbercloud_vpc.test.id
  cidr       = "%s"
  gateway_ip = "%s"
}`, rName, cidr, rName, cidr, gatewayIp)
}
