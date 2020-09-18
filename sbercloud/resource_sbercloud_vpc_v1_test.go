package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v1/vpcs"
)

func TestAccVpcV1_basic(t *testing.T) {
	var vpc vpcs.Vpc

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists("sbercloud_vpc_v1.vpc_1", &vpc),
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_v1.vpc_1", "name", "terraform_test_vpc_1"),
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_v1.vpc_1", "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_v1.vpc_1", "status", "OK"),
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_v1.vpc_1", "shared", "false"),
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_v1.vpc_1", "tags.foo", "bar"),
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_v1.vpc_1", "tags.key", "value"),
				),
			},
		},
	})
}

func TestAccVpcV1_update(t *testing.T) {
	var vpc vpcs.Vpc

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists("sbercloud_vpc_v1.vpc_1", &vpc),
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_v1.vpc_1", "name", "terraform_test_vpc_1"),
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_v1.vpc_1", "tags.key", "value"),
				),
			},
			{
				Config: testAccVpcV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists("sbercloud_vpc_v1.vpc_1", &vpc),
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_v1.vpc_1", "name", "terraform_test_vpc_1_updated"),
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_v1.vpc_1", "tags.key", "value_updated"),
				),
			},
		},
	})
}

func TestAccVpcV1_timeout(t *testing.T) {
	var vpc vpcs.Vpc

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists("sbercloud_vpc_v1.vpc_1", &vpc),
				),
			},
		},
	})
}

func testAccCheckVpcV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vpcClient, err := config.networkingV1Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud vpc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_vpc_v1" {
			continue
		}

		_, err := vpcs.Get(vpcClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Vpc still exists")
		}
	}

	return nil
}

func testAccCheckVpcV1Exists(n string, vpc *vpcs.Vpc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		vpcClient, err := config.networkingV1Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud vpc client: %s", err)
		}

		found, err := vpcs.Get(vpcClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("vpc not found")
		}

		*vpc = *found

		return nil
	}
}

const testAccVpcV1_basic = `
resource "sbercloud_vpc_v1" "vpc_1" {
  name = "terraform_test_vpc_1"
  cidr="192.168.0.0/16"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`

const testAccVpcV1_update = `
resource "sbercloud_vpc_v1" "vpc_1" {
  name = "terraform_test_vpc_1_updated"
  cidr="192.168.0.0/16"

  tags = {
    foo = "bar"
    key = "value_updated"
  }
}
`
const testAccVpcV1_timeout = `
resource "sbercloud_vpc_v1" "vpc_1" {
  name = "terraform_test_vpc_1"
  cidr="192.168.0.0/16"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
