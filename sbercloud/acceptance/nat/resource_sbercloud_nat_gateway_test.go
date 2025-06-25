package sbercloud

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/natgateways"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccNatGateway_basic(t *testing.T) {
	randSuffix := acctest.RandString(5)
	resourceName := "sbercloud_nat_gateway.nat_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNatV2GatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatV2Gateway_basic(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatV2GatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("nat-gateway-basic-%s", randSuffix)),
					resource.TestCheckResourceAttr(resourceName, "description", "test for terraform"),
					resource.TestCheckResourceAttr(resourceName, "spec", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNatV2Gateway_update(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("nat-gateway-updated-%s", randSuffix)),
					resource.TestCheckResourceAttr(resourceName, "description", "test for terraform updated"),
					resource.TestCheckResourceAttr(resourceName, "spec", "2"),
				),
			},
		},
	})
}

func TestAccNatGateway_withEpsId(t *testing.T) {
	randSuffix := acctest.RandString(5)
	resourceName := "sbercloud_nat_gateway.nat_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNatV2GatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatV2Gateway_epsId(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatV2GatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.SBC_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckNatV2GatewayDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	natClient, err := config.NatGatewayClient(acceptance.SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud nat client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_nat_gateway" {
			continue
		}

		_, err := natgateways.Get(natClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Nat gateway still exists")
		}
	}

	return nil
}

func testAccCheckNatV2GatewayExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		natClient, err := config.NatGatewayClient(acceptance.SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud nat client: %s", err)
		}

		found, err := natgateways.Get(natClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Nat gateway not found")
		}

		return nil
	}
}

func testAccNatPreCondition(suffix string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc" "vpc_1" {
  name = "nat-vpc-%s"
  cidr = "172.16.0.0/16"
}

resource "sbercloud_vpc_subnet" "subnet_1" {
  name       = "nat-sunnet-%s"
  cidr       = "172.16.10.0/24"
  gateway_ip = "172.16.10.1"
  vpc_id     = sbercloud_vpc.vpc_1.id
  dns_list   = ["100.125.13.59"]
}
	`, suffix, suffix)
}

func testAccNatV2Gateway_basic(suffix string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_nat_gateway" "nat_1" {
  name        = "nat-gateway-basic-%s"
  description = "test for terraform"
  spec        = "1"
  vpc_id      = sbercloud_vpc.vpc_1.id
  subnet_id   = sbercloud_vpc_subnet.subnet_1.id
}
	`, testAccNatPreCondition(suffix), suffix)
}

func testAccNatV2Gateway_update(suffix string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_nat_gateway" "nat_1" {
  name        = "nat-gateway-updated-%s"
  description = "test for terraform updated"
  spec        = "2"
  vpc_id      = sbercloud_vpc.vpc_1.id
  subnet_id   = sbercloud_vpc_subnet.subnet_1.id
}
	`, testAccNatPreCondition(suffix), suffix)
}

func testAccNatV2Gateway_epsId(suffix string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_nat_gateway" "nat_1" {
  name                  = "nat-gateway-eps-%s"
  description           = "test for terraform"
  spec                  = "1"
  vpc_id                = sbercloud_vpc.vpc_1.id
  subnet_id             = sbercloud_vpc_subnet.subnet_1.id
  enterprise_project_id = "%s"
}
	`, testAccNatPreCondition(suffix), suffix, acceptance.SBC_ENTERPRISE_PROJECT_ID_TEST)
}
