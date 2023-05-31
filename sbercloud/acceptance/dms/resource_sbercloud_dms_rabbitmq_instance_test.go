package dms

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/dms/v2/rabbitmq/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccDmsRabbitmqInstances_basic(t *testing.T) {
	var instance instances.Instance
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	updateName := rName + "update"
	resourceName := "sbercloud_dms_rabbitmq_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsRabbitmqInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitmqInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsRabbitmqInstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "rabbitmq"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"used_storage_space",
				},
			},
			{
				Config: testAccDmsRabbitmqInstance_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsRabbitmqInstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "rabbitmq test update"),
				),
			},
		},
	})
}

func TestAccDmsRabbitmqInstances_withEpsId(t *testing.T) {
	var instance instances.Instance
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_dms_rabbitmq_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsRabbitmqInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitmqInstance_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsRabbitmqInstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.SBC_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckDmsRabbitmqInstanceDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	dmsClient, err := config.DmsV2Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating SberCloud dms instance client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_dms_rabbitmq_instance" {
			continue
		}

		_, err := instances.Get(dmsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("The Dms rabbitmq instance still exists.")
		}
	}
	return nil
}

func testAccCheckDmsRabbitmqInstanceExists(n string, instance instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		dmsClient, err := config.DmsV2Client(acceptance.SBC_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating SberCloud dms instance client: %s", err)
		}

		v, err := instances.Get(dmsClient, rs.Primary.ID).Extract()
		if err != nil {
			return fmtp.Errorf("Error getting SberCloud dms rabbitmq instance: %s, err: %s", rs.Primary.ID, err)
		}

		if v.InstanceID != rs.Primary.ID {
			return fmtp.Errorf("The Dms rabbitmq instance not found.")
		}
		instance = *v
		return nil
	}
}

func testAccDmsRabbitmqInstance_Base(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_dms_az" "test" {}

data "sbercloud_vpc" "test" {
  name = "vpc-default"
}

data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "sbercloud_dms_product" "test" {
  engine        = "rabbitmq"
  instance_type = "cluster"
  version       = "3.7.17"
}

resource "sbercloud_networking_secgroup" "test" {
  name        = "%s"
  description = "secgroup for rabbitmq"
}
`, rName)
}

func testAccDmsRabbitmqInstance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_dms_rabbitmq_instance" "test" {
  name              = "%s"
  description       = "rabbitmq test"
  access_user       = "user"
  password          = "Rabbitmqtest@123"
  vpc_id            = data.sbercloud_vpc.test.id
  network_id        = data.sbercloud_vpc_subnet.test.id
  security_group_id = sbercloud_networking_secgroup.test.id
  available_zones   = [data.sbercloud_dms_az.test.id]
  product_id        = data.sbercloud_dms_product.test.id
  engine_version    = data.sbercloud_dms_product.test.version
  storage_space     = data.sbercloud_dms_product.test.storage
  storage_spec_code = data.sbercloud_dms_product.test.storage_spec_code
}
`, testAccDmsRabbitmqInstance_Base(rName), rName)
}

func testAccDmsRabbitmqInstance_update(rName, updateName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_dms_rabbitmq_instance" "test" {
  name              = "%s"
  description       = "rabbitmq test update"
  access_user       = "user"
  password          = "Rabbitmqtest@123"
  vpc_id            = data.sbercloud_vpc.test.id
  network_id        = data.sbercloud_vpc_subnet.test.id
  security_group_id = sbercloud_networking_secgroup.test.id
  available_zones   = [data.sbercloud_dms_az.test.id]
  product_id        = data.sbercloud_dms_product.test.id
  engine_version    = data.sbercloud_dms_product.test.version
  storage_space     = data.sbercloud_dms_product.test.storage
  storage_spec_code = data.sbercloud_dms_product.test.storage_spec_code
}
`, testAccDmsRabbitmqInstance_Base(rName), updateName)
}

func testAccDmsRabbitmqInstance_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_dms_rabbitmq_instance" "test" {
  name                  = "%s"
  description           = "rabbitmq test"
  access_user           = "user"
  password              = "Rabbitmqtest@123"
  vpc_id                = data.sbercloud_vpc.test.id
  network_id            = data.sbercloud_vpc_subnet.test.id
  security_group_id     = sbercloud_networking_secgroup.test.id
  available_zones       = [data.sbercloud_dms_az.test.id]
  product_id            = data.sbercloud_dms_product.test.id
  engine_version        = data.sbercloud_dms_product.test.version
  storage_space         = data.sbercloud_dms_product.test.storage
  storage_spec_code     = data.sbercloud_dms_product.test.storage_spec_code
  enterprise_project_id = "%s"
}
`, testAccDmsRabbitmqInstance_Base(rName), rName, acceptance.SBC_ENTERPRISE_PROJECT_ID_TEST)
}
