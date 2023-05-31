package dms

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDmsAZV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsAZV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsAZV1DataSourceID("data.sbercloud_dms_az.az1"),
					resource.TestCheckResourceAttr(
						"data.sbercloud_dms_az.az1", "name", "可用区1"),
					resource.TestCheckResourceAttr(
						"data.sbercloud_dms_az.az1", "port", "443"),
					resource.TestCheckResourceAttr(
						"data.sbercloud_dms_az.az1", "code", "ru-moscow-1a"),
				),
			},
		},
	})
}

func testAccCheckDmsAZV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Dms az data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Dms az data source ID not set")
		}

		return nil
	}
}

var testAccDmsAZV1DataSource_basic = `
data "sbercloud_availability_zones" "test" {}

data "sbercloud_dms_az" "az1" {
name = "可用区1"
port = "443"
code = data.sbercloud_availability_zones.test.names[0]
}`
