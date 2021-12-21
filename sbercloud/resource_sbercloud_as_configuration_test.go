package sbercloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/configurations"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccASV1Configuration_basic(t *testing.T) {
	var asConfig configurations.Configuration
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckASV1ConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASV1Configuration_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1ConfigurationExists("sbercloud_as_configuration.hth_as_config", &asConfig),
				),
			},
		},
	})
}

func testAccCheckASV1ConfigurationDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	asClient, err := config.AutoscalingV1Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating sbercloud autoscaling client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_as_configuration" {
			continue
		}

		_, err := configurations.Get(asClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("AS configuration still exists")
		}
	}

	log.Printf("[DEBUG] testAccCheckASV1ConfigurationDestroy success!")

	return nil
}

func testAccCheckASV1ConfigurationExists(n string, configuration *configurations.Configuration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		asClient, err := config.AutoscalingV1Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating sbercloud autoscaling client: %s", err)
		}

		found, err := configurations.Get(asClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Autoscaling Configuration not found")
		}
		log.Printf("[DEBUG] test found is: %#v", found)
		configuration = &found

		return nil
	}
}

func testAccASV1Configuration_basic(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_availability_zones" "test" {}

data "sbercloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "sbercloud_compute_flavors" "test" {
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "sbercloud_compute_keypair" "hth_key" {
  name = "%s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "sbercloud_as_configuration" "hth_as_config"{
  scaling_configuration_name = "%s"
  instance_config {
	image = data.sbercloud_images_image.test.id
	flavor = data.sbercloud_compute_flavors.test.ids[0]
    disk {
      size = 40
      volume_type = "SATA"
      disk_type = "SYS"
    }
    key_name = "${sbercloud_compute_keypair.hth_key.id}"
  }
}
`, rName, rName)
}
