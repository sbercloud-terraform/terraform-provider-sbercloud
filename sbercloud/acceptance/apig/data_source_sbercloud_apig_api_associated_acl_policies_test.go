package apig

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceApiAssociatedAclPolicies_basic(t *testing.T) {
	var (
		rName = "data.sbercloud_apig_api_associated_acl_policies.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		byId   = "data.sbercloud_apig_api_associated_acl_policies.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.sbercloud_apig_api_associated_acl_policies.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNotFoundName   = "data.sbercloud_apig_api_associated_acl_policies.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		byType   = "data.sbercloud_apig_api_associated_acl_policies.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byEnvId   = "data.sbercloud_apig_api_associated_acl_policies.filter_by_env_id"
		dcByEnvId = acceptance.InitDataSourceCheck(byEnvId)

		byEnvName   = "data.sbercloud_apig_api_associated_acl_policies.filter_by_env_name"
		dcByEnvName = acceptance.InitDataSourceCheck(byEnvName)

		byEntityType   = "data.sbercloud_apig_api_associated_acl_policies.filter_by_entity_type"
		dcByEntityType = acceptance.InitDataSourceCheck(byEntityType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
			acceptance.TestAccPreCheckApigChannelRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApiAssociatedAclPolicies_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "policies.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_not_found_filter_useful", "true"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					dcByEnvId.CheckResourceExists(),
					resource.TestCheckOutput("is_env_id_filter_useful", "true"),
					dcByEnvName.CheckResourceExists(),
					resource.TestCheckOutput("is_env_name_filter_useful", "true"),
					dcByEntityType.CheckResourceExists(),
					resource.TestCheckOutput("is_entity_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceApiAssociatedAclPolicies_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

resource "sbercloud_apig_acl_policy_associate" "test" {
  instance_id = local.instance_id
  policy_id   = sbercloud_apig_acl_policy.test.id

  publish_ids = [
    sbercloud_apig_api_publishment.test[0].publish_id
  ]
}
`, testAccAclPolicyAssociate_base(name), name)
}

func testAccDataSourceApiAssociatedAclPolicies_basic() string {
	return fmt.Sprintf(`
%[1]s

data "sbercloud_apig_api_associated_acl_policies" "test" {
  depends_on = [
    sbercloud_apig_acl_policy_associate.test,
  ]

  instance_id = local.instance_id
  api_id      = sbercloud_apig_api.test.id
}

# Filter by ID
locals {
  policy_id = sbercloud_apig_acl_policy.test.id
}

data "sbercloud_apig_api_associated_acl_policies" "filter_by_id" {
  depends_on = [
    sbercloud_apig_acl_policy_associate.test,
  ]

  instance_id = local.instance_id
  api_id      = sbercloud_apig_api.test.id

  policy_id = local.policy_id
}

locals {
  id_filter_result = [
    for v in data.sbercloud_apig_api_associated_acl_policies.filter_by_id.policies[*].id : v == local.policy_id
  ]
}

output "is_id_filter_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

# Filter by name
locals {
  policy_name = sbercloud_apig_acl_policy.test.name
}

data "sbercloud_apig_api_associated_acl_policies" "filter_by_name" {
  depends_on = [
    sbercloud_apig_acl_policy_associate.test,
  ]

  instance_id = local.instance_id
  api_id      = sbercloud_apig_api.test.id

  name = local.policy_name
}

locals {
  name_filter_result = [
    for v in data.sbercloud_apig_api_associated_acl_policies.filter_by_name.policies[*].name : v == local.policy_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by name (not found)
locals {
  not_found_name = "not_found"
}

data "sbercloud_apig_api_associated_acl_policies" "filter_by_not_found_name" {
  depends_on = [
    sbercloud_apig_acl_policy_associate.test,
  ]

  instance_id = local.instance_id
  api_id      = sbercloud_apig_api.test.id

  name = local.not_found_name
}

locals {
  not_found_name_filter_result = [
    for v in data.sbercloud_apig_api_associated_acl_policies.filter_by_not_found_name.policies[*].name : strcontains(v, local.not_found_name)
  ]
}

output "is_name_not_found_filter_useful" {
  value = length(local.not_found_name_filter_result) == 0
}

# Filter by type
locals {
  policy_type = sbercloud_apig_acl_policy.test.type
}

data "sbercloud_apig_api_associated_acl_policies" "filter_by_type" {
  depends_on = [
    sbercloud_apig_acl_policy_associate.test,
  ]

  instance_id = local.instance_id
  api_id      = sbercloud_apig_api.test.id

  type = local.policy_type
}

locals {
  type_filter_result = [
    for v in data.sbercloud_apig_api_associated_acl_policies.filter_by_type.policies[*].type : v == local.policy_type
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

# Filter by env ID
locals {
  env_id = sbercloud_apig_environment.test[0].id
}

data "sbercloud_apig_api_associated_acl_policies" "filter_by_env_id" {
  depends_on = [
    sbercloud_apig_acl_policy_associate.test,
  ]

  instance_id = local.instance_id
  api_id      = sbercloud_apig_api.test.id

  env_id = local.env_id
}

locals {
  env_id_filter_result = [
    for v in data.sbercloud_apig_api_associated_acl_policies.filter_by_env_id.policies[*].env_id : v == local.env_id
  ]
}

output "is_env_id_filter_useful" {
  value = length(local.env_id_filter_result) > 0 && alltrue(local.env_id_filter_result)
}

# Filter by env name
locals {
  env_name = sbercloud_apig_environment.test[0].name
}

data "sbercloud_apig_api_associated_acl_policies" "filter_by_env_name" {
  depends_on = [
    sbercloud_apig_acl_policy_associate.test,
  ]

  instance_id = local.instance_id
  api_id      = sbercloud_apig_api.test.id

  env_name = local.env_name
}

locals {
  env_name_filter_result = [
    for v in data.sbercloud_apig_api_associated_acl_policies.filter_by_env_name.policies[*].env_name : v == local.env_name
  ]
}

output "is_env_name_filter_useful" {
  value = length(local.env_name_filter_result) > 0 && alltrue(local.env_name_filter_result)
}

# Filter by entity type
locals {
  entity_type = sbercloud_apig_acl_policy.test.entity_type
}

data "sbercloud_apig_api_associated_acl_policies" "filter_by_entity_type" {
  depends_on = [
    sbercloud_apig_acl_policy_associate.test,
  ]

  instance_id = local.instance_id
  api_id      = sbercloud_apig_api.test.id

  entity_type = local.entity_type
}

locals {
  entity_type_filter_result = [
    for v in data.sbercloud_apig_api_associated_acl_policies.filter_by_entity_type.policies[*].entity_type : v == local.entity_type
  ]
}

output "is_entity_type_filter_useful" {
  value = length(local.entity_type_filter_result) > 0 && alltrue(local.entity_type_filter_result)
}
`, testAccDataSourceApiAssociatedAclPolicies_base())
}
