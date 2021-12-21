package sbercloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVpcSubnetV1DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dNameByID := "data.sbercloud_vpc_subnet.by_id"
	dNameByCIDR := "data.sbercloud_vpc_subnet.by_cidr"
	dNameByName := "data.sbercloud_vpc_subnet.by_name"
	dNameByVpcID := "data.sbercloud_vpc_subnet.by_vpc_id"
	tmp := strconv.Itoa(acctest.RandIntRange(1, 254))
	cidr := fmt.Sprintf("172.16.%s.0/24", tmp)
	gateway := fmt.Sprintf("172.16.%s.1", tmp)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetV1DataSource_basic(rName, cidr, gateway),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetV1DataSourceID(dNameByID),
					resource.TestCheckResourceAttr(dNameByID, "name", rName),
					resource.TestCheckResourceAttr(dNameByID, "cidr", cidr),
					resource.TestCheckResourceAttr(dNameByID, "gateway_ip", gateway),
					resource.TestCheckResourceAttr(dNameByID, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dNameByID, "dhcp_enable", "true"),
					testAccCheckVpcSubnetV1DataSourceID(dNameByCIDR),
					resource.TestCheckResourceAttr(dNameByCIDR, "name", rName),
					resource.TestCheckResourceAttr(dNameByCIDR, "cidr", cidr),
					resource.TestCheckResourceAttr(dNameByCIDR, "gateway_ip", gateway),
					resource.TestCheckResourceAttr(dNameByCIDR, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dNameByCIDR, "dhcp_enable", "true"),
					testAccCheckVpcSubnetV1DataSourceID(dNameByName),
					resource.TestCheckResourceAttr(dNameByName, "name", rName),
					resource.TestCheckResourceAttr(dNameByName, "cidr", cidr),
					resource.TestCheckResourceAttr(dNameByName, "gateway_ip", gateway),
					resource.TestCheckResourceAttr(dNameByName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dNameByName, "dhcp_enable", "true"),
					testAccCheckVpcSubnetV1DataSourceID(dNameByVpcID),
					resource.TestCheckResourceAttr(dNameByVpcID, "name", rName),
					resource.TestCheckResourceAttr(dNameByVpcID, "cidr", cidr),
					resource.TestCheckResourceAttr(dNameByVpcID, "gateway_ip", gateway),
					resource.TestCheckResourceAttr(dNameByVpcID, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dNameByVpcID, "dhcp_enable", "true"),
				),
			},
		},
	})
}

func testAccCheckVpcSubnetV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find %s in state", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Vpc Subnet data source ID not set")
		}

		return nil
	}
}

func testAccVpcSubnetV1DataSource_basic(rName, cidr, gateway string) string {
	return fmt.Sprintf(`
data "sbercloud_availability_zones" "test" {}

resource "sbercloud_vpc" "test" {
  name = "%s"
  cidr = "%s"
}

resource "sbercloud_vpc_subnet" "test" {
  name              = "%s"
  cidr              = "%s"
  gateway_ip        = "%s"
  vpc_id            = sbercloud_vpc.test.id
  availability_zone = data.sbercloud_availability_zones.test.names[0]
}

data "sbercloud_vpc_subnet" "by_id" {
  id = sbercloud_vpc_subnet.test.id
}

data "sbercloud_vpc_subnet" "by_cidr" {
  cidr = sbercloud_vpc_subnet.test.cidr
}

data "sbercloud_vpc_subnet" "by_name" {
  name = sbercloud_vpc_subnet.test.name
}

data "sbercloud_vpc_subnet" "by_vpc_id" {
  vpc_id = sbercloud_vpc_subnet.test.vpc_id
}
`, rName, cidr, rName, cidr, gateway)
}
