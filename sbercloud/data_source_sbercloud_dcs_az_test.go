package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDcsAZV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsAZV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsAZV1DataSourceID("data.sbercloud_dcs_az.az1"),
					resource.TestCheckResourceAttr("data.sbercloud_dcs_az.az1", "code", "ru-moscow-1a"),
					resource.TestCheckResourceAttr("data.sbercloud_dcs_az.az1", "port", "443"),
					testAccCheckDcsAZV1DataSourceID("data.sbercloud_dcs_az.az2"),
					resource.TestCheckResourceAttr("data.sbercloud_dcs_az.az2", "code", "ru-moscow-1b"),
					resource.TestCheckResourceAttr("data.sbercloud_dcs_az.az2", "port", "443"),
				),
			},
		},
	})
}

func testAccCheckDcsAZV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Dcs az data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Dcs az data source ID not set")
		}

		return nil
	}
}

var testAccDcsAZV1DataSource_basic = fmt.Sprintf(`
data "sbercloud_availability_zones" "test" {}

data "sbercloud_dcs_az" "az1" {
  code = data.sbercloud_availability_zones.test.names[0]
  port = "443"
}

data "sbercloud_dcs_az" "az2" {
  code = data.sbercloud_availability_zones.test.names[1]
  port = "443"
}
`)
