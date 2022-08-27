package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccCCEAddonTemplateV3DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddonTemplateV3DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sbercloud_cce_addon_template.coredns_test", "spec"),
					resource.TestCheckResourceAttrSet("data.sbercloud_cce_addon_template.metrics-server_test", "spec"),
				),
			},
		},
	})
}

func testAccCCEAddonTemplateV3DataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_cce_addon_template" "coredns_test" {
  cluster_id = sbercloud_cce_cluster.test.id
  name       = "coredns"
  version    = "1.17.15"
}

data "sbercloud_cce_addon_template" "metrics-server_test" {
  cluster_id = sbercloud_cce_cluster.test.id
  name       = "metrics-server"
  version    = "1.1.10"
}
`, testAccCCEClusterV3_basic(rName))
}
