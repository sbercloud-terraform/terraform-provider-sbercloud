package cfw

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCfwAddressGroupMembers_basic(t *testing.T) {
	dataSource := "data.sbercloud_cfw_address_group_members.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCfwAddressGroupMembers_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.item_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.address"),
					resource.TestCheckResourceAttrSet(dataSource, "records.1.item_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.1.address"),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_address_filter_useful", "true"),
					resource.TestCheckOutput("is_keyword_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCfwAddressGroupMembers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  id       = sbercloud_cfw_address_group_member.m2.id
  address  = sbercloud_cfw_address_group_member.m2.address
  keyword = "member 1"
}

data "sbercloud_cfw_address_group_members" "test" {
  group_id = sbercloud_cfw_address_group.test.id

  depends_on = [
    sbercloud_cfw_address_group_member.m1,
    sbercloud_cfw_address_group_member.m2,
  ]
}

output "is_default_filter_useful" {
  value = length(data.sbercloud_cfw_address_group_members.test.records) >= 2
}

data "sbercloud_cfw_address_group_members" "filter_by_id" {
  group_id = sbercloud_cfw_address_group.test.id
  item_id  = local.id

  depends_on = [
    sbercloud_cfw_address_group_member.m1,
    sbercloud_cfw_address_group_member.m2,
  ]
}

output "is_id_filter_useful" {
  value = length(data.sbercloud_cfw_address_group_members.filter_by_id.records) >= 1 && alltrue(
    [for v in data.sbercloud_cfw_address_group_members.filter_by_id.records[*] : v.item_id == local.id]
  )
}

data "sbercloud_cfw_address_group_members" "filter_by_address" {
  group_id = sbercloud_cfw_address_group.test.id
  address  = local.address
  
  depends_on = [
    sbercloud_cfw_address_group_member.m1,
    sbercloud_cfw_address_group_member.m2,
  ]
}

output "is_address_filter_useful" {
  value = length(data.sbercloud_cfw_address_group_members.filter_by_address.records) >= 1 && alltrue(
    [for v in data.sbercloud_cfw_address_group_members.filter_by_address.records[*] : v.address == local.address]
  )
}

data "sbercloud_cfw_address_group_members" "filter_by_keyword" {
  group_id = sbercloud_cfw_address_group.test.id
  key_word = local.keyword
  
  depends_on = [
    sbercloud_cfw_address_group_member.m1,
    sbercloud_cfw_address_group_member.m2,
  ]
}

output "is_keyword_filter_useful" {
  value = length(data.sbercloud_cfw_address_group_members.filter_by_keyword.records) >= 1 && alltrue(
    [for v in data.sbercloud_cfw_address_group_members.filter_by_keyword.records[*] : can(regex(local.keyword, v.description))]
  )
}
`, testDataSourceCfwAddressGroupMembers_base(name))
}

// sbercloud_cfw_address_group_members
func testDataSourceCfwAddressGroupMembers_base(name string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_cfw_address_group_member" "m1" {
  group_id    = sbercloud_cfw_address_group.test.id
  address     = "192.168.0.1"
  description = "member 1"
}

resource "sbercloud_cfw_address_group_member" "m2" {
  group_id    = sbercloud_cfw_address_group.test.id
  address     = "192.168.0.2"
  description = "member 2"
}
`, testAddressGroup_basic(name))
}
