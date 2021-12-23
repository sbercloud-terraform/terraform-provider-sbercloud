package vpc

import (
	"fmt"
	"testing"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVpcSubnetDataSource_ipv4Basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.sbercloud_vpc_subnet.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetDataSource_ipv4Base(randName, randCidr, randGatewayIp),
			},
			{
				Config: testAccVpcSubnetDataSource_ipv4Basic(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "dhcp_enable", "true"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpc_id",
						"${sbercloud_vpc.test.id}"),
				),
			},
		},
	})
}

func TestAccVpcSubnetDataSource_ipv4ByCidr(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.sbercloud_vpc_subnet.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetDataSource_ipv4Base(randName, randCidr, randGatewayIp),
			},
			{
				Config: testAccVpcSubnetDataSource_ipv4ByCidr(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "dhcp_enable", "true"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpc_id",
						"${sbercloud_vpc.test.id}"),
				),
			},
		},
	})
}

func TestAccVpcSubnetDataSource_ipv4ByName(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.sbercloud_vpc_subnet.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetDataSource_ipv4Base(randName, randCidr, randGatewayIp),
			},
			{
				Config: testAccVpcSubnetDataSource_ipv4ByName(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "dhcp_enable", "true"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpc_id",
						"${sbercloud_vpc.test.id}"),
				),
			},
		},
	})
}

func TestAccVpcSubnetDataSource_ipv4ByVpcId(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.sbercloud_vpc_subnet.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetDataSource_ipv4Base(randName, randCidr, randGatewayIp),
			},
			{
				Config: testAccVpcSubnetDataSource_ipv4ByVpcId(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "dhcp_enable", "true"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpc_id",
						"${sbercloud_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccVpcSubnetDataSource_ipv4Base(rName, cidr, gatewayIp string) string {
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

func testAccVpcSubnetDataSource_ipv4Basic(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc_subnet" "test" {
  id = sbercloud_vpc_subnet.test.id
}
`, testAccVpcSubnetDataSource_ipv4Base(rName, cidr, gatewayIp))
}

func testAccVpcSubnetDataSource_ipv4ByCidr(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc_subnet" "test" {
  cidr = sbercloud_vpc_subnet.test.cidr
}
`, testAccVpcSubnetDataSource_ipv4Base(rName, cidr, gatewayIp))
}

func testAccVpcSubnetDataSource_ipv4ByName(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc_subnet" "test" {
  name = sbercloud_vpc_subnet.test.name
}
`, testAccVpcSubnetDataSource_ipv4Base(rName, cidr, gatewayIp))
}

func testAccVpcSubnetDataSource_ipv4ByVpcId(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc_subnet" "test" {
  vpc_id = sbercloud_vpc_subnet.test.vpc_id
}
`, testAccVpcSubnetDataSource_ipv4Base(rName, cidr, gatewayIp))
}
