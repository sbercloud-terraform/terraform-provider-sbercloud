package rds

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"log"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/rds/v3/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccRdsInstanceV3_basic(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceType := "sbercloud_rds_instance"
	resourceName := "sbercloud_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.x1.large.2"),
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
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.x1.large.2"),
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
					"availability_zone",
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
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3_epsId(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.SBC_ENTERPRISE_PROJECT_ID_TEST),
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
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3_ha(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.x1.large.2.ha"),
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
		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.RdsV3Client(acceptance.SBC_REGION_NAME)
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

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.RdsV3Client(acceptance.SBC_REGION_NAME)
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
  primary_dns   = "100.125.13.59"
  secondary_dns = "100.125.65.14"
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
  flavor            = "rds.pg.x1.large.2"
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
    type = "CLOUDSSD"
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
  flavor            = "rds.pg.x1.large.2"
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
    type = "CLOUDSSD"
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
  flavor                = "rds.pg.x1.large.2"
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
    type = "CLOUDSSD"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
}
`, testAccRdsInstanceV3_base(name), name, acceptance.SBC_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccRdsInstanceV3_ha(name string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_rds_instance" "test" {
  name                = "%s"
  flavor              = "rds.pg.x1.large.2.ha"
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
    type = "CLOUDSSD"
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

// if the instance flavor has been changed, then a temp instance will be kept for 12 hours,
// the binding relationship between instance and security group or subnet cannot be unbound
// when deleting the instance in this period time, so we cannot create a new vpc, subnet and
// security group in the test case, otherwise, they cannot be deleted when destroy the resource
func testAccRdsInstance_mysql_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

data "sbercloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 8
}

resource "sbercloud_rds_instance" "test" {
  name                   = "%[2]s"
  flavor                 = data.sbercloud_rds_flavors.test.flavors[0].name
  security_group_id      = sbercloud_networking_secgroup.test.id
  subnet_id              = data.sbercloud_vpc_subnet.test.id
  vpc_id                 = data.sbercloud_vpc.test.id
  availability_zone      = slice(sort(data.sbercloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)
  ssl_enable             = true  
  binlog_retention_hours = "12"
  read_write_permissions = "readonly"

  # seconds_level_monitoring_enabled  = false
  # seconds_level_monitoring_interval = 1

  db {
    type     = "MySQL"
    version  = "8.0"
    port     = 3306
  }

  backup_strategy {
    start_time = "08:15-09:15"
    keep_days  = 3
    period     = 1
  }

  volume {
    type              = "CLOUDSSD"
    size              = 40
    limit_size        = 0
    trigger_threshold = 15
  }

  parameters {
    name  = "back_log"
    value = "2000"
  }
}
`, testAccRdsInstance_base(name), name)
}

func testAccRdsInstance_mysql_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

data "sbercloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "sbercloud_rds_instance" "test" {
  name                   = "%[3]s"
  flavor                 = data.sbercloud_rds_flavors.test.flavors[1].name
  security_group_id      = sbercloud_networking_secgroup.test.id
  subnet_id              = data.sbercloud_vpc_subnet.test.id
  vpc_id                 = data.sbercloud_vpc.test.id
  availability_zone      = slice(sort(data.sbercloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)
  ssl_enable             = false
  param_group_id         = sbercloud_rds_parametergroup.pg_1.id
  binlog_retention_hours = "0"
  read_write_permissions = "readwrite"

  seconds_level_monitoring_enabled  = true
  seconds_level_monitoring_interval = 5

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = 3308
  }

  backup_strategy {
    start_time = "18:15-19:15"
    keep_days  = 5
    period     = 3
  }

  volume {
    type              = "CLOUDSSD"
    size              = 40
    limit_size        = 500
    trigger_threshold = 20
  }

  parameters {
    name  = "connect_timeout"
    value = "14"
  }
}
`, testAccRdsInstance_base(name), testAccRdsConfig_basic(name), name)
}

func testAccRdsInstance_sqlserver(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 8634
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = sbercloud_networking_secgroup.test.id
}

data "sbercloud_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2017_SE"
  instance_mode = "single"
  # group_type    = "dedicated"
  vcpus         = 4
}

resource "sbercloud_rds_instance" "test" {
  depends_on        = [sbercloud_networking_secgroup_rule.ingress]
  name              = "%[2]s"
  flavor            = data.sbercloud_rds_flavors.test.flavors[0].name
  security_group_id = sbercloud_networking_secgroup.test.id
  subnet_id         = data.sbercloud_vpc_subnet.test.id
  vpc_id            = data.sbercloud_vpc.test.id
  collation         = "Chinese_PRC_CI_AS"
  # tde_enabled       = true

  availability_zone = [
    data.sbercloud_availability_zones.test.names[0],
  ]

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2017_SE"
    port     = 8634
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, testAccRdsInstance_base(name), name)
}

func testAccRdsInstance_sqlserver_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 8634
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = sbercloud_networking_secgroup.test.id
}

data "sbercloud_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2017_SE"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "sbercloud_rds_instance" "test" {
  depends_on        = [sbercloud_networking_secgroup_rule.ingress]
  name              = "%[2]s"
  flavor            = data.sbercloud_rds_flavors.test.flavors[0].name
  security_group_id = sbercloud_networking_secgroup.test.id
  subnet_id         = data.sbercloud_vpc_subnet.test.id
  vpc_id            = data.sbercloud_vpc.test.id
  collation         = "Chinese_PRC_CI_AI"

  availability_zone = [
    data.sbercloud_availability_zones.test.names[0],
  ]

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2017_SE"
    port     = 8634
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, testAccRdsInstance_base(name), name)
}

func testAccRdsInstance_sqlserver_msdtcHosts_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 8634
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = sbercloud_networking_secgroup.test.id
}

data "sbercloud_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2019_SE"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

data "sbercloud_compute_flavors" "test" {
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "sbercloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}
`, testAccRdsInstance_base(name), name)
}

func testAccRdsInstance_sqlserver_msdtcHosts(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_compute_instance" "ecs_1" {
  name               = "%[2]s_ecs_1"
  image_id           = data.sbercloud_images_image.test.id
  flavor_id          = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids = [sbercloud_networking_secgroup.test.id]
  availability_zone  = data.sbercloud_availability_zones.test.names[0]

  network {
    uuid = data.sbercloud_vpc_subnet.test.id
  }
}

resource "sbercloud_rds_instance" "test" {
  depends_on        = [sbercloud_networking_secgroup_rule.ingress]
  name              = "%[2]s"
  flavor            = data.sbercloud_rds_flavors.test.flavors[0].name
  security_group_id = sbercloud_networking_secgroup.test.id
  subnet_id         = data.sbercloud_vpc_subnet.test.id
  vpc_id            = data.sbercloud_vpc.test.id
  collation         = "Chinese_PRC_CI_AS"

  availability_zone = [
    data.sbercloud_availability_zones.test.names[0],
  ]

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2019_SE"
    port     = 8634
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }

  msdtc_hosts {
    ip        = sbercloud_compute_instance.ecs_1.access_ip_v4
    host_name = "msdtc-host-name-1"
  }
}
`, testAccRdsInstance_sqlserver_msdtcHosts_base(name), name)
}

func testAccRdsInstance_sqlserver_msdtcHosts_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_compute_instance" "ecs_1" {
  name               = "%[2]s_ecs_1"
  image_id           = data.sbercloud_images_image.test.id
  flavor_id          = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids = [sbercloud_networking_secgroup.test.id]
  availability_zone  = data.sbercloud_availability_zones.test.names[0]

  network {
    uuid = data.sbercloud_vpc_subnet.test.id
  }
}

resource "sbercloud_compute_instance" "ecs_2" {
  name               = "%[2]s_ecs_2"
  image_id           = data.sbercloud_images_image.test.id
  flavor_id          = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids = [sbercloud_networking_secgroup.test.id]
  availability_zone  = data.sbercloud_availability_zones.test.names[0]

  network {
    uuid = data.sbercloud_vpc_subnet.test.id
  }
}

resource "sbercloud_rds_instance" "test" {
  depends_on        = [sbercloud_networking_secgroup_rule.ingress]
  name              = "%[2]s"
  flavor            = data.sbercloud_rds_flavors.test.flavors[0].name
  security_group_id = sbercloud_networking_secgroup.test.id
  subnet_id         = data.sbercloud_vpc_subnet.test.id
  vpc_id            = data.sbercloud_vpc.test.id
  collation         = "Chinese_PRC_CI_AS"

  availability_zone = [
    data.sbercloud_availability_zones.test.names[0],
  ]

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2019_SE"
    port     = 8634
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }

  msdtc_hosts {
    ip        = sbercloud_compute_instance.ecs_1.access_ip_v4
    host_name = "msdtc-host-name-1"
  }
  msdtc_hosts {
    ip        = sbercloud_compute_instance.ecs_2.access_ip_v4
    host_name = "msdtc-host-name-2"
  }
}
`, testAccRdsInstance_sqlserver_msdtcHosts_base(name), name)
}
