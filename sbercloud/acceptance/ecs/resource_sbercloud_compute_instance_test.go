package ecs

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccComputeInstance_basic(t *testing.T) {
	var instance cloudservers.CloudServer

	rName := acceptance.RandomAccResourceName()
	resourceName := "sbercloud_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test"),
					resource.TestCheckResourceAttr(resourceName, "hostname", "hostname-test"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(resourceName, "system_disk_id"),
					resource.TestCheckResourceAttrSet(resourceName, "security_groups.#"),
					resource.TestCheckResourceAttrSet(resourceName, "volume_attached.#"),
					resource.TestCheckResourceAttrSet(resourceName, "network.#"),
					resource.TestCheckResourceAttrSet(resourceName, "network.0.port"),
					resource.TestCheckResourceAttrSet(resourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "network.0.source_dest_check", "false"),
					resource.TestCheckResourceAttr(resourceName, "stop_before_destroy", "true"),
					resource.TestCheckResourceAttr(resourceName, "delete_eip_on_termination", "true"),
					resource.TestCheckResourceAttr(resourceName, "system_disk_size", "50"),
					resource.TestCheckResourceAttr(resourceName, "agency_name", "test111"),
					resource.TestCheckResourceAttr(resourceName, "agent_list", "hss"),
					resource.TestCheckResourceAttr(resourceName, "metadata.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					// resource.TestCheckResourceAttr(resourceName, "auto_terminate_time", "2025-10-10T11:11:00Z"),
				),
			},
			{
				Config: testAccComputeInstance_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test update"),
					resource.TestCheckResourceAttr(resourceName, "hostname", "hostname-update"),
					resource.TestCheckResourceAttr(resourceName, "system_disk_size", "60"),
					resource.TestCheckResourceAttr(resourceName, "agency_name", "test222"),
					resource.TestCheckResourceAttr(resourceName, "agent_list", "ces"),
					resource.TestCheckResourceAttr(resourceName, "metadata.foo", "bar2"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key2", "value2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
					resource.TestCheckResourceAttr(resourceName, "network.0.source_dest_check", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_terminate_time", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"stop_before_destroy", "delete_eip_on_termination", "data_disks", "metadata", "user_data",
				},
			},
		},
	})
}

func TestAccComputeInstance_powerAction(t *testing.T) {
	var instance cloudservers.CloudServer

	rName := acceptance.RandomAccResourceName()
	resourceName := "sbercloud_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_powerAction(rName, "OFF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "OFF"),
					resource.TestCheckResourceAttr(resourceName, "status", "SHUTOFF"),
				),
			},
			{
				Config: testAccComputeInstance_powerAction(rName, "ON"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "ON"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccComputeInstance_powerAction(rName, "REBOOT"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "REBOOT"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccComputeInstance_powerAction(rName, "FORCE-REBOOT"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "FORCE-REBOOT"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccComputeInstance_powerAction(rName, "FORCE-OFF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "FORCE-OFF"),
					resource.TestCheckResourceAttr(resourceName, "status", "SHUTOFF"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"stop_before_destroy",
					"delete_eip_on_termination",
					"power_action",
				},
			},
		},
	})
}

func TestAccComputeInstance_disk_encryption(t *testing.T) {
	var instance cloudservers.CloudServer

	rName := acceptance.RandomAccResourceName()
	resourceName := "sbercloud_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckKms(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_disk_encryption(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(resourceName, "volume_attached.1.kms_key_id",
						"sbercloud_kms_key.test", "id"),
				),
			},
		},
	})
}

func TestAccComputeInstance_withEPS(t *testing.T) {
	var instance cloudservers.CloudServer

	srcEPS := acceptance.SBC_ENTERPRISE_PROJECT_ID_TEST
	destEPS := acceptance.SBC_ENTERPRISE_MIGRATE_PROJECT_ID_TEST
	rName := acceptance.RandomAccResourceName()
	resourceName := "sbercloud_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_withEPS(rName, srcEPS),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testAccComputeInstance_withEPS(rName, destEPS),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", destEPS),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"stop_before_destroy", "delete_eip_on_termination",
				},
			},
		},
	})
}

func testAccCheckComputeInstanceDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	computeClient, err := cfg.ComputeV1Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_compute_instance" {
			continue
		}

		server, err := cloudservers.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "DELETED" {
				return fmt.Errorf("instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckComputeInstanceExists(n string, instance *cloudservers.CloudServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		computeClient, err := cfg.ComputeV1Client(acceptance.SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating compute client: %s", err)
		}

		found, err := cloudservers.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("instance not found")
		}

		*instance = *found

		return nil
	}
}

const testAccCompute_data = `
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

data "sbercloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "sbercloud_networking_secgroup" "test" {
  name = "default"
}
`

func testAccComputeInstance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_compute_instance" "test" {
  name                = "%s"
  description         = "terraform test"
  hostname            = "hostname-test"
  image_id            = data.sbercloud_images_image.test.id
  flavor_id           = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids  = [data.sbercloud_networking_secgroup.test.id]
  stop_before_destroy = true
  agency_name         = "test111"
  agent_list          = "hss"

  user_data = <<EOF
#! /bin/bash
echo user_test > /home/user.txt
EOF

  network {
    uuid              = data.sbercloud_vpc_subnet.test.id
    source_dest_check = false
  }

  system_disk_type = "SAS"
  system_disk_size = 50

  data_disks {
    type = "SAS"
    size = "10"
  }

  metadata = {
    foo = "bar"
    key = "value"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCompute_data, rName)
}

func testAccComputeInstance_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_compute_instance" "test" {
  name                = "%s-update"
  description         = "terraform test update"
  hostname            = "hostname-update"
  image_id            = data.sbercloud_images_image.test.id
  flavor_id           = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids  = [data.sbercloud_networking_secgroup.test.id]
  stop_before_destroy = true
  agency_name         = "test222"
  agent_list          = "ces"
  # auto_terminate_time = ""

  network {
    uuid              = data.sbercloud_vpc_subnet.test.id
    source_dest_check = true
  }

  system_disk_type = "SAS"
  system_disk_size = 60

  data_disks {
    type = "SAS"
    size = "10"
  }

  metadata = {
    foo  = "bar2"
    key2 = "value2"
  }

  tags = {
    foo  = "bar2"
    key2 = "value2"
  }
}
`, testAccCompute_data, rName)
}

func testAccComputeInstance_powerAction(rName, powerAction string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.sbercloud_images_image.test.id
  flavor_id          = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids = [data.sbercloud_networking_secgroup.test.id]
  availability_zone  = data.sbercloud_availability_zones.test.names[0]
  power_action       = "%s"
  system_disk_type    = "SSD"

  network {
    uuid = data.sbercloud_vpc_subnet.test.id
  }
}
`, testAccCompute_data, rName, powerAction)
}

func testAccComputeInstance_disk_encryption(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_kms_key" "test" {
  key_alias       = "%s"
  pending_days    = "7"
  key_description = "first test key"
  is_enabled      = true
}

resource "sbercloud_compute_instance" "test" {
  name                = "%s"
  image_id            = data.sbercloud_images_image.test.id
  flavor_id           = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids  = [data.sbercloud_networking_secgroup.test.id]
  stop_before_destroy = true
  agent_list          = "hss"

  network {
    uuid              = data.sbercloud_vpc_subnet.test.id
    source_dest_check = false
  }

  system_disk_type = "SAS"
  system_disk_size = 50

  data_disks {
    type = "SAS"
    size = "10"
    kms_key_id = sbercloud_kms_key.test.id
  }
}
`, testAccCompute_data, rName, rName)
}

func testAccComputeInstance_withEPS(rName, epsID string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_compute_instance" "test" {
  name                  = "%s"
  description           = "terraform test"
  image_id              = data.sbercloud_images_image.test.id
  flavor_id             = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids    = [data.sbercloud_networking_secgroup.test.id]
  enterprise_project_id = "%s"
  system_disk_type      = "SAS"
  system_disk_size      = 40

  network {
    uuid              = data.sbercloud_vpc_subnet.test.id
    source_dest_check = false
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCompute_data, rName, epsID)
}
