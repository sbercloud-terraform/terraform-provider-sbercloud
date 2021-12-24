package vpc

import (
	"fmt"
	"testing"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVpcIdsDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.sbercloud_vpc_ids.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcIds_base(randName),
			},
			{
				Config: testAccDataSourceVpcIds_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccDataSourceVpcIds_base(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc" "test" {
  name = "%s"
  cidr = "172.16.9.0/24"
}
`, rName)
}

func testAccDataSourceVpcIds_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc_ids" "test" {}
`, testAccDataSourceVpcIds_base(rName))
}
