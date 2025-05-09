package sfsturbo

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSfsTurboObsTargets_basic(t *testing.T) {
	var (
		dataSource = "data.sbercloud_sfs_turbo_obs_targets.test"
		rName      = acceptance.RandomAccResourceName()
		randInt    = acctest.RandInt()
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byTargetId   = "data.sbercloud_sfs_turbo_obs_targets.filter_by_target_id"
		dcByTargetId = acceptance.InitDataSourceCheck(byTargetId)

		byStatus   = "data.sbercloud_sfs_turbo_obs_targets.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byBucket   = "data.sbercloud_sfs_turbo_obs_targets.filter_by_bucket"
		dcByBucket = acceptance.InitDataSourceCheck(byBucket)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBSEndpoint(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSfsTurboObsTargets_basic(rName, randInt),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					dcByTargetId.CheckResourceExists(),
					resource.TestCheckOutput("target_id_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),

					dcByBucket.CheckResourceExists(),
					resource.TestCheckOutput("bucket_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceSfsTurboObsTargets_basic(name string, randInt int) string {
	return fmt.Sprintf(`
%[1]s

data "sbercloud_sfs_turbo_obs_targets" "test" {
  depends_on = [
    sbercloud_sfs_turbo_obs_target.test
  ]

  share_id = sbercloud_sfs_turbo.test.id
}

locals {
  target_id = data.sbercloud_sfs_turbo_obs_targets.test.targets[0].id
}

data "sbercloud_sfs_turbo_obs_targets" "filter_by_target_id" {
  share_id  = sbercloud_sfs_turbo.test.id
  target_id = local.target_id
}

locals {
  target_id_filter_result = [
    for v in data.sbercloud_sfs_turbo_obs_targets.filter_by_target_id.targets[*].id : v == local.target_id
  ]
}

output "target_id_filter_is_useful" {
  value = alltrue(local.target_id_filter_result) && length(local.target_id_filter_result) > 0
}

locals {
  status = data.sbercloud_sfs_turbo_obs_targets.test.targets[0].status
}

data "sbercloud_sfs_turbo_obs_targets" "filter_by_status" {
  share_id = sbercloud_sfs_turbo.test.id
  status   = local.status
}

locals {
  status_filter_result = [ 
    for v in data.sbercloud_sfs_turbo_obs_targets.filter_by_status.targets[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

locals {
  bucket = data.sbercloud_sfs_turbo_obs_targets.test.targets[0].obs[0].bucket
}

data "sbercloud_sfs_turbo_obs_targets" "filter_by_bucket" {
  share_id = sbercloud_sfs_turbo.test.id
  bucket   = local.bucket
}

locals {
  bucket_filter_result = [
    for v in data.sbercloud_sfs_turbo_obs_targets.filter_by_bucket.targets[*].obs[0].bucket : v == local.bucket
  ]
}

output "bucket_filter_is_useful" {
  value = alltrue(local.bucket_filter_result) && length(local.bucket_filter_result) > 0
}
`, testAccOBSTarget_basic(name, randInt))
}
