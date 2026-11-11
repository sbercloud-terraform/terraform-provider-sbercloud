package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccDataSourceEcsComputePassword_basic(t *testing.T) {
	dataSource := "data.sbercloud_compute_password.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEcsComputePassword_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
				),
			},
		},
	})
}

func testDataSourceEcsComputePassword_base(name string) string {
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

  os         = "Windows"
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
`, name)
}

func testDataSourceEcsComputePassword_basic(name string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_compute_password" "test" {
  server_id = sbercloud_compute_instance.test.id
}
`, testDataSourceEcsComputePassword_base(name))
}
