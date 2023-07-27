package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccEvsVolumesDataSource_basic(t *testing.T) {
	dataSourceName := "data.sbercloud_evs_volumes.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolumesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "volumes.#", "5"),
				),
			},
		},
	})
}

func testAccEvsVolumesDataSource_base(rName string) string {
	return fmt.Sprintf(`
variable "volume_configuration" {
  type = list(object({
    suffix      = string
    size        = number
    device_type = string
    multiattach = bool
  }))
  default = [
    {suffix = "vbd_normal_volume", size = 100, device_type = "VBD", multiattach = false},
    {suffix = "vbd_share_volume", size = 100, device_type = "VBD", multiattach = true},
    {suffix = "scsi_normal_volume", size = 100, device_type = "SCSI", multiattach = false},
    {suffix = "scsi_share_volume", size = 100, device_type = "SCSI", multiattach = true},
  ]
}

%[1]s

resource "sbercloud_compute_instance" "test" {
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  name              = "%[2]s"
  image_id          = data.sbercloud_images_image.test.id
  flavor_id         = data.sbercloud_compute_flavors.test.ids[0]

  system_disk_type = "SSD"
  system_disk_size = 50

  security_group_ids = [
    sbercloud_networking_secgroup.test.id
  ]

  network {
    uuid = sbercloud_vpc_subnet.test.id
  }
}

resource "sbercloud_evs_volume" "test" {
  count = length(var.volume_configuration)
  
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  volume_type       = "SSD"
  name              = "%[2]s_${var.volume_configuration[count.index].suffix}"
  size              = var.volume_configuration[count.index].size
  device_type       = var.volume_configuration[count.index].device_type
  multiattach       = var.volume_configuration[count.index].multiattach

  tags = {
    index = tostring(count.index)
  }
}

resource "sbercloud_compute_volume_attach" "test" {
  count = length(sbercloud_evs_volume.test)

  instance_id = sbercloud_compute_instance.test.id
  volume_id   = sbercloud_evs_volume.test[count.index].id
}
`, acceptance.TestBaseComputeResources(rName), rName)
}

func testAccEvsVolumesDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_evs_volumes" "test" {
  depends_on = [sbercloud_compute_volume_attach.test]

  availability_zone = data.sbercloud_availability_zones.test.names[0]
  server_id         = sbercloud_compute_instance.test.id
  status            = "in-use"
}
`, testAccEvsVolumesDataSource_base(rName))
}
