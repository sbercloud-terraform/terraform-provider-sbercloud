package vpc

import (
	"fmt"
	"log"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccNetworkingV2VIPAssociate_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2VIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2VIPAssociateConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("sbercloud_networking_vip_associate.vip_associate",
						"vip_id", "sbercloud_networking_vip.vip", "id"),
					resource.TestCheckOutput("port_ids_check", "true"),
				),
			},
			{
				Config:            testAccNetworkingV2VIPAssociateConfig_config(rName),
				ResourceName:      "sbercloud_networking_vip_associate.vip_associate",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNetworkingV2VIPAssociateImportStateIdFunc(),
			},
			{
				Config: testAccNetworkingV2VIPAssociateConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("sbercloud_networking_vip_associate.vip_associate",
						"vip_id", "sbercloud_networking_vip.vip", "id"),
					resource.TestCheckOutput("port_ids_check", "true"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2VIPAssociateDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV2Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_networking_vip_associate" {
			continue
		}

		vipID := rs.Primary.Attributes["vip_id"]
		_, err = ports.Get(networkingClient, vipID).Extract()
		if err != nil {
			// If the error is a 404, then the vip port does not exist,
			// and therefore the floating IP cannot be associated to it.
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}
	}

	log.Printf("[DEBUG] Destroy NetworkingVIPAssociated success!")
	return nil
}

func testAccNetworkingV2VIPAssociateImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		vip, ok := s.RootModule().Resources["sbercloud_networking_vip.vip"]
		if !ok {
			return "", fmt.Errorf("vip not found: %s", vip)
		}
		instance, ok := s.RootModule().Resources["sbercloud_compute_instance.test.0"]
		if !ok {
			return "", fmt.Errorf("port not found: %s", instance)
		}
		if vip.Primary.ID == "" || instance.Primary.Attributes["network.0.port"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s", vip.Primary.ID,
				instance.Primary.Attributes["network.0.port"])
		}
		return fmt.Sprintf("%s/%s", vip.Primary.ID, instance.Primary.Attributes["network.0.port"]), nil
	}
}

const testAccCompute_data = `
data "sbercloud_availability_zones" "test" {}

data "sbercloud_compute_flavors" "test" {
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "sbercloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "sbercloud_networking_secgroup" "test" {
  name = "default"
}
`

func testAccNetworkingV2VIPAssociateConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_compute_instance" "test" {
  count               = 1
  name                = "%s-${count.index}"
  image_id            = data.sbercloud_images_image.test.id
  flavor_id           = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids  = [data.sbercloud_networking_secgroup.test.id]
  stop_before_destroy = true
  system_disk_type    = "SSD"

  network {
    uuid              = data.sbercloud_vpc_subnet.test.id
    source_dest_check = false
  }
}

resource "sbercloud_networking_vip" "vip" {
  network_id = data.sbercloud_vpc_subnet.test.id
}

resource "sbercloud_networking_vip_associate" "vip_associate" {
  vip_id   = sbercloud_networking_vip.vip.id
  port_ids = [sbercloud_compute_instance.test[0].network[0].port]
}

locals {
  port_ids_result = [
    for v in sbercloud_compute_instance.test[*].network[0].port : contains(sbercloud_networking_vip_associate.vip_associate.port_ids, v)]
}

output "port_ids_check" {
  value = alltrue(local.port_ids_result)
}
`, testAccCompute_data, rName)
}

func testAccNetworkingV2VIPAssociateConfig_config(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_compute_instance" "test" {
  count               = 1
  name                = "%s-${count.index}"
  image_id            = data.sbercloud_images_image.test.id
  flavor_id           = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids  = [data.sbercloud_networking_secgroup.test.id]
  stop_before_destroy = true
  system_disk_type    = "SSD"

  network {
    uuid              = data.sbercloud_vpc_subnet.test.id
    source_dest_check = false
  }
}

resource "sbercloud_networking_vip" "vip" {
  network_id = data.sbercloud_vpc_subnet.test.id
}

resource "sbercloud_networking_vip_associate" "vip_associate" {
  vip_id   = sbercloud_networking_vip.vip.id
  port_ids = [sbercloud_compute_instance.test[0].network[0].port]
}
`, testAccCompute_data, rName)
}

func testAccNetworkingV2VIPAssociateConfig_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_compute_instance" "test" {
  count               = 2
  name                = "%s-${count.index}"
  image_id            = data.sbercloud_images_image.test.id
  flavor_id           = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids  = [data.sbercloud_networking_secgroup.test.id]
  stop_before_destroy = true
  system_disk_type    = "SSD"

  network {
    uuid              = data.sbercloud_vpc_subnet.test.id
    source_dest_check = false
  }
}

resource "sbercloud_networking_vip" "vip" {
  network_id = data.sbercloud_vpc_subnet.test.id
}

resource "sbercloud_networking_vip_associate" "vip_associate" {
  vip_id   = sbercloud_networking_vip.vip.id
  port_ids = [
    sbercloud_compute_instance.test[0].network[0].port,
    sbercloud_compute_instance.test[1].network[0].port,
  ]
}

locals {
  port_ids_result = [
    for v in sbercloud_compute_instance.test[*].network[0].port : contains(sbercloud_networking_vip_associate.vip_associate.port_ids, v)]
}

output "port_ids_check" {
  value = alltrue(local.port_ids_result)
}
`, testAccCompute_data, rName)
}
