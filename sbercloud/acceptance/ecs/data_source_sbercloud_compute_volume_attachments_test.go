package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccDataSourceEcsComputeVolumeAttachments_basic(t *testing.T) {
	dataSource := "data.sbercloud_compute_volume_attachments.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEcsComputeVolumeAttachments_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "volume_attachments.#"),
					resource.TestCheckResourceAttrSet(dataSource, "volume_attachments.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "volume_attachments.0.device"),
					resource.TestCheckResourceAttrSet(dataSource, "volume_attachments.0.server_id"),
					resource.TestCheckResourceAttrSet(dataSource, "volume_attachments.0.volume_id"),
				),
			},
		},
	})
}

func testDataSourceEcsComputeVolumeAttachments_base(name string) string {
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

resource "sbercloud_evs_volume" "test" {
  count = 2
  name              = "%[1]s_${count.index}"
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 10
}

resource "sbercloud_compute_volume_attach" "test" {
  count = 2

  instance_id           = sbercloud_compute_instance.test.id
  volume_id             = sbercloud_evs_volume.test[count.index].id
  delete_on_termination = "true"
}
`, name)
}

func testDataSourceEcsComputeVolumeAttachments_basic(name string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_compute_volume_attachments" "test" {
  depends_on = [
    sbercloud_compute_volume_attach.test[0],
    sbercloud_compute_volume_attach.test[1]
  ]
  server_id = sbercloud_compute_instance.test.id
}
`, testDataSourceEcsComputeVolumeAttachments_base(name))
}
