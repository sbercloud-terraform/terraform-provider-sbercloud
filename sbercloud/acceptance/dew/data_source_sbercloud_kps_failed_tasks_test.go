package dew

import (
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceKpsFailedTasks_basic(t *testing.T) {
	var (
		dataSource = "data.sbercloud_kps_failed_tasks.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDewFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKpsFailedTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.server_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.server_name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.operate_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.keypair_name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.task_error_msg"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.task_time"),
				),
			},
		},
	})
}

func testDataSourceKpsFailedTasks_basic() string {
	return (`
data "sbercloud_kps_failed_tasks" "test" {}
`)
}
