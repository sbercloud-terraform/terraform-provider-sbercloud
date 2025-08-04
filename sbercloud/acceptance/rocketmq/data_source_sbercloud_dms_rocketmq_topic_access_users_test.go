package rocketmq

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDmsRocketmqTopicAccessUsers_basic(t *testing.T) {
	dataSource := "data.sbercloud_dms_rocketmq_topic_access_users.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDmsRocketmqTopicAccessUsers_basic(rName),
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

func testDataSourceDmsRocketmqTopicAccessUsers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_dms_rocketmq_topic" "test" {
  instance_id = sbercloud_dms_rocketmq_instance.test.id
  name        = "%[2]s"
  queue_num   = 3

  brokers {
    name = "broker-0"
  }
}

data "sbercloud_dms_rocketmq_topic_access_users" "test" {
  depends_on = [sbercloud_dms_rocketmq_user.test, sbercloud_dms_rocketmq_topic.test]

  instance_id = sbercloud_dms_rocketmq_instance.test.id
  topic       = sbercloud_dms_rocketmq_topic.test.name
}
`, testDmsRocketMQUser_basic(name), name)
}
