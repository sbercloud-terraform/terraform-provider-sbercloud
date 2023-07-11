package vpc

import (
	"fmt"
	"testing"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVpcEipsDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.sbercloud_vpc_eips.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcEips_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.status", "UNBOUND"),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.type", "5_bgp"),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.ip_version", "4"),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.bandwidth_name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.bandwidth_size", "5"),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.bandwidth_share_type", "PER"),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.tags.key", "value"),
					resource.TestCheckResourceAttrPair(dataSourceName, "eips.0.id",
						"sbercloud_vpc_eip.test", "id"),
				),
			},
		},
	})
}

func testAccDataSourceVpcEips_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc_eips" "test" {
  public_ips = [sbercloud_vpc_eip.test.address]
}
`, testAccVpcEip_tags(rName))
}

func TestAccVpcEipsDataSource_byTag(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.sbercloud_vpc_eips.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcEips_byTag(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.status", "UNBOUND"),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.type", "5_bgp"),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.ip_version", "4"),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.bandwidth_name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.bandwidth_size", "5"),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.bandwidth_share_type", "PER"),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.tags.key", "value"),
					resource.TestCheckResourceAttrPair(dataSourceName, "eips.0.id",
						"sbercloud_vpc_eip.test", "id"),
				),
			},
		},
	})
}

func testAccDataSourceVpcEips_byTag(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc_eips" "test" {
  public_ips = [sbercloud_vpc_eip.test.address]

  tags = {
    foo = "bar"
  }
}
`, testAccVpcEip_tags(rName))
}
