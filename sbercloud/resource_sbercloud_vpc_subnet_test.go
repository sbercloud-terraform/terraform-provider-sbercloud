package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/subnets"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccVpcSubnetV1_basic(t *testing.T) {
	var subnet subnets.Subnet

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_vpc_subnet.test"
	rNameUpdate := rName + "-updated"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetV1_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetV1Exists(resourceName, &subnet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "gateway_ip", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccVpcSubnetV1_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_updated"),
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

func testAccCheckVpcSubnetV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	subnetClient, err := config.NetworkingV1Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud vpc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_vpc_subnet" {
			continue
		}

		_, err := subnets.Get(subnetClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Subnet still exists")
		}
	}

	return nil
}
func testAccCheckVpcSubnetV1Exists(n string, subnet *subnets.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		subnetClient, err := config.NetworkingV1Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating huaweicloud Vpc client: %s", err)
		}

		found, err := subnets.Get(subnetClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Subnet not found")
		}

		*subnet = *found

		return nil
	}
}

func testAccVpcSubnetV1_basic(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_availability_zones" "test" {}

resource "sbercloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id     = sbercloud_vpc.test.id

  availability_zone = data.sbercloud_availability_zones.test.names[0]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName, rName)
}

func testAccVpcSubnetV1_update(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_availability_zones" "test" {}

resource "sbercloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id     = sbercloud_vpc.test.id

  availability_zone = data.sbercloud_availability_zones.test.names[0]

  tags = {
    foo = "bar"
    key = "value_updated"
  }
}
`, rName, rName)
}
