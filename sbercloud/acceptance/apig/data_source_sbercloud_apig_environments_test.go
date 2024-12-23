package apig

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceEnvironments_basic(t *testing.T) {
	var (
		dataSourceName = "data.sbercloud_apig_environments.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		name           = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEnvironments_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "environments.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

func testAccDataSourceEnvironments_basic(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_apig_environment" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  description = "Created by script"
}

data "sbercloud_apig_environments" "test" {
  depends_on = [sbercloud_apig_environment.test]

  instance_id = "%[1]s"
  name        = sbercloud_apig_environment.test.name
}
`, acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, name)
}
