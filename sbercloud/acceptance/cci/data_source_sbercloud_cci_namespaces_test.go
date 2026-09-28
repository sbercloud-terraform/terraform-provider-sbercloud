package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccDataSourceV2Namespaces_basic(t *testing.T) {
	dataSourceName := "data.sbercloud_cci_namespaces.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceV2Namespaces_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.0.name", rName),
					resource.TestCheckResourceAttrSet(dataSourceName, "namespaces.0.api_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "namespaces.0.kind"),
					resource.TestCheckResourceAttrSet(dataSourceName, "namespaces.0.annotations.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "namespaces.0.labels.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "namespaces.0.creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSourceName, "namespaces.0.resource_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "namespaces.0.uid"),
					resource.TestCheckResourceAttrSet(dataSourceName, "namespaces.0.finalizers.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "namespaces.0.status"),
				),
			},
		},
	})
}

func testAccDataSourceV2Namespaces_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "sbercloud_cci_namespaces" "test" {
  depends_on = [sbercloud_cci_namespace.test]

  name = "%[2]s"
}
`, testAccV2Namespace_basic(rName), rName)
}
