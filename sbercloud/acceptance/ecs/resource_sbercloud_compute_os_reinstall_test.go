package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccComputeOsReinstall_Basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeOsReinstall_basic(rName),
			},
		},
	})
}

func TestAccComputeOsReinstall_with_cloud_init(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeOsReinstall_with_cloud_init(rName),
			},
		},
	})
}

func testAccComputeOsReinstall_base(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_availability_zones" "test" {}

resource "sbercloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_subnet" "test" {
  vpc_id     = sbercloud_vpc.test.id
  name       = "%[1]s"
  cidr       = cidrsubnet(sbercloud_vpc.test.cidr, 4, 0)
  gateway_ip = cidrhost(cidrsubnet(sbercloud_vpc.test.cidr, 4, 0), 1)
}

resource "sbercloud_networking_secgroup" "test" {
  name = "%[1]s"
}

data "sbercloud_compute_flavors" "test" {
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "sbercloud_images_images" "test" {
  flavor_id = data.sbercloud_compute_flavors.test.ids[0]

  os         = "Ubuntu"
  visibility = "public"
}

resource "sbercloud_compute_instance" "test" {
  name               = "%[1]s"
  image_id           = data.sbercloud_images_images.test.images[0].id
  flavor_id          = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids = [sbercloud_networking_secgroup.test.id]
  availability_zone  = data.sbercloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = sbercloud_vpc_subnet.test.id
  }
}
`, rName)
}
func testAccComputeOsReinstall_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_compute_os_reinstall" "test" {
  cloud_init_installed = "false"
  server_id            = sbercloud_compute_instance.test.id

  os_reinstall {
    userid = "test"
    mode   = "withStopServer"

    metadata {
      __system__encrypted = "0"
    }
  }
}
`, testAccComputeOsReinstall_base(rName))
}

func testAccComputeOsReinstall_with_cloud_init(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_compute_os_reinstall" "test" {
  cloud_init_installed = "true"
  server_id            = sbercloud_compute_instance.test.id

  os_reinstall {
    userid = "test"
    mode   = "withStopServer"

    metadata {
      user_data           = "IyEvYmluL2Jhc2gKZWNobyB1c2VyX3Rlc3QgPiAvaG9tZS91c2VyLnR4dA=="
      __system__encrypted = "0"
    }
  }
}
`, testAccComputeOsReinstall_base(rName))
}
