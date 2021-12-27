package vpc

import (
	"fmt"
	"testing"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVpcDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.sbercloud_vpc.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpc_basic(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "status", "OK"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "id", "${sbercloud_vpc.test.id}"),
				),
			},
		},
	})
}

func TestAccVpcDataSource_byCidr(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.sbercloud_vpc.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpc_base(randName, randCidr),
			},
			{
				Config: testAccDataSourceVpc_byCidr(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "status", "OK"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "id", "${sbercloud_vpc.test.id}"),
				),
			},
		},
	})
}

func TestAccVpcDataSource_byName(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.sbercloud_vpc.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpc_base(randName, randCidr),
			},
			{
				Config: testAccDataSourceVpc_byName(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "status", "OK"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "id", "${sbercloud_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccDataSourceVpc_base(rName, cidr string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc" "test" {
  name = "%s"
  cidr = "%s"
}
`, rName, cidr)
}

func testAccDataSourceVpc_basic(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc" "test" {
  id = sbercloud_vpc.test.id
}
`, testAccDataSourceVpc_base(rName, cidr))
}

func testAccDataSourceVpc_byCidr(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc" "test" {
  cidr = sbercloud_vpc.test.cidr
}
`, testAccDataSourceVpc_base(rName, cidr))
}

func testAccDataSourceVpc_byName(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc" "test" {
  name = sbercloud_vpc.test.name
}
`, testAccDataSourceVpc_base(rName, cidr))
}
