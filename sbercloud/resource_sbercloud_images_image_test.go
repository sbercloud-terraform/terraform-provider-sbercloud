package sbercloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/ims/v2/cloudimages"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccImsImage_basic(t *testing.T) {
	var image cloudimages.Image

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameUpdate := rName + "-update"
	resourceName := "sbercloud_images_image.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckImsImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImage_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageExists(resourceName, &image),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				Config: testAccImsImage_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageExists(resourceName, &image),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
				),
			},
		},
	})
}

func getCloudimage(client *golangsdk.ServiceClient, id string) (*cloudimages.Image, error) {
	listOpts := &cloudimages.ListOpts{
		ID:    id,
		Limit: 1,
	}
	allPages, err := cloudimages.List(client, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("Unable to query images: %s", err)
	}

	allImages, err := cloudimages.ExtractImages(allPages)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve images: %s", err)
	}

	if len(allImages) < 1 {
		return nil, fmt.Errorf("Unable to find images %s: Maybe not existed", id)
	}

	img := allImages[0]
	if img.ID != id {
		return nil, fmt.Errorf("Unexpected images ID")
	}
	log.Printf("[DEBUG] Retrieved Image %s: %#v", id, img)
	return &img, nil
}

func testAccCheckImsImageDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*huaweicloud.Config)
	imageClient, err := config.ImageV2Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Sbercloud Image: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_images_image" {
			continue
		}

		_, err := getCloudimage(imageClient, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Image still exists")
		}
	}

	return nil
}

func testAccCheckImsImageExists(n string, image *cloudimages.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("IMS Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*huaweicloud.Config)
		imageClient, err := config.ImageV2Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating Sbercloud Image: %s", err)
		}

		found, err := getCloudimage(imageClient, rs.Primary.ID)
		if err != nil {
			return err
		}

		*image = *found
		return nil
	}
}

func testAccImsImage_basic(rName string) string {
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

func testAccImsImage_update(rName string) string {
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
