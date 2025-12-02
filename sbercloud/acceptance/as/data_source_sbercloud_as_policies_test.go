package as

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPoliciesDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.sbercloud_as_policies.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byScalingPolicyID   = "data.sbercloud_as_policies.scaling_policy_id_filter"
		dcByScalingPolicyID = acceptance.InitDataSourceCheck(byScalingPolicyID)

		byScalingPolicyName   = "data.sbercloud_as_policies.scaling_policy_name_filter"
		dcByScalingPolicyName = acceptance.InitDataSourceCheck(byScalingPolicyName)

		byScalingPolicyType   = "data.sbercloud_as_policies.scaling_policy_type_filter"
		dcByScalingPolicyType = acceptance.InitDataSourceCheck(byScalingPolicyType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare the AS group containing the policies in advance and configure the AS group ID into
			// the environment variable.
			acceptance.TestAccPreCheckASScalingGroupID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPoliciesDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.scaling_group_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.action.0.operation"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.action.0.instance_number"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.cool_down_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.scheduled_policy.#"),

					dcByScalingPolicyID.CheckResourceExists(),
					resource.TestCheckOutput("is_scaling_policy_id_filter_useful", "true"),

					dcByScalingPolicyName.CheckResourceExists(),
					resource.TestCheckOutput("is_scaling_policy_name_filter_useful", "true"),

					dcByScalingPolicyType.CheckResourceExists(),
					resource.TestCheckOutput("is_scaling_policy_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccPoliciesDataSource_basic() string {
	return fmt.Sprintf(`
data "sbercloud_as_policies" "test" {
  scaling_group_id = "%[1]s"
}

// Filter using scaling policy ID.
locals {
  scaling_policy_id = data.sbercloud_as_policies.test.policies[0].id
}

data "sbercloud_as_policies" "scaling_policy_id_filter" {
  scaling_group_id  = "%[1]s"
  scaling_policy_id = local.scaling_policy_id
}

locals {
  scaling_policy_id_filter_result = [
    for v in data.sbercloud_as_policies.scaling_policy_id_filter.policies[*].id : v == local.scaling_policy_id
  ]
}

output "is_scaling_policy_id_filter_useful" {
  value = alltrue(local.scaling_policy_id_filter_result) && length(local.scaling_policy_id_filter_result) > 0
}

// Filter using scaling policy name.
locals {
  scaling_policy_name = data.sbercloud_as_policies.test.policies[0].name
}

data "sbercloud_as_policies" "scaling_policy_name_filter" {
  scaling_group_id    = "%[1]s"
  scaling_policy_name = local.scaling_policy_name
}

locals {
  scaling_policy_name_filter_result = [
    for v in data.sbercloud_as_policies.scaling_policy_name_filter.policies[*].name : v == local.scaling_policy_name
  ]
}

output "is_scaling_policy_name_filter_useful" {
  value = alltrue(local.scaling_policy_name_filter_result) && length(local.scaling_policy_name_filter_result) > 0
}

// Filter using scaling policy type.
locals {
  scaling_policy_type = data.sbercloud_as_policies.test.policies[0].type
}

data "sbercloud_as_policies" "scaling_policy_type_filter" {
  scaling_group_id    = "%[1]s"
  scaling_policy_type = local.scaling_policy_type
}

locals {
  scaling_policy_type_filter_result = [
    for v in data.sbercloud_as_policies.scaling_policy_type_filter.policies[*].type : v == local.scaling_policy_type
  ]
}

output "is_scaling_policy_type_filter_useful" {
  value = alltrue(local.scaling_policy_type_filter_result) && length(local.scaling_policy_type_filter_result) > 0
}
`, acceptance.SBC_AS_SCALING_GROUP_ID)
}
