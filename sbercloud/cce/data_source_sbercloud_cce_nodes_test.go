package cce

import (
	"fmt"
	"testing"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCCENodesDataSource_basic(t *testing.T) {
	dataSourceName := "data.sbercloud_cce_nodes.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "nodes.0.name", rName),
				),
			},
		},
	})
}

func testAccCCENodesDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_cce_nodes" "test" {
  cluster_id = sbercloud_cce_cluster.test.id
  name       = sbercloud_cce_node.test.name

  depends_on = [sbercloud_cce_node.test]
}
`, testAccCceCluster_config(rName))
}
