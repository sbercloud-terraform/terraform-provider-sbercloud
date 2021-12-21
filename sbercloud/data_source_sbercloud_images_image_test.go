package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSbercloudImagesV2ImageDataSource_basic(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSbercloudImagesV2ImageDataSource_ubuntu(rName),
			},
			{
				Config: testAccSbercloudImagesV2ImageDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.sbercloud_images_image.test"),
					resource.TestCheckResourceAttr(
						"data.sbercloud_images_image.test", "name", rName),
					resource.TestCheckResourceAttr(
						"data.sbercloud_images_image.test", "protected", "false"),
					resource.TestCheckResourceAttr(
						"data.sbercloud_images_image.test", "visibility", "private"),
				),
			},
		},
	})
}

func TestAccSbercloudImagesV2ImageDataSource_testQueries(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSbercloudImagesV2ImageDataSource_ubuntu(rName),
			},
			{
				Config: testAccSbercloudImagesV2ImageDataSource_querySizeMin(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.sbercloud_images_image.test"),
				),
			},
			{
				Config: testAccSbercloudImagesV2ImageDataSource_querySizeMax(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.sbercloud_images_image.test"),
				),
			},
			{
				Config: testAccSbercloudImagesV2ImageDataSource_ubuntu(rName),
			},
		},
	})
}

func testAccCheckImagesV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find image data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Image data source ID not set")
		}

		return nil
	}
}

func testAccSbercloudImagesV2ImageDataSource_ubuntu(rName string) string {
	return fmt.Sprintf(`
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

resource "sbercloud_compute_instance" "test" {
  name              = "%s"
  image_name        = "Ubuntu 18.04 server 64bit"
  flavor_id         = data.sbercloud_compute_flavors.test.ids[0]
  security_groups   = ["default"]
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  system_disk_type  = "SSD"

  network {
    uuid = data.sbercloud_vpc_subnet.test.id
  }
}

resource "sbercloud_images_image" "test" {
  name        = "%s"
  instance_id = sbercloud_compute_instance.test.id
  description = "created by TerraformAccTest"
}

`, rName, rName)
}

func testAccSbercloudImagesV2ImageDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_images_image" "test" {
	most_recent = true
	name = sbercloud_images_image.test.name
}
`, testAccSbercloudImagesV2ImageDataSource_ubuntu(rName))
}

func testAccSbercloudImagesV2ImageDataSource_querySizeMin(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_images_image" "test" {
	most_recent = true
	visibility = "private"
	size_min = "13000000"
}
`, testAccSbercloudImagesV2ImageDataSource_ubuntu(rName))
}

func testAccSbercloudImagesV2ImageDataSource_querySizeMax(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_images_image" "test" {
	most_recent = true
	visibility = "private"
	size_max = "23000000"
}
`, testAccSbercloudImagesV2ImageDataSource_ubuntu(rName))
}
