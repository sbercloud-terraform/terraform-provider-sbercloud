package sbercloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"

	"github.com/huaweicloud/golangsdk/openstack/rds/v3/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccRdsInstanceV3_basic(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceType := "sbercloud_rds_instance"
	resourceName := "sbercloud_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.c6.large.4"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip", "192.168.0.58"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
				),
			},
			{
				Config: testAccRdsInstanceV3_update(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "2"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.c6.xlarge.4"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "100"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_updated"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"db",
					"status",
				},
			},
		},
	})
}

func TestAccRdsInstanceV3_withEpsId(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceType := "sbercloud_rds_instance"
	resourceName := "sbercloud_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3_epsId(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", SBC_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3_ha(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceType := "sbercloud_rds_instance"
	resourceName := "sbercloud_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3_ha(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.c6.large.4.ha"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip", "192.168.0.58"),
					resource.TestCheckResourceAttr(resourceName, "ha_replication_mode", "async"),
				),
			},
		},
	})
}

func testAccCheckRdsInstanceV3Destroy(rsType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config.Config)
		client, err := config.RdsV3Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud rds client: %s", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != rsType {
				continue
			}

			id := rs.Primary.ID
			instance, err := getRdsInstanceByID(client, id)
			if err != nil {
				return err
			}
			if instance.Id != "" {
				return fmt.Errorf("%s (%s) still exists", rsType, id)
			}
		}
		return nil
	}
}

func testAccCheckRdsInstanceV3Exists(name string, instance *instances.RdsInstanceResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		client, err := config.RdsV3Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud rds client: %s", err)
		}

		found, err := getRdsInstanceByID(client, id)
		if err != nil {
			return fmt.Errorf("Error checking %s exist, err=%s", name, err)
		}
		if found.Id == "" {
			return fmt.Errorf("resource %s does not exist", name)
		}

		instance = found
		return nil
	}
}

func getRdsInstanceByID(client *golangsdk.ServiceClient, instanceID string) (*instances.RdsInstanceResponse, error) {
	listOpts := instances.ListOpts{
		Id: instanceID,
	}
	pages, err := instances.List(client, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("An error occured while querying rds instance %s: %s", instanceID, err)
	}

	resp, err := instances.ExtractRdsInstances(pages)
	if err != nil {
		return nil, err
	}

	instanceList := resp.Instances
	if len(instanceList) == 0 {
		// return an empty rds instance
		log.Printf("[WARN] can not find the specified rds instance %s", instanceID)
		instance := new(instances.RdsInstanceResponse)
		return instance, nil
	}

	if len(instanceList) > 1 {
		return nil, fmt.Errorf("retrieving more than one rds instance by %s", instanceID)
	}
	if instanceList[0].Id != instanceID {
		return nil, fmt.Errorf("the id of rds instance was expected %s, but got %s",
			instanceID, instanceList[0].Id)
	}

	return &instanceList[0], nil
}

func testAccRdsInstanceV3_base(name string) string {
	return fmt.Sprintf(`
data "sbercloud_availability_zones" "test" {}

resource "sbercloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_subnet" "test" {
  name          = "%s"
  cidr          = "192.168.0.0/24"
  gateway_ip    = "192.168.0.1"
  primary_dns   = "100.125.1.250"
  secondary_dns = "100.125.21.250"
  vpc_id        = sbercloud_vpc.test.id
}

resource "sbercloud_networking_secgroup" "test" {
  name = "%s"
}
`, name, name, name)
}

func testAccRdsInstanceV3_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_rds_instance" "test" {
  name              = "%s"
  flavor            = "rds.pg.c6.large.4"
  availability_zone = [data.sbercloud_availability_zones.test.names[0]]
  security_group_id = sbercloud_networking_secgroup.test.id
  subnet_id         = sbercloud_vpc_subnet.test.id
  vpc_id            = sbercloud_vpc.test.id
  time_zone         = "UTC+08:00"
  fixed_ip          = "192.168.0.58"

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8635
  }
  volume {
    type = "HIGH"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testAccRdsInstanceV3_base(name), name)
}

// volume.size, backup_strategy, flavor and tags will be updated
func testAccRdsInstanceV3_update(name string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_rds_instance" "test" {
  name              = "%s"
  flavor            = "rds.pg.c6.xlarge.4"
  availability_zone = [data.sbercloud_availability_zones.test.names[0]]
  security_group_id = sbercloud_networking_secgroup.test.id
  subnet_id         = sbercloud_vpc_subnet.test.id
  vpc_id            = sbercloud_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8635
  }
  volume {
    type = "HIGH"
    size = 100
  }
  backup_strategy {
    start_time = "09:00-10:00"
    keep_days  = 2
  }

  tags = {
    key1 = "value"
    foo  = "bar_updated"
  }
}
`, testAccRdsInstanceV3_base(name), name)
}

func testAccRdsInstanceV3_epsId(name string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_rds_instance" "test" {
  name                  = "%s"
  flavor                = "rds.pg.c6.large.4"
  availability_zone     = [data.sbercloud_availability_zones.test.names[0]]
  security_group_id     = sbercloud_networking_secgroup.test.id
  subnet_id             = sbercloud_vpc_subnet.test.id
  vpc_id                = sbercloud_vpc.test.id
  enterprise_project_id = "%s"

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8635
  }
  volume {
    type = "HIGH"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
}
`, testAccRdsInstanceV3_base(name), name, SBC_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccRdsInstanceV3_ha(name string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_rds_instance" "test" {
  name                = "%s"
  flavor              = "rds.pg.c6.large.4.ha"
  security_group_id   = sbercloud_networking_secgroup.test.id
  subnet_id           = sbercloud_vpc_subnet.test.id
  vpc_id              = sbercloud_vpc.test.id
  time_zone           = "UTC+08:00"
  fixed_ip            = "192.168.0.58"
  ha_replication_mode = "async"
  availability_zone   = [
    data.sbercloud_availability_zones.test.names[0],
    data.sbercloud_availability_zones.test.names[1],
  ]

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8635
  }
  volume {
    type = "HIGH"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testAccRdsInstanceV3_base(name), name)
}
