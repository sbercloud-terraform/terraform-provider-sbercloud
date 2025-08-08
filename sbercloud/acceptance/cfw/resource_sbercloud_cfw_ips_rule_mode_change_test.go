package cfw

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCfwIpsRuleModeChange_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testResourceCfwIpsRuleModeChange_basic(),
			},
		},
	})
}

func testResourceCfwIpsRuleModeChange_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_cfw_ips_rule_mode_change" "test" {
  object_id = data.sbercloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  ips_ids   = [340710, 340922]
  status    = "CLOSE"
}
`, testAccDatasourceFirewalls_basic())
}
