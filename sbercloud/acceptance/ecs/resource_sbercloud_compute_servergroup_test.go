package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/servergroups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccComputeServerGroup_basic(t *testing.T) {
	var sg servergroups.ServerGroup
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_compute_servergroup.sg_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeServerGroupExists(resourceName, &sg),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeServerGroup_scheduler(t *testing.T) {
	var instance cloudservers.CloudServer
	var sg servergroups.ServerGroup
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_compute_servergroup.sg_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroup_scheduler(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeServerGroupExists(resourceName, &sg),
					testAccCheckComputeInstanceExists("sbercloud_compute_instance.instance_1", &instance),
					testAccCheckComputeInstanceInServerGroup(&instance, &sg),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func TestAccComputeServerGroup_members(t *testing.T) {
	var instance cloudservers.CloudServer
	var sg servergroups.ServerGroup
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_compute_servergroup.sg_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroup_members(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeServerGroupExists(resourceName, &sg),
					testAccCheckComputeInstanceExists("sbercloud_compute_instance.instance_1", &instance),
					testAccCheckComputeInstanceInServerGroup(&instance, &sg),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func testAccCheckComputeServerGroupDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	ecsClient, err := cfg.ComputeV1Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_compute_servergroup" {
			continue
		}

		_, err := servergroups.Get(ecsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("server group still exists")
		}
	}

	return nil
}

func testAccCheckComputeServerGroupExists(n string, kp *servergroups.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		ecsClient, err := cfg.ComputeV1Client(acceptance.SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating compute client: %s", err)
		}

		found, err := servergroups.Get(ecsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("server group not found")
		}

		*kp = *found

		return nil
	}
}

func testAccCheckComputeInstanceInServerGroup(instance *cloudservers.CloudServer, sg *servergroups.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(sg.Members) > 0 {
			for _, m := range sg.Members {
				if m == instance.ID {
					return nil
				}
			}
		}

		return fmt.Errorf("instance %s does not belong to server group %s", instance.ID, sg.ID)
	}
}

func testAccComputeServerGroup_basic(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_compute_servergroup" "sg_1" {
  name     = "%s"
  policies = ["anti-affinity"]
}
`, rName)
}

func testAccComputeServerGroup_scheduler(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_compute_servergroup" "sg_1" {
  name     = "%s"
  policies = ["anti-affinity"]
}

resource "sbercloud_compute_instance" "instance_1" {
  name               = "%s"
  image_id           = data.sbercloud_images_image.test.id
  flavor_id          = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids = [data.sbercloud_networking_secgroup.test.id]
  availability_zone  = data.sbercloud_availability_zones.test.names[0]

  scheduler_hints {
    group = sbercloud_compute_servergroup.sg_1.id
  }
  network {
    uuid = data.sbercloud_vpc_subnet.test.id
  }
}
`, testAccCompute_data, rName, rName)
}

func testAccComputeServerGroup_members(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_compute_servergroup" "sg_1" {
  name     = "%s"
  policies = ["anti-affinity"]
  members  = [sbercloud_compute_instance.instance_1.id]
}

resource "sbercloud_compute_instance" "instance_1" {
  name               = "%s"
  image_id           = data.sbercloud_images_image.test.id
  flavor_id          = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids = [data.sbercloud_networking_secgroup.test.id]
  availability_zone  = data.sbercloud_availability_zones.test.names[0]

  network {
    uuid = data.sbercloud_vpc_subnet.test.id
  }
}
`, testAccCompute_data, rName, rName)
}
