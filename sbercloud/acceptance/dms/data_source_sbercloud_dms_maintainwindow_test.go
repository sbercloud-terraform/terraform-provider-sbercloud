package dms

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDmsMaintainWindowV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsMaintainWindowV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsMaintainWindowV1DataSourceID("data.sbercloud_dms_maintainwindow.maintainwindow1"),
					resource.TestCheckResourceAttr(
						"data.sbercloud_dms_maintainwindow.maintainwindow1", "seq", "1"),
					resource.TestCheckResourceAttr(
						"data.sbercloud_dms_maintainwindow.maintainwindow1", "begin", "22"),
				),
			},
		},
	})
}

func testAccCheckDmsMaintainWindowV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Dms maintainwindow data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Dms maintainwindow data source ID not set")
		}

		return nil
	}
}

var testAccDmsMaintainWindowV1DataSource_basic = fmt.Sprintf(`
data "sbercloud_dms_maintainwindow" "maintainwindow1" {
seq = 1
}
`)
