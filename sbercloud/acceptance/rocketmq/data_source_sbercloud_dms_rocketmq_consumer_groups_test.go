package rocketmq

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDmsRocketMQConsumerGroups_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSourceName := "data.sbercloud_dms_rocketmq_consumer_groups.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQSearchConsumerGroups(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.broadcast"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.retry_max_times"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("enabled_filter_is_useful", "true"),
					resource.TestCheckOutput("broadcast_filter_is_useful", "true"),
					resource.TestCheckOutput("retry_max_times_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDmsRocketMQSearchConsumerGroups(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_dms_rocketmq_consumer_group" "test" {
  instance_id = sbercloud_dms_rocketmq_instance.test.id
  
  brokers = [
    "broker-0"
  ]
  
  name            = "%[2]s"
  retry_max_times = 3
  description     = "add description."
}

data "sbercloud_dms_rocketmq_consumer_groups" "test" {
  depends_on  = [sbercloud_dms_rocketmq_consumer_group.test]
  instance_id = sbercloud_dms_rocketmq_instance.test.id
}

data "sbercloud_dms_rocketmq_consumer_groups" "name_filter" {
  depends_on  = [sbercloud_dms_rocketmq_consumer_group.test]
  instance_id = sbercloud_dms_rocketmq_instance.test.id
  name        = "%[2]s"
}
  
output "name_filter_is_useful" {
  value = length(data.sbercloud_dms_rocketmq_consumer_groups.name_filter.groups) > 0 && alltrue(
  [for v in data.sbercloud_dms_rocketmq_consumer_groups.name_filter.groups[*].name : v == "%[2]s"]
  )  
}

data "sbercloud_dms_rocketmq_consumer_groups" "enabled_filter" {
  depends_on  = [sbercloud_dms_rocketmq_consumer_group.test]
  instance_id = sbercloud_dms_rocketmq_instance.test.id
  enabled     = sbercloud_dms_rocketmq_consumer_group.test.enabled
}

locals {
  enabled = sbercloud_dms_rocketmq_consumer_group.test.enabled
}
    
output "enabled_filter_is_useful" {
  value = length(data.sbercloud_dms_rocketmq_consumer_groups.enabled_filter.groups) > 0 && alltrue(
  [for v in data.sbercloud_dms_rocketmq_consumer_groups.enabled_filter.groups[*].enabled : v == local.enabled]
  )  
}

data "sbercloud_dms_rocketmq_consumer_groups" "broadcast_filter" {
  depends_on  = [sbercloud_dms_rocketmq_consumer_group.test]
  instance_id = sbercloud_dms_rocketmq_instance.test.id
  broadcast   = sbercloud_dms_rocketmq_consumer_group.test.broadcast
}

locals {
  broadcast = sbercloud_dms_rocketmq_consumer_group.test.broadcast
}
    
output "broadcast_filter_is_useful" {
  value = length(data.sbercloud_dms_rocketmq_consumer_groups.broadcast_filter.groups) > 0 && alltrue(
  [for v in data.sbercloud_dms_rocketmq_consumer_groups.broadcast_filter.groups[*].broadcast : v == local.broadcast]
  )  
}

data "sbercloud_dms_rocketmq_consumer_groups" "retry_max_times_filter" {
  depends_on      = [sbercloud_dms_rocketmq_consumer_group.test]
  instance_id     = sbercloud_dms_rocketmq_instance.test.id
  retry_max_times = sbercloud_dms_rocketmq_consumer_group.test.retry_max_times
}

locals {
  retry_max_times = sbercloud_dms_rocketmq_consumer_group.test.retry_max_times
}
    
output "retry_max_times_filter_is_useful" {
  value = length(data.sbercloud_dms_rocketmq_consumer_groups.retry_max_times_filter.groups) > 0 && alltrue(
  [for v in data.sbercloud_dms_rocketmq_consumer_groups.retry_max_times_filter.groups[*].retry_max_times : v == local.retry_max_times]
  )
}

`, testAccDmsRocketmqConsumerGroup_version4(name), name)
}
