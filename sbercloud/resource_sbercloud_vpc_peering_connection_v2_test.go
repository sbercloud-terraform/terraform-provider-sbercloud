package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/peerings"
)

func TestAccVpcPeeringConnectionV2_basic(t *testing.T) {
	var peering peerings.Peering

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcPeeringConnectionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnectionV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcPeeringConnectionV2Exists("sbercloud_vpc_peering_connection_v2.peering_1", &peering),
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_peering_connection_v2.peering_1", "name", "terraform_test_peering"),
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_peering_connection_v2.peering_1", "status", "ACTIVE"),
				),
			},
			{
				Config: testAccVpcPeeringConnectionV2_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sbercloud_vpc_peering_connection_v2.peering_1", "name", "terraform_test_peering_1"),
				),
			},
		},
	})
}

func TestAccVpcPeeringConnectionV2_timeout(t *testing.T) {
	var peering peerings.Peering

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcPeeringConnectionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnectionV2_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcPeeringConnectionV2Exists("sbercloud_vpc_peering_connection_v2.peering_1", &peering),
				),
			},
		},
	})
}

func testAccCheckVpcPeeringConnectionV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	peeringClient, err := config.networkingHwV2Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud Peering client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_vpc_peering_connection_v2" {
			continue
		}

		_, err := peerings.Get(peeringClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Vpc Peering Connection still exists")
		}
	}

	return nil
}

func testAccCheckVpcPeeringConnectionV2Exists(n string, peering *peerings.Peering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		peeringClient, err := config.networkingHwV2Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud Peering client: %s", err)
		}

		found, err := peerings.Get(peeringClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Vpc peering Connection not found")
		}

		*peering = *found

		return nil
	}
}

const testAccVpcPeeringConnectionV2_basic = `
resource "sbercloud_vpc_v1" "vpc_1" {
  name = "terraform_test_vpc_1"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_v1" "vpc_2" {
  name = "terraform_test_vpc_2"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_peering_connection_v2" "peering_1" {
  name = "terraform_test_peering"
  vpc_id = "${sbercloud_vpc_v1.vpc_1.id}"
  peer_vpc_id = "${sbercloud_vpc_v1.vpc_2.id}"
}
`
const testAccVpcPeeringConnectionV2_update = `
resource "sbercloud_vpc_v1" "vpc_1" {
  name = "terraform_test_vpc_1"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_v1" "vpc_2" {
  name = "terraform_test_vpc_2"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_peering_connection_v2" "peering_1" {
  name = "terraform_test_peering_1"
  vpc_id = "${sbercloud_vpc_v1.vpc_1.id}"
  peer_vpc_id = "${sbercloud_vpc_v1.vpc_2.id}"
}
`
const testAccVpcPeeringConnectionV2_timeout = `
resource "sbercloud_vpc_v1" "vpc_1" {
  name = "terraform_test_vpc_1"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_v1" "vpc_2" {
  name = "terraform_test_vpc_2"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_peering_connection_v2" "peering_1" {
  name = "terraform_test_peering_1"
  vpc_id = "${sbercloud_vpc_v1.vpc_1.id}"
  peer_vpc_id = "${sbercloud_vpc_v1.vpc_2.id}"

 timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
