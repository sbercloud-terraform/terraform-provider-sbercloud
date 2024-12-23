package apig

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceInstanceSupportedFeatures_basic(t *testing.T) {
	var (
		rName = "data.sbercloud_apig_instance_supported_features.test"
		dc    = acceptance.InitDataSourceCheck(rName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceInstanceSupportedFeatures_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "features.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
		},
	})
}

func testAccDataSourceInstanceSupportedFeatures_basic() string {
	return fmt.Sprintf(`
resource "sbercloud_apig_instance_feature" "test" {
  instance_id = "%[1]s"
  name        = "ratelimit"
  enabled     = true

  config = jsonencode({
    api_limits = 200
  })
}

data "sbercloud_apig_instance_supported_features" "test" {
  instance_id = "%[1]s"
}
`, acceptance.SBC_APIG_DEDICATED_INSTANCE_ID)
}
