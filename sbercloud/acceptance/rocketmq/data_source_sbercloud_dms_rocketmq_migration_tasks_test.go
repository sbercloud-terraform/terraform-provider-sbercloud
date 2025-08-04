package rocketmq

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDmsRocketmqMigrationTasks_basic(t *testing.T) {
	dataSource := "data.sbercloud_dms_rocketmq_migration_tasks.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDmsRocketmqMigrationTasks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.#"),
					resource.TestCheckOutput("task_id_validation", "true"),
					resource.TestCheckOutput("name_validation", "true"),
					resource.TestCheckOutput("type_validation", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDmsRocketmqMigrationTasks_basic(name string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_dms_rocketmq_migration_tasks" "all" {
  depends_on = [sbercloud_dms_rocketmq_migration_task.test]

  instance_id = sbercloud_dms_rocketmq_instance.test.id
}

// filter
data "sbercloud_dms_rocketmq_migration_tasks" "test" {
  instance_id = sbercloud_dms_rocketmq_instance.test.id
  task_id     = sbercloud_dms_rocketmq_migration_task.test.id
  type        = sbercloud_dms_rocketmq_migration_task.test.type
  name        = sbercloud_dms_rocketmq_migration_task.test.name
}

locals {
  filter_results = data.sbercloud_dms_rocketmq_migration_tasks.test
}

output "task_id_validation" {
  value = alltrue([for v in local.filter_results.tasks[*].id : v == sbercloud_dms_rocketmq_migration_task.test.id])
}

output "type_validation" {
  value = alltrue([for v in local.filter_results.tasks[*].name : v == sbercloud_dms_rocketmq_migration_task.test.name])
}

output "name_validation" {
  value = alltrue([for v in local.filter_results.tasks[*].type : v == sbercloud_dms_rocketmq_migration_task.test.type])
}
`, testAccRockemqMigrationTask_rocketmq(name))
}
