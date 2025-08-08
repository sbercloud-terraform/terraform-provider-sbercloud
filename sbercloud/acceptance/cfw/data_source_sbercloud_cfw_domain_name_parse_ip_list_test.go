package cfw

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCfwDomainNameParseIpList_basic(t *testing.T) {
	dataSource := "data.sbercloud_cfw_domain_name_parse_ip_list.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDomainNameParseIpList_domainName(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
				),
			},
			{
				Config: testDataSourceDomainNameParseIpList_domainNameInTheGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
				),
			},
		},
	})
}

func testDataSourceDomainNameParseIpList_domainName() string {
	return `
data "sbercloud_cfw_domain_name_parse_ip_list" "test" {
  domain_name = "www.baidu.com"
}`
}

func testDataSourceDomainNameParseIpList_domainNameInTheGroup(name string) string {
	return fmt.Sprintf(`
%[1]s

data "sbercloud_cfw_domain_name_parse_ip_list" "test" {
  domain_address_id = data.sbercloud_cfw_domain_name_groups.test.records[0].domain_names[0].domain_address_id
  group_id          = sbercloud_cfw_domain_name_group.test.id
  fw_instance_id    = "%[2]s"
}
`, testDataSourceDomainNameParseIpList_base(name), acceptance.SBC_CFW_INSTANCE_ID)
}

func testDataSourceDomainNameParseIpList_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_cfw_domain_name_group" "test" {
  fw_instance_id = "%[2]s"
  object_id      = data.sbercloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name           = "%[3]s"
  type           = 1
  description    = "network domain name group"
  
  domain_names {
    domain_name = "www.baidu.com"
    description = "baidu"
  }
}

data "sbercloud_cfw_domain_name_groups" "test" {
  fw_instance_id = "%[2]s"
  object_id      = data.sbercloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  group_id       = sbercloud_cfw_domain_name_group.test.id
  
  depends_on = [sbercloud_cfw_domain_name_group.test]
}
`, testAccDatasourceFirewalls_basic(), acceptance.SBC_CFW_INSTANCE_ID, name)
}
