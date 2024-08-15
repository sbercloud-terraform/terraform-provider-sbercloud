package cbr

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataBackup_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceNameWithDash()
	dataSourceName := "data.sbercloud_cbr_backup.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataBackup_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccDataBackup_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.sbercloud_images_image.test.id
  flavor_id          = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids = [sbercloud_networking_secgroup.test.id]
  availability_zone  = data.sbercloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = sbercloud_vpc_subnet.test.id
  }

  data_disks {
    type = "SAS"
    size = "10"
  }
}

resource "sbercloud_cbr_vault" "test" {
  name             = "%[2]s"
  type             = "server"
  consistent_level = "app_consistent"
  protection_type  = "backup"
  size             = 200
}

resource "sbercloud_images_image" "test" {
  name        = "%[2]s"
  instance_id = sbercloud_compute_instance.test.id
  vault_id    = sbercloud_cbr_vault.test.id
}

data "sbercloud_cbr_backup" "test" {
  id = sbercloud_images_image.test.backup_id
}
`, acceptance.TestBaseComputeResources(name), name)
}
