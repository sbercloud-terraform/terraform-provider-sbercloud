package cce

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccClusterUpgrade_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterUpgrade_basic(name),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testAccClusterUpgrade_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "sbercloud_availability_zones" "test" {}

data "sbercloud_compute_flavors" "test" {
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "sbercloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  cluster_version        = "v1.28"
  vpc_id                 = sbercloud_vpc.test.id
  subnet_id              = sbercloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  container_network_cidr = "172.16.0.0/24"
  service_network_cidr   = "172.17.0.0/16"

  lifecycle {
    ignore_changes = [
      cluster_version
    ]
  }
}

resource "sbercloud_cce_node" "test" {
  count = 2

  cluster_id        = sbercloud_cce_cluster.test.id
  name              = "%[2]s-${count.index}"
  flavor_id         = data.sbercloud_compute_flavors.test.ids[0]
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  password          = "Test@1234"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }

  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}

resource "sbercloud_cce_node_pool" "test" {
  cluster_id               = sbercloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "CentOS 7.6"
  flavor_id                = data.sbercloud_compute_flavors.test.ids[0]
  initial_node_count       = 2
  availability_zone        = data.sbercloud_availability_zones.test.names[0]
  password                 = "Test@1234"
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

`, acceptance.TestVpc(name), name)
}

func testAccClusterUpgrade_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_cce_cluster_upgrade" "test" {
  cluster_id     = sbercloud_cce_cluster.test.id
  target_version = "v1.32"

  strategy {
    type = "inPlaceRollingUpdate"
    in_place_rolling_update {
      user_defined_step = 20
    }
  }

  node_order = {
    "DefaultPool" = jsonencode(
      [
        {
          "nodeSelector" : {
            "key" : "name",
            "value" : [
              "${sbercloud_cce_node.test[0].name}"
            ],
            "operator" : "="
          },
          "priority" : 1
        },
        {
          "nodeSelector" : {
            "key" : "name",
            "value" : [
              "${sbercloud_cce_node.test[1].name}"
            ],
            "operator" : "="
          },
          "priority" : 2
        }
      ]
    )
  }

  nodepool_order = {
    "DefaultPool"                          = 1
    "${sbercloud_cce_node_pool.test.id}" = 2
  }
}
`, testAccClusterUpgrade_base(name), name)
}
