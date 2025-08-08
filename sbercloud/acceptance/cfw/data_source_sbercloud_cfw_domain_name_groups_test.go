package cfw

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCfwDomainNameGroups_basic(t *testing.T) {
	dataSource := "data.sbercloud_cfw_domain_name_groups.filter_by_id"
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
				Config: testDataSourceCfwDomainNameGroups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.domain_names.0.domain_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.domain_names.0.domain_address_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.rules.0.name"),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_config_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCfwDomainNameGroups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  id            = sbercloud_cfw_domain_name_group.g2.id
  name          = sbercloud_cfw_domain_name_group.g2.name
  type          = tostring(sbercloud_cfw_domain_name_group.g1.type)
  config_status = "3"
}

data "sbercloud_cfw_domain_name_groups" "test" {
  fw_instance_id = "%[2]s"
  object_id      = data.sbercloud_cfw_firewalls.test.records[0].protect_objects[0].object_id

  depends_on = [
    sbercloud_cfw_protection_rule.test,
  ]
}

output "is_default_filter_useful" {
  value = length(data.sbercloud_cfw_domain_name_groups.test.records) >= 2
}

data "sbercloud_cfw_domain_name_groups" "filter_by_id" {
  fw_instance_id = "%[2]s"
  object_id      = data.sbercloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  group_id       = local.id

  depends_on = [
    sbercloud_cfw_protection_rule.test,
  ]
}

output "is_id_filter_useful" {
  value = length(data.sbercloud_cfw_domain_name_groups.filter_by_id.records) >= 1 && alltrue(
    [for v in data.sbercloud_cfw_domain_name_groups.filter_by_id.records[*] : v.group_id == local.id]
  )
}

data "sbercloud_cfw_domain_name_groups" "filter_by_name" {
  fw_instance_id = "%[2]s"
  object_id      = data.sbercloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name           = local.name

  depends_on = [
    sbercloud_cfw_protection_rule.test,
  ]
}

output "is_name_filter_useful" {
  value = length(data.sbercloud_cfw_domain_name_groups.filter_by_name.records) >= 1 && alltrue(
    [for v in data.sbercloud_cfw_domain_name_groups.filter_by_name.records[*] : v.name == local.name]
  )
}

data "sbercloud_cfw_domain_name_groups" "filter_by_type" {
  fw_instance_id = "%[2]s"
  object_id      = data.sbercloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  type           = local.type

  depends_on = [
    sbercloud_cfw_protection_rule.test,
  ]
}

output "is_type_filter_useful" {
  value = length(data.sbercloud_cfw_domain_name_groups.filter_by_type.records) >= 1 && alltrue(
    [for v in data.sbercloud_cfw_domain_name_groups.filter_by_type.records[*] : v.type == local.type]
  )
}

data "sbercloud_cfw_domain_name_groups" "filter_by_config_status" {
  fw_instance_id = "%[2]s"
  object_id      = data.sbercloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  config_status  = local.config_status

  depends_on = [
    sbercloud_cfw_protection_rule.test,
  ]
}

output "is_config_status_filter_useful" {
  value = length(data.sbercloud_cfw_domain_name_groups.filter_by_config_status.records) >= 1 && alltrue([
    for v in data.sbercloud_cfw_domain_name_groups.filter_by_config_status.records[*] : 
      v.config_status == local.config_status
  ])
}
`, testDataSourceCfwDomainNameGroups_base(name), acceptance.SBC_CFW_INSTANCE_ID)
}

func testDataSourceCfwDomainNameGroups_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_cfw_domain_name_group" "g1" {
  fw_instance_id = "%[2]s"
  object_id      = data.sbercloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name           = "%[3]s1"
  type           = 0
  description    = "created by terraform"
	  
  domain_names {
   domain_name = "www.test1.com"
   description = "test domain 1"
  }

  domain_names {
    domain_name = "www.test2.com"
    description = "test domain 2"
  }
}

resource "sbercloud_cfw_domain_name_group" "g2" {
  fw_instance_id = "%[2]s"
  object_id      = data.sbercloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name           = "%[3]s2"
  type           = 1
  description    = "created by terraform"
		
  domain_names {
    domain_name = "www.test3.com"
    description = "test domain 3"
  }

  domain_names {
   domain_name = "www.test4.com"
   description = "test domain 4"
  }
}

resource "sbercloud_cfw_protection_rule" "test" {
  name                = "%[3]s"
  object_id           = data.sbercloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  description         = "terraform test"
  type                = 0
  address_type        = 0
  action_type         = 0
  long_connect_enable = 0
  status              = 1
  direction           = 1

  source {
    type    = 0
    address = "1.1.1.1"
  }

  destination {
    type            = 6
    domain_set_id   = sbercloud_cfw_domain_name_group.g2.id
    domain_set_name = sbercloud_cfw_domain_name_group.g2.name
  }

  service {
    type = 2

    custom_service {
      protocol    = 6
      source_port = 80
      dest_port   = 80			
    }

    custom_service {
      protocol    = 6
      source_port = 8080
      dest_port   = 8080
    }
  }

  sequence {
    top = 1
  }

  tags = {
    key = "value"
  }

  depends_on = [
    sbercloud_cfw_domain_name_group.g1,
    sbercloud_cfw_domain_name_group.g2,
  ] 
}
`, testAccDatasourceFirewalls_basic(), acceptance.SBC_CFW_INSTANCE_ID, name)
}
