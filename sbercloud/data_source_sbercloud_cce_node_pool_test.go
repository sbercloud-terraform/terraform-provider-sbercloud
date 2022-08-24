package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCCENodePoolV3DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "data.sbercloud_cce_node_pool.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePoolV3DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolV3DataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func testAccCheckCCENodePoolV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find node pools data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Node pool data source ID not set ")
		}

		return nil
	}
}

func testAccCCENodePoolV3DataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_cce_node_pool" "test" {
  cluster_id = sbercloud_cce_cluster.test.id
  name       = sbercloud_cce_node_pool.test.name
}
`, testAccCCENodePool_basic(rName))
}
