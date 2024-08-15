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

func TestAccComputeInstancesDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.sbercloud_compute_instances.test"
	var instance cloudservers.CloudServer

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstancesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("sbercloud_compute_instance.test", &instance),
					testAccCheckComputeInstancesDataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "instances.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.image_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.flavor_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.flavor_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.availability_zone"),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.security_group_ids.#", "1"),
				),
			},
		},
	})
}

func testAccCheckComputeInstancesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find compute instances data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Data source ID not set")
		}

		return nil
	}
}

func testAccComputeInstancesDataSource_basic(rName string) string {
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

  tags = {
    foo = "bar"
  }
}

data "sbercloud_compute_instances" "test" {
  name = sbercloud_compute_instance.test.name

  depends_on = [
    sbercloud_compute_instance.test
  ]
}
`, testAccCompute_data, rName)
}
