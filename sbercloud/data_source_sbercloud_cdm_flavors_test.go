package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCdmFlavorV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdmFlavorV1DataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdmFlavorV1DataSourceID("data.sbercloud_cdm_flavors.flavor"),
				),
			},
		},
	})
}

func testAccCheckCdmFlavorV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find cdm data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Cdm data source ID not set ")
		}

		return nil
	}
}

func testAccCdmFlavorV1DataSource_basic() string {
	return fmt.Sprintf(`
data "sbercloud_cdm_flavors" "flavor" {
  region = "%s"
}
`, SBC_REGION_NAME)
}
