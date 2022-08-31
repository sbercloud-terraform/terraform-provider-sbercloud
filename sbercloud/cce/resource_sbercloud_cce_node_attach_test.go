package cce

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/cce/v3/nodes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccCCENodeAttachV3_basic(t *testing.T) {
	var node nodes.Nodes

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameUpdate := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_cce_node_attach.test"
	//clusterName here is used to provide the cluster id to fetch cce node.
	clusterName := "sbercloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeAttachV3_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "os", "CentOS 7.6"),
				),
			},
			{
				Config: testAccCCENodeAttachV3_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "os", "CentOS 7.6"),
				),
			},
		},
	})
}

func testAccCCENodeAttachV3_Base(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_availability_zones" "test" {}

data "sbercloud_images_image" "test" {
  name = "CentOS 7.6 64bit"
  most_recent = true
}

data "sbercloud_compute_flavors" "test" {
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "sbercloud_compute_keypair" "test" {
  name = "%s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "sbercloud_compute_instance" "test" {
  name                        = "%s"
  image_id                    = data.sbercloud_images_image.test.id
  flavor_id                   = data.sbercloud_compute_flavors.test.ids[0]
  availability_zone           = data.sbercloud_availability_zones.test.names[0]
  key_pair                    = sbercloud_compute_keypair.test.name
  delete_disks_on_termination = true

  system_disk_type = "SAS"
  system_disk_size = 50

  data_disks {
	type = "SAS"
	size = "100"
  }

  network {
	uuid = sbercloud_vpc_subnet.test.id
  }

  lifecycle {
    ignore_changes = [
      image_id, tags, name
    ]
  }
}

resource "sbercloud_cce_cluster" "test" {
  name                   = "%s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = sbercloud_vpc.test.id
  subnet_id              = sbercloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
}
`, testAccCCEClusterV3_Base(rName), rName, rName, rName)
}

func testAccCCENodeAttachV3_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_cce_node_attach" "test" {
  cluster_id = sbercloud_cce_cluster.test.id
  server_id  = sbercloud_compute_instance.test.id
  key_pair   = sbercloud_compute_keypair.test.name
  os         = "CentOS 7.6"
  name       = "%s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCCENodeAttachV3_Base(rName), rName)
}

func testAccCCENodeAttachV3_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_cce_node_attach" "test" {
  cluster_id = sbercloud_cce_cluster.test.id
  server_id  = sbercloud_compute_instance.test.id
  key_pair   = sbercloud_compute_keypair.test.name
  os         = "CentOS 7.6"
  name       = "%s"

  tags = {
    foo        = "bar_update"
    key_update = "value_update"
  }
}
`, testAccCCENodeAttachV3_Base(rName), rNameUpdate)
}
