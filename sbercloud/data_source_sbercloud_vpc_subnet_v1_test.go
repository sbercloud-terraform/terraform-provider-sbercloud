package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVpcSubnetV1DataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcSubnetV1Config,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceVpcSubnetV1Check("data.sbercloud_vpc_subnet_v1.by_id", "terraform_test_subnet", "192.168.0.0/16",
						"192.168.0.1", TEST_SBC_AVAILABILITY_ZONE),
					testAccDataSourceVpcSubnetV1Check("data.sbercloud_vpc_subnet_v1.by_cidr", "terraform_test_subnet", "192.168.0.0/16",
						"192.168.0.1", TEST_SBC_AVAILABILITY_ZONE),
					testAccDataSourceVpcSubnetV1Check("data.sbercloud_vpc_subnet_v1.by_name", "terraform_test_subnet", "192.168.0.0/16",
						"192.168.0.1", TEST_SBC_AVAILABILITY_ZONE),
					testAccDataSourceVpcSubnetV1Check("data.sbercloud_vpc_subnet_v1.by_vpc_id", "terraform_test_subnet", "192.168.0.0/16",
						"192.168.0.1", TEST_SBC_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(
						"data.sbercloud_vpc_subnet_v1.by_id", "status", "ACTIVE"),
					resource.TestCheckResourceAttr(
						"data.sbercloud_vpc_subnet_v1.by_id", "dhcp_enable", "true"),
				),
			},
		},
	})
}

func testAccDataSourceVpcSubnetV1Check(n, name, cidr, gateway_ip, availability_zone string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", n)
		}

		subnetRs, ok := s.RootModule().Resources["sbercloud_vpc_subnet_v1.subnet_1"]
		if !ok {
			return fmt.Errorf("can't find sbercloud_vpc_subnet_v1.subnet_1 in state")
		}

		attr := rs.Primary.Attributes

		if attr["id"] != subnetRs.Primary.Attributes["id"] {
			return fmt.Errorf(
				"id is %s; want %s",
				attr["id"],
				subnetRs.Primary.Attributes["id"],
			)
		}

		if attr["cidr"] != cidr {
			return fmt.Errorf("bad subnet cidr %s, expected: %s", attr["cidr"], cidr)
		}
		if attr["name"] != name {
			return fmt.Errorf("bad subnet name %s", attr["name"])
		}
		if attr["gateway_ip"] != gateway_ip {
			return fmt.Errorf("bad subnet gateway_ip %s", attr["gateway_ip"])
		}
		if attr["availability_zone"] != availability_zone {
			return fmt.Errorf("bad subnet availability_zone %s", attr["availability_zone"])
		}

		return nil
	}
}

var testAccDataSourceVpcSubnetV1Config = fmt.Sprintf(`
resource "sbercloud_vpc_v1" "vpc_1" {
	name = "terraform_test_vpc"
	cidr= "192.168.0.0/16"
}

resource "sbercloud_vpc_subnet_v1" "subnet_1" {
	name = "terraform_test_subnet"
	cidr = "192.168.0.0/16"
	gateway_ip = "192.168.0.1"
	vpc_id = "${sbercloud_vpc_v1.vpc_1.id}"
	availability_zone = "%s"
 }

data "sbercloud_vpc_subnet_v1" "by_id" {
  id = "${sbercloud_vpc_subnet_v1.subnet_1.id}"
}

data "sbercloud_vpc_subnet_v1" "by_cidr" {
  cidr = "${sbercloud_vpc_subnet_v1.subnet_1.cidr}"
}

data "sbercloud_vpc_subnet_v1" "by_name" {
	name = "${sbercloud_vpc_subnet_v1.subnet_1.name}"
}

data "sbercloud_vpc_subnet_v1" "by_vpc_id" {
	vpc_id = "${sbercloud_vpc_subnet_v1.subnet_1.vpc_id}"
}
`, TEST_SBC_AVAILABILITY_ZONE)
