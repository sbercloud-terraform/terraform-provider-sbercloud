// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file at
//     https://www.github.com/huaweicloud/magic-modules
//
// ----------------------------------------------------------------------------

package sbercloud

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccNatDnat_basic(t *testing.T) {
	randSuffix := acctest.RandString(5)
	resourceName := "sbercloud_nat_dnat_rule.dnat"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatDnatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatV2DnatRule_basic(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatDnatExists(),
					resource.TestCheckResourceAttr(resourceName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
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

func TestAccNatDnat_protocol(t *testing.T) {
	randSuffix := acctest.RandString(5)
	resourceName := "sbercloud_nat_dnat_rule.dnat"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatDnatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatV2DnatRule_protocol(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatDnatExists(),
					resource.TestCheckResourceAttr(resourceName, "protocol", "any"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
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

func testAccCheckNatDnatDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.NatGatewayClient(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_nat_dnat_rule" {
			continue
		}

		url, err := replaceVarsForTest(rs, "dnat_rules/{id}")
		if err != nil {
			return err
		}
		url = client.ServiceURL(url)

		_, err = client.Get(
			url, nil,
			&golangsdk.RequestOpts{MoreHeaders: map[string]string{"Accept": "application/json"}})
		if err == nil {
			return fmt.Errorf("sbercloud dnat rule still exists at %s", url)
		}
	}

	return nil
}

func testAccCheckNatDnatExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config.Config)
		client, err := config.NatGatewayClient(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["sbercloud_nat_dnat_rule.dnat"]
		if !ok {
			return fmt.Errorf("Error checking sbercloud_nat_dnat_rule.dnat exist, err=not found sbercloud_nat_dnat_rule.dnat")
		}

		url, err := replaceVarsForTest(rs, "dnat_rules/{id}")
		if err != nil {
			return fmt.Errorf("Error checking sbercloud_nat_dnat_rule.dnat exist, err=building url failed: %s", err)
		}
		url = client.ServiceURL(url)

		_, err = client.Get(
			url, nil,
			&golangsdk.RequestOpts{MoreHeaders: map[string]string{"Accept": "application/json"}})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmt.Errorf("sbercloud_nat_dnat_rule.dnat is not exist")
			}
			return fmt.Errorf("Error checking sbercloud_nat_dnat_rule.dnat exist, err=send request failed: %s", err)
		}
		return nil
	}
}

func replaceVarsForTest(rs *terraform.ResourceState, linkTmpl string) (string, error) {
	re := regexp.MustCompile("{([[:word:]]+)}")

	replaceFunc := func(s string) string {
		m := re.FindStringSubmatch(s)[1]
		if m == "project" {
			return "replace_holder"
		}
		if rs != nil {
			if m == "id" {
				return rs.Primary.ID
			}
			v, ok := rs.Primary.Attributes[m]
			if ok {
				return v
			}
		}
		return ""
	}

	s := re.ReplaceAllStringFunc(linkTmpl, replaceFunc)
	return strings.Replace(s, "replace_holder/", "", 1), nil
}

func testAccNatV2DnatRule_base(suffix string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc_eip" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

data "sbercloud_availability_zones" "test" {}

data "sbercloud_compute_flavors" "test" {
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "sbercloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "sbercloud_compute_instance" "instance_1" {
  name              = "instance-acc-test-%s"
  image_id          = data.sbercloud_images_image.test.id
  flavor_id         = data.sbercloud_compute_flavors.test.ids[0]
  security_groups   = ["default"]
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  system_disk_type  = "SSD"

  network {
    uuid = sbercloud_vpc_subnet.subnet_1.id
  }
  
  tags = {
    foo = "bar"
  }
}
`, suffix)
}

func testAccNatV2DnatRule_basic(suffix string) string {
	return fmt.Sprintf(`
%s

%s

resource "sbercloud_nat_dnat_rule" "dnat" {
  nat_gateway_id = sbercloud_nat_gateway.nat_1.id
  floating_ip_id = sbercloud_vpc_eip.eip_1.id
  private_ip     = sbercloud_compute_instance.instance_1.network.0.fixed_ip_v4
  protocol       = "tcp"
  internal_service_port = 993
  external_service_port = 242
}
`, testAccNatV2Gateway_basic(suffix), testAccNatV2DnatRule_base(suffix))
}

func testAccNatV2DnatRule_protocol(suffix string) string {
	return fmt.Sprintf(`
%s

%s

resource "sbercloud_nat_dnat_rule" "dnat" {
  nat_gateway_id = sbercloud_nat_gateway.nat_1.id
  floating_ip_id = sbercloud_vpc_eip.eip_1.id
  private_ip     = sbercloud_compute_instance.instance_1.network.0.fixed_ip_v4
  protocol       = "any"
  internal_service_port = 0
  external_service_port = 0
}
`, testAccNatV2Gateway_basic(suffix), testAccNatV2DnatRule_base(suffix))
}
