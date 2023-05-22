package dcs

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDcsProductV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsProductV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsProductV1DataSourceID("data.sbercloud_dcs_product.product1"),
					resource.TestCheckResourceAttr(
						"data.sbercloud_dcs_product.product1", "spec_code", "dcs.single_node"),
				),
			},
		},
	})
}

func testAccCheckDcsProductV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Dcs product data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Dcs product data source ID not set")
		}

		return nil
	}
}

var testAccDcsProductV1DataSource_basic = fmt.Sprintf(`
data "sbercloud_dcs_product" "product1" {
spec_code = "dcs.single_node"
}
`)
