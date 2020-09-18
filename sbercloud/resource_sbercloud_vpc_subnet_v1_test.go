package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v1/subnets"
)

func TestAccVpcSubnetV1_basic(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetV1Exists("sbercloud_vpc_subnet_v1.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_subnet_v1.subnet_1", "name", "terraform_test_subnet"),
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_subnet_v1.subnet_1", "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_subnet_v1.subnet_1", "gateway_ip", "192.168.0.1"),
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_subnet_v1.subnet_1", "availability_zone", TEST_SBC_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_subnet_v1.subnet_1", "tags.foo", "bar"),
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_subnet_v1.subnet_1", "tags.key", "value"),
				),
			},
			{
				Config: testAccVpcSubnetV1_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_subnet_v1.subnet_1", "name", "terraform_test_subnet_1"),
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_subnet_v1.subnet_1", "tags.key", "value_updated"),
				),
			},
		},
	})
}

func TestAccVpcSubnetV1_timeout(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetV1Exists("sbercloud_vpc_subnet_v1.subnet_1", &subnet),
				),
			},
		},
	})
}

func testAccCheckVpcSubnetV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	subnetClient, err := config.networkingV1Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud vpc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_vpc_subnet_v1" {
			continue
		}

		_, err := subnets.Get(subnetClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Subnet still exists")
		}
	}

	return nil
}
func testAccCheckVpcSubnetV1Exists(n string, subnet *subnets.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		subnetClient, err := config.networkingV1Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud Vpc client: %s", err)
		}

		found, err := subnets.Get(subnetClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Subnet not found")
		}

		*subnet = *found

		return nil
	}
}

var testAccVpcSubnetV1_basic = fmt.Sprintf(`
resource "sbercloud_vpc_v1" "vpc_1" {
  name = "terraform_test_vpc_1"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_subnet_v1" "subnet_1" {
  name = "terraform_test_subnet"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${sbercloud_vpc_v1.vpc_1.id}"
  availability_zone = "%s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, TEST_SBC_AVAILABILITY_ZONE)
var testAccVpcSubnetV1_update = fmt.Sprintf(`
resource "sbercloud_vpc_v1" "vpc_1" {
  name = "terraform_test_vpc_1"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_subnet_v1" "subnet_1" {
  name = "terraform_test_subnet_1"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${sbercloud_vpc_v1.vpc_1.id}"
  availability_zone = "%s"

  tags = {
    foo = "bar"
    key = "value_updated"
  }
}
`, TEST_SBC_AVAILABILITY_ZONE)

var testAccVpcSubnetV1_timeout = fmt.Sprintf(`
resource "sbercloud_vpc_v1" "vpc_1" {
  name = "terraform_test_vpc_1"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_subnet_v1" "subnet_1" {
  name = "terraform_test_subnet_1"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${sbercloud_vpc_v1.vpc_1.id}"
  availability_zone = "%s"

 timeouts {
    create = "5m"
    delete = "5m"
  }

}
`, TEST_SBC_AVAILABILITY_ZONE)
