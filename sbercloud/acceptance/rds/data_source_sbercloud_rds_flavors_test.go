package rds

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccRdsFlavorV3DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsFlavorV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsFlavorV3DataSourceID("data.sbercloud_rds_flavors.flavor"),
				),
			},
		},
	})
}

func testAccCheckRdsFlavorV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find rds data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Rds data source ID not set ")
		}

		return nil
	}
}

var testAccRdsFlavorV3DataSource_basic = `

data "sbercloud_rds_flavors" "flavor" {
  db_type = "PostgreSQL"
  db_version = "10"
  instance_mode = "ha"
}
`
