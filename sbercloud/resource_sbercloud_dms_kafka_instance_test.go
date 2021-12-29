package sbercloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/dms/v2/kafka/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccDmsKafkaInstances_basic(t *testing.T) {
	var instance instances.Instance
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	updateName := rName + "update"
	resourceName := "sbercloud_dms_kafka_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDmsKafkaInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsKafkaInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsKafkaInstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"manager_password",
					"used_storage_space",
				},
			},
			{
				Config: testAccDmsKafkaInstance_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsKafkaInstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "kafka test update"),
				),
			},
		},
	})
}

func TestAccDmsKafkaInstances_withEpsId(t *testing.T) {
	var instance instances.Instance
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_dms_kafka_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDmsKafkaInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsKafkaInstance_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsKafkaInstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", SBC_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckDmsKafkaInstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	dmsClient, err := config.DmsV2Client(SBC_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating SberCloud dms instance client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_dms_kafka_instance" {
			continue
		}

		_, err := instances.Get(dmsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("The Dms kafka instance still exists.")
		}
	}
	return nil
}

func testAccCheckDmsKafkaInstanceExists(n string, instance instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		dmsClient, err := config.DmsV2Client(SBC_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating SberCloud dms instance client: %s", err)
		}

		v, err := instances.Get(dmsClient, rs.Primary.ID).Extract()
		if err != nil {
			return fmtp.Errorf("Error getting SberCloud dms kafka instance: %s, err: %s", rs.Primary.ID, err)
		}

		if v.InstanceID != rs.Primary.ID {
			return fmtp.Errorf("The Dms kafka instance not found.")
		}
		instance = *v
		return nil
	}
}

func testAccDmsKafkaInstance_Base(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_dms_az" "test" {}

data "sbercloud_vpc" "test" {
  name = "vpc-default"
}

data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "sbercloud_dms_product" "test" {
  engine        = "kafka"
  instance_type = "cluster"
  version       = "2.3.0"
}

resource "sbercloud_networking_secgroup" "test" {
  name        = "%s"
  description = "secgroup for kafka"
}
`, rName)
}

func testAccDmsKafkaInstance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_dms_kafka_instance" "test" {
  name              = "%s"
  description       = "kafka test"
  access_user       = "user"
  password          = "Kafkatest@123"
  vpc_id            = data.sbercloud_vpc.test.id
  network_id        = data.sbercloud_vpc_subnet.test.id
  security_group_id = sbercloud_networking_secgroup.test.id
  available_zones   = [data.sbercloud_dms_az.test.id]
  product_id        = data.sbercloud_dms_product.test.id
  engine_version    = data.sbercloud_dms_product.test.version
  bandwidth         = data.sbercloud_dms_product.test.bandwidth
  storage_space     = data.sbercloud_dms_product.test.storage
  storage_spec_code = data.sbercloud_dms_product.test.storage_spec_code
  manager_user      = "kafka-user"
  manager_password  = "Kafkatest@123"
}
`, testAccDmsKafkaInstance_Base(rName), rName)
}

func testAccDmsKafkaInstance_update(rName, updateName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_dms_kafka_instance" "test" {
  name              = "%s"
  description       = "kafka test update"
  access_user       = "user"
  password          = "Kafkatest@123"
  vpc_id            = data.sbercloud_vpc.test.id
  network_id        = data.sbercloud_vpc_subnet.test.id
  security_group_id = sbercloud_networking_secgroup.test.id
  available_zones   = [data.sbercloud_dms_az.test.id]
  product_id        = data.sbercloud_dms_product.test.id
  engine_version    = data.sbercloud_dms_product.test.version
  bandwidth         = data.sbercloud_dms_product.test.bandwidth
  storage_space     = data.sbercloud_dms_product.test.storage
  storage_spec_code = data.sbercloud_dms_product.test.storage_spec_code
  manager_user      = "kafka-user"
  manager_password  = "Kafkatest@123"
}
`, testAccDmsKafkaInstance_Base(rName), updateName)
}

func testAccDmsKafkaInstance_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_dms_kafka_instance" "test" {
  name                  = "%s"
  description           = "kafka test"
  access_user           = "user"
  password              = "Kafkatest@123"
  vpc_id                = data.sbercloud_vpc.test.id
  network_id            = data.sbercloud_vpc_subnet.test.id
  security_group_id     = sbercloud_networking_secgroup.test.id
  available_zones       = [data.sbercloud_dms_az.test.id]
  product_id            = data.sbercloud_dms_product.test.id
  engine_version        = data.sbercloud_dms_product.test.version
  bandwidth             = data.sbercloud_dms_product.test.bandwidth
  storage_space         = data.sbercloud_dms_product.test.storage
  storage_spec_code     = data.sbercloud_dms_product.test.storage_spec_code
  manager_user          = "kafka-user"
  manager_password      = "Kafkatest@123"
  enterprise_project_id = "%s"
}
`, testAccDmsKafkaInstance_Base(rName), rName, SBC_ENTERPRISE_PROJECT_ID_TEST)
}
