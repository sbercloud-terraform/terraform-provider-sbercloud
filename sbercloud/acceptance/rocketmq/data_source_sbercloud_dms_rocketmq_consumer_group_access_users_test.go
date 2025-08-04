package rocketmq

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDmsRocketmqConsumerGroupAccessUsers_basic(t *testing.T) {
	dataSource := "data.sbercloud_dms_rocketmq_consumer_group_access_users.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDmsRocketmqConsumerGroupAccessUsers_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "policies.#"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.admin"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.perm"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.access_key"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.white_remote_address"),
				),
			},
		},
	})
}

func testDataSourceDmsRocketmqConsumerGroupAccessUsers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_dms_rocketmq_consumer_group" "test" {
  instance_id = sbercloud_dms_rocketmq_instance.test.id
  broadcast   = true

  brokers = [
    "broker-0",
  ]

  name            = "%[2]s"
  retry_max_times = "3"
  description     = "add description."
}

data "sbercloud_dms_rocketmq_consumer_group_access_users" "test" {
  depends_on = [sbercloud_dms_rocketmq_user.test, sbercloud_dms_rocketmq_consumer_group.test]

  instance_id = sbercloud_dms_rocketmq_instance.test.id
  group       = sbercloud_dms_rocketmq_consumer_group.test.name
}
`, testDmsRocketMQUser_basic(name), name)
}
