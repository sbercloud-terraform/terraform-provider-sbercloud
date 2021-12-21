package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/compute/v2/extensions/servergroups"
	"github.com/chnsz/golangsdk/openstack/compute/v2/servers"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccComputeV2ServerGroup_basic(t *testing.T) {
	var sg servergroups.ServerGroup
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2ServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2ServerGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2ServerGroupExists("sbercloud_compute_servergroup.sg_1", &sg),
				),
			},
			{
				ResourceName:      "sbercloud_compute_servergroup.sg_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeV2ServerGroup_affinity(t *testing.T) {
	var instance servers.Server
	var sg servergroups.ServerGroup
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2ServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2ServerGroup_affinity(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2ServerGroupExists("sbercloud_compute_servergroup.sg_1", &sg),
					testAccCheckComputeV2InstanceExists("sbercloud_compute_instance.instance_1", &instance),
					testAccCheckComputeV2InstanceInServerGroup(&instance, &sg),
				),
			},
		},
	})
}

func testAccCheckComputeV2ServerGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	computeClient, err := config.ComputeV2Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Sbercloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_compute_servergroup" {
			continue
		}

		_, err := servergroups.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("ServerGroup still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2ServerGroupExists(n string, kp *servergroups.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		computeClient, err := config.ComputeV2Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating Sbercloud compute client: %s", err)
		}

		found, err := servergroups.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("ServerGroup not found")
		}

		*kp = *found

		return nil
	}
}

func testAccCheckComputeV2InstanceInServerGroup(instance *servers.Server, sg *servergroups.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(sg.Members) > 0 {
			for _, m := range sg.Members {
				if m == instance.ID {
					return nil
				}
			}
		}

		return fmt.Errorf("Instance %s is not part of Server Group %s", instance.ID, sg.ID)
	}
}

func testAccComputeV2ServerGroup_basic(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_compute_servergroup" "sg_1" {
  name = "%s"
  policies = ["affinity"]
}
`, rName)
}

func testAccComputeV2ServerGroup_affinity(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_compute_servergroup" "sg_1" {
  name = "%s"
  policies = ["affinity"]
}

resource "sbercloud_compute_instance" "instance_1" {
  name = "%s"
  image_id = data.sbercloud_images_image.test.id
  flavor_id = data.sbercloud_compute_flavors.test.ids[0]
  security_groups = ["default"]
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  system_disk_type  = "SSD"
  scheduler_hints {
    group = sbercloud_compute_servergroup.sg_1.id
  }
  network {
    uuid = data.sbercloud_vpc_subnet.test.id
  }
}
`, testAccCompute_data, rName, rName)
}
