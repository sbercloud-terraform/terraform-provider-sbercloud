package as

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/groups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccASV1Group_basic(t *testing.T) {
	var asGroup groups.Group
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_as_group.hth_as_group"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASV1GroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASV1Group_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1GroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "lbaas_listeners.0.protocol_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
		},
	})
}

func testAccCheckASV1GroupDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	asClient, err := config.AutoscalingV1Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating sbercloud autoscaling client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_as_group" {
			continue
		}

		_, err := groups.Get(asClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("AS group still exists")
		}
	}

	log.Printf("[DEBUG] testCheckASV1GroupDestroy success!")

	return nil
}

func testAccCheckASV1GroupExists(n string, group *groups.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		asClient, err := config.AutoscalingV1Client(acceptance.SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating sbercloud autoscaling client: %s", err)
		}

		found, err := groups.Get(asClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Autoscaling Group not found")
		}
		log.Printf("[DEBUG] test found is: %#v", found)
		group = &found

		return nil
	}
}

func testASV1Group_basic(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_availability_zones" "test" {}

data "sbercloud_vpc" "test" {
  name = "vpc-default"
}

data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

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

resource "sbercloud_networking_secgroup" "secgroup" {
  name        = "%s"
  description = "This is a terraform test security group"
}

resource "sbercloud_compute_keypair" "hth_key" {
  name       = "%s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "sbercloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "%s"
  vip_subnet_id = data.sbercloud_vpc_subnet.test.subnet_id
}

resource "sbercloud_lb_listener" "listener_1" {
  name            = "%s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = sbercloud_lb_loadbalancer.loadbalancer_1.id
}

resource "sbercloud_lb_pool" "pool_1" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = sbercloud_lb_listener.listener_1.id
}

resource "sbercloud_as_configuration" "hth_as_config"{
  scaling_configuration_name = "%s"
  instance_config {
	image    = data.sbercloud_images_image.test.id
	flavor   = data.sbercloud_compute_flavors.test.ids[0]
    key_name = sbercloud_compute_keypair.hth_key.id
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }
  }
}

resource "sbercloud_as_group" "hth_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = sbercloud_as_configuration.hth_as_config.id
  vpc_id                   = data.sbercloud_vpc.test.id

  networks {
    id = data.sbercloud_vpc_subnet.test.id
  }
  security_groups {
    id = sbercloud_networking_secgroup.secgroup.id
  }
  lbaas_listeners {
    pool_id       = sbercloud_lb_pool.pool_1.id
    protocol_port = sbercloud_lb_listener.listener_1.protocol_port
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName, rName, rName, rName, rName, rName, rName)
}
