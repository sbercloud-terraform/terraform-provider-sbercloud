package dms

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dms/v1/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccDmsInstancesV1_Rabbitmq(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dms_instance_%s", acctest.RandString(5))
	resourceName := "sbercloud_dms_instance.instance_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1Instance_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
		},
	})
}

func TestAccDmsInstancesV1_Kafka(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dms_instance_%s", acctest.RandString(5))
	resourceName := "sbercloud_dms_instance.instance_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1Instance_KafkaInstance(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
		},
	})
}

func testAccCheckDmsV1InstanceDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	dmsClient, err := config.DmsV1Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud instance client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_dms_instance" {
			continue
		}

		_, err := instances.Get(dmsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("The Dms instance still exists.")
		}
	}
	return nil
}

func testAccCheckDmsV1InstanceExists(n string, instance instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		dmsClient, err := config.DmsV1Client(acceptance.SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud instance client: %s", err)
		}

		v, err := instances.Get(dmsClient, rs.Primary.ID).Extract()
		if err != nil {
			return fmt.Errorf("Error getting SberCloud instance: %s, err: %s", rs.Primary.ID, err)
		}

		if v.InstanceID != rs.Primary.ID {
			return fmt.Errorf("The Dms instance not found.")
		}
		instance = *v
		return nil
	}
}

func testAccDmsV1Instance_base(name string) string {
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

func testAccDmsV1Instance_basic(instanceName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_dms_az" "az_1" {
  code = data.sbercloud_availability_zones.test.names[1]
}

data "sbercloud_dms_product" "product_1" {
  engine        = "rabbitmq"
  instance_type = "single"
  version       = "3.7.17"
}

resource "sbercloud_dms_instance" "instance_1" {
  name              = "%s"
  engine            = "rabbitmq"
  access_user       = "user"
  password          = "Dmstest@123"
  vpc_id            = sbercloud_vpc.test.id
  subnet_id         = sbercloud_vpc_subnet.test.id
  security_group_id = sbercloud_networking_secgroup.test.id
  available_zones   = [data.sbercloud_availability_zones.test.names[1]]
  product_id        = data.sbercloud_dms_product.product_1.id
  engine_version    = data.sbercloud_dms_product.product_1.version
  storage_space     = data.sbercloud_dms_product.product_1.storage
  storage_spec_code = data.sbercloud_dms_product.product_1.storage_spec_code

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
	`, testAccDmsV1Instance_base(instanceName), instanceName)
}

func testAccDmsV1Instance_KafkaInstance(instanceName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_dms_az" "az_1" {
  code = data.sbercloud_availability_zones.test.names[1]
}

data "sbercloud_dms_product" "product_1" {
  engine        = "kafka"
  instance_type = "cluster"
  version       = "1.1.0"
}

resource "sbercloud_dms_instance" "instance_1" {
  name              = "%s"
  engine            = "kafka"
  vpc_id            = sbercloud_vpc.test.id
  subnet_id         = sbercloud_vpc_subnet.test.id
  subnet_id         = sbercloud_vpc_subnet.test.id
  security_group_id = sbercloud_networking_secgroup.test.id
  available_zones   = [data.sbercloud_availability_zones.test.names[1]]
  product_id        = data.sbercloud_dms_product.product_1.id
  engine_version    = data.sbercloud_dms_product.product_1.version
  specification     = data.sbercloud_dms_product.product_1.bandwidth
  partition_num     = data.sbercloud_dms_product.product_1.partition_num
  storage_space     = data.sbercloud_dms_product.product_1.storage
  storage_spec_code = data.sbercloud_dms_product.product_1.storage_spec_code

  tags = {
    key   = "value"
    owner = "terraform"
  }
}`, testAccDmsV1Instance_base(instanceName), instanceName)
}
