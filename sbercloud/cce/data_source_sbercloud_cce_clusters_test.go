package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccCCEClustersDataSource_basic(t *testing.T) {
	dataSourceName := "data.sbercloud_cce_clusters.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClustersV3DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "clusters.0.name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "clusters.0.status", "Available"),
					resource.TestCheckResourceAttr(dataSourceName, "clusters.0.cluster_type", "VirtualMachine"),
				),
			},
		},
	})
}

func testAccCCEClustersV3DataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_cce_clusters" "test" {
  name = sbercloud_cce_cluster.test.name

  depends_on = [sbercloud_cce_cluster.test]
}
`, testAccCceCluster_config(rName))
}
