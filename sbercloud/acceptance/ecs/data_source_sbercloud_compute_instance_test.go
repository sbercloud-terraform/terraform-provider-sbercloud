package ecs

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComputeInstanceDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("ecs-data-test-%s", acctest.RandString(5))
	resourceName := "data.sbercloud_compute_instance.this"
	var instance cloudservers.CloudServer

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("sbercloud_compute_instance.test", &instance),
					testAccCheckComputeInstanceDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "system_disk_id"),
					resource.TestCheckResourceAttrSet(resourceName, "security_groups.#"),
					resource.TestCheckResourceAttrSet(resourceName, "network.#"),
					resource.TestCheckResourceAttrSet(resourceName, "volume_attached.#"),
				),
			},
		},
	})
}

func testAccCheckComputeInstanceDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find compute instance data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Compute instance data source ID not set")
		}

		return nil
	}
}

func testAccComputeInstanceDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.sbercloud_images_image.test.id
  flavor_id          = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids = [data.sbercloud_networking_secgroup.test.id]
  availability_zone  = data.sbercloud_availability_zones.test.names[0]
  system_disk_type  = "SSD"

  network {
    uuid = data.sbercloud_vpc_subnet.test.id
  }
}

data "sbercloud_compute_instance" "this" {
  name = sbercloud_compute_instance.test.name

  depends_on = [
    sbercloud_compute_instance.test
  ]
}
`, testAccCompute_data, rName)
}
