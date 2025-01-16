package cce

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccClusterCertificateDataSource_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	datasourceName := "data.sbercloud_cce_cluster_certificate.test"
	dc := acceptance.InitDataSourceCheck(datasourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testClousterCertificate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(datasourceName, "duration", "30"),
					resource.TestCheckResourceAttr(datasourceName, "clusters.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "users.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "contexts.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "current_context"),
					resource.TestCheckResourceAttrSet(datasourceName, "kube_config_raw"),
				),
			},
		},
	})
}

func testClousterCertificate_basic(name string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_cce_cluster_certificate" "test" {
  cluster_id = sbercloud_cce_cluster.test.id
  duration   = 30
}`, testAccCCEClusterV3_basic(name))
}
