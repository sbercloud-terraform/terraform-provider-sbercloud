package cce

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNodesRemove_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()
	)

	baseConfig := testAccNodePool_base(name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccNodesRemove_basic(baseConfig, name),
				// there is nothing to check, if no error occurred, that means the test is successful
				ExpectError: regexp.MustCompile("Insufficient nodes blocks"),
			},
		},
	})
}

func testAccNodesRemove_base(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_cce_node_pool" "test" {
  cluster_id               = sbercloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "CentOS 7.6"
  flavor_id                = data.sbercloud_compute_flavors.test.ids[0]
  initial_node_count       = 2
  availability_zone        = data.sbercloud_availability_zones.test.names[0]
  key_pair                 = sbercloud_kps_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}
  `, baseConfig, name)
}

func testAccNodesRemove_basic(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

data "sbercloud_cce_nodes" "test" {
  cluster_id = sbercloud_cce_cluster.test.id

  depends_on = [ sbercloud_cce_node_pool.test ]
}

resource "sbercloud_cce_nodes_remove" "test" {
  cluster_id = sbercloud_cce_cluster.test.id

  dynamic "nodes" {
    for_each = data.sbercloud_cce_nodes.test.ids
    content {
      uid = nodes.value
    }
  }
}
`, testAccNodesRemove_base(baseConfig, name), name)
}
