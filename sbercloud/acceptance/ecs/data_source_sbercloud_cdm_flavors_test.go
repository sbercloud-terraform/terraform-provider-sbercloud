package ecs

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCdmFlavorV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
`, acceptance.SBC_REGION_NAME)
}
