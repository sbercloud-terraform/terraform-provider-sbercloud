package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/vpcs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccVpcV1_basic(t *testing.T) {
	var vpc vpcs.Vpc

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_vpc.test"
	rNameUpdate := rName + "-updated"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists(resourceName, &vpc),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "status", "OK"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccVpcV1_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists(resourceName, &vpc),
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

func TestAccVpcV1_WithEpsId(t *testing.T) {
	var vpc vpcs.Vpc

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_vpc.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1_epsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists(resourceName, &vpc),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "status", "OK"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", SBC_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckVpcV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	vpcClient, err := config.NetworkingV1Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud vpc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_vpc" {
			continue
		}

		_, err := vpcs.Get(vpcClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Vpc still exists")
		}
	}

	return nil
}

func testAccCheckVpcV1Exists(n string, vpc *vpcs.Vpc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		vpcClient, err := config.NetworkingV1Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud vpc client: %s", err)
		}

		found, err := vpcs.Get(vpcClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("vpc not found")
		}

		*vpc = *found

		return nil
	}
}

func testAccVpcV1_basic(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName)
}

func testAccVpcV1_update(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc" "test" {
  name = "%s"
  cidr="192.168.0.0/16"

  tags = {
    foo = "bar"
    key = "value_updated"
  }
}
`, rName)
}

func testAccVpcV1_epsId(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
  enterprise_project_id = "%s"
}
`, rName, SBC_ENTERPRISE_PROJECT_ID_TEST)
}
