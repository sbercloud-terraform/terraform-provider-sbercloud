package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccComputeOsChange_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckECSID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeOsChange_basic(),
			},
		},
	})
}

func TestAccComputeOsChange_with_cloud_init(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeOsChange_with_cloud_init(),
			},
		},
	})
}

func testAccComputeOsChange_basic() string {
	return fmt.Sprintf(`
data "sbercloud_images_images" "test" {
  os         = "Ubuntu"
  visibility = "public"
}

resource "sbercloud_compute_os_change" "test" {
  cloud_init_installed = "false"
  server_id            = "%s"

  os_change {
    imageid = data.sbercloud_images_images.test.images[0].id
    userid  = "test"
    mode    = "withStopServer"

    metadata {
      __system__encrypted = "0"
    }
  }
}
`, acceptance.SBC_ECS_ID)
}

func testAccComputeOsChange_with_cloud_init() string {
	return fmt.Sprintf(`
data "sbercloud_images_images" "test" {
  os         = "Ubuntu"
  visibility = "public"
}

resource "sbercloud_compute_os_change" "test" {
  cloud_init_installed = "true"
  server_id            = "%s"

  os_change {
    imageid = data.sbercloud_images_images.test.images[0].id
    userid  = "test"
    mode    = "withStopServer"

    metadata {
      user_data           = "IyEvYmluL2Jhc2gKZWNobyB1c2VyX3Rlc3QgPiAvaG9tZS91c2VyLnR4dA=="
      __system__encrypted = "0"
    }
  }
}
`, acceptance.SBC_ECS_ID)
}
