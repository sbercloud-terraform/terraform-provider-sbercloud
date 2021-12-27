package vpc

import (
	"fmt"
	"testing"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVpcsDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.sbercloud_vpcs.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcs_base(randName, randCidr),
			},
			{
				Config: testAccDataSourceVpcs_basic(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.status", "OK"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpcs.0.id",
						"${sbercloud_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccDataSourceVpcs_base(rName, cidr string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc" "test" {
  name = "%s"
  cidr = "%s"
}
`, rName, cidr)
}

func testAccDataSourceVpcs_basic(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpcs" "test" {
  id = sbercloud_vpc.test.id
}
`, testAccDataSourceVpcs_base(rName, cidr))
}

func TestAccVpcsDataSource_byCidr(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.sbercloud_vpcs.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcs_base(randName, randCidr),
			},
			{
				Config: testAccDataSourceVpcs_byCidr(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.status", "OK"),
				),
			},
		},
	})
}

func testAccDataSourceVpcs_byCidr(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpcs" "test" {
  cidr = sbercloud_vpc.test.cidr
}
`, testAccDataSourceVpcs_base(rName, cidr))
}

func TestAccVpcsDataSource_byName(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.sbercloud_vpcs.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcs_base(randName, randCidr),
			},
			{
				Config: testAccDataSourceVpcs_byName(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.status", "OK"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpcs.0.id",
						"${sbercloud_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccDataSourceVpcs_byName(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpcs" "test" {
  name = sbercloud_vpc.test.name
}
`, testAccDataSourceVpcs_base(rName, cidr))
}

func TestAccVpcsDataSource_byAll(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.sbercloud_vpcs.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcs_base(randName, randCidr),
			},
			{
				Config: testAccDataSourceVpcs_byAll(randName, randCidr, acceptance.SBC_ENTERPRISE_PROJECT_ID),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.enterprise_project_id",
						acceptance.SBC_ENTERPRISE_PROJECT_ID),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.status", "OK"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpcs.0.id",
						"${sbercloud_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccDataSourceVpcs_byAll(rName, cidr, enterpriseProjectID string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpcs" "test" {
  id                    = sbercloud_vpc.test.id
  name                  = sbercloud_vpc.test.name
  cidr                  = sbercloud_vpc.test.cidr
  enterprise_project_id = "%s"
  status                = "OK"
}
`, testAccDataSourceVpcs_base(rName, cidr), enterpriseProjectID)
}

func TestAccVpcsDataSource_tags(t *testing.T) {
	randName1 := acceptance.RandomAccResourceName()
	randName2 := acceptance.RandomAccResourceName()
	dataSourceName := "data.sbercloud_vpcs.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcs_tags(randName1, randName2),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "tags.foo", randName1),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.name", randName1),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.status", "OK"),
				),
			},
		},
	})
}

func testAccDataSourceVpcs_tags(rName1, rName2 string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc" "test1" {
  name = "%s"
  cidr = "172.16.0.0/24"
  tags = {
    foo = "%s"
  }
}

resource "sbercloud_vpc" "test2" {
  name = "%s"
  cidr = "10.12.2.0/24"
  tags = {
    foo = "%s"
  }
}

data "sbercloud_vpcs" "test" {
  tags = {
    foo = "%s"
  }
  depends_on = [
    sbercloud_vpc.test1,
    sbercloud_vpc.test2,
  ]
}
`, rName1, rName1, rName2, rName2, rName1)
}
