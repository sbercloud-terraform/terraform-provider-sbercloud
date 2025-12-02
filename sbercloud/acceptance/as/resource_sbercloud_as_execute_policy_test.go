package as

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceExecutePolicy_basic(t *testing.T) {
	// Avoid CheckDestroy because this resource is a one-time action resource and there is nothing in the destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckASScalingPolicyID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceExecutePolicy_basic(),
			},
		},
	})
}

func testAccResourceExecutePolicy_basic() string {
	return fmt.Sprintf(`
resource "sbercloud_as_execute_policy" "test" {
  scaling_policy_id = "%s"
}
`, acceptance.SBC_AS_SCALING_POLICY_ID)
}
