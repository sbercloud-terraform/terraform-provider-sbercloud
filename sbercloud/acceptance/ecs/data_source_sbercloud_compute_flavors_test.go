package ecs

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccEcsFlavorsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsFlavorsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEcsFlavorDataSourceID("data.sbercloud_compute_flavors.this"),
				),
			},
		},
	})
}

func testAccCheckEcsFlavorDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find compute flavors data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Compute Flavors data source ID not set")
		}

		return nil
	}
}

const testAccEcsFlavorsDataSource_basic = `
data "sbercloud_compute_flavors" "this" {
	performance_type = "normal"
	cpu_core_count   = 2
	memory_size      = 4
}
`
